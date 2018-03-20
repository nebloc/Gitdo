package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/nebbers1111/gitdo/diffparse"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	ErrNotGitDir = errors.New("directory is not a git repo")
	ErrNoDiff    = errors.New("diff is empty")
)

// TODO: Change diff method to be io.reader and pass file reader or exec reader instead

// GetDiffFromCmd runs the git diff command on the OS and returns a string of
// the result or the error that the cmd produced.
func GetDiffFromCmd(_ *cli.Context) (string, error) {
	log.WithFields(log.Fields{
		"cached": cachedFlag,
	}).Debug("Running Git diff")

	// Run a git diff to look for changes --cached to be added for
	// precommit hook
	var cmd *exec.Cmd
	if cachedFlag {
		cmd = exec.Command("git", "diff", "--cached")
	} else {
		cmd = exec.Command("git", "diff")
	}
	resp, err := cmd.CombinedOutput()

	// If error running git diff abort all
	if err != nil {
		if err.Error() == "exit status 129" {
			return "", ErrNotGitDir
		}
		if err, ok := err.(*exec.ExitError); ok {
			log.WithFields(log.Fields{
				"exit code": err,
				"stderr":    fmt.Sprintf("%s", resp),
			}).Error("error exiting git")
			return "", err
		} else {
			log.WithError(err).Fatal("Git diff couldn't be ran")
			return "", err
		}
	}

	// TODO: Is printing the diff or length helpful?
	diff := fmt.Sprintf("%s", resp)
	log.WithFields(log.Fields{
		"length": len(diff),
	}).Debug("Returned diff")

	return diff, nil
}

// GetDiffFromFile reads in the file path specified in the config and returns a
// string of the contents and any read errors
func GetDiffFromFile(_ *cli.Context) (string, error) {
	bDiff, err := ioutil.ReadFile(config.DiffFrom)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s", bDiff), nil
}

// HandleDiffSource checks the current config and gets the diff from the
// specified source (command or file)
func HandleDiffSource(ctx *cli.Context) (string, error) {
	GetDiff := GetDiffFromFile
	if config.DiffFrom == "cmd" {
		GetDiff = GetDiffFromCmd
	}
	rawDiff, err := GetDiff(ctx)
	if err != nil {
		log.Error("error getting diff")
		return "", err
	} else if rawDiff == "" {
		log.Warn("No git diff output")
		return "", ErrNoDiff
	}
	return rawDiff, nil
}

// WriteStagedTasks writes the given task array to a staged tasks file
func WriteStagedTasks(tasks []Task) error {
	if len(tasks) == 0 {
		return nil
	}

	var existingTasks Tasks
	bExisting, err := ioutil.ReadFile(StagedTasksFile)
	if err != nil {
		log.WithError(err).Debug("No existing tasks")
	} else {
		err = json.Unmarshal(bExisting, &existingTasks)
		if err != nil {
			log.Error("Poorly formatted staged JSON")
			return err
		}

		tasks = append(existingTasks.Staged, tasks...)
	}

	existingTasks.Staged = tasks
	btask, err := json.MarshalIndent(existingTasks, "", "\t")
	if err != nil {
		log.Error("couldn't marshal tasks")
		return err
	}
	err = ioutil.WriteFile(StagedTasksFile, btask, os.ModePerm)
	if err != nil {
		log.Error("couldn't write new staged tasks")
		return err
	}
	return nil
}

// Commit is called when commit mode. It gathers the git diff, parses it in to
// source lines and starts the processing for tasks and writing of staged tasks.
func Commit(ctx *cli.Context) error {
	rawDiff, err := HandleDiffSource(ctx)
	if err != nil {
		return err
	}

	// Parse diff output
	lines, err := diffparse.ParseGitDiff(rawDiff)
	if err != nil {
		log.Errorf("Error processing diff: %v", err)
		return err
	}

	taskChan := make(chan Task, 2)
	done := make(chan bool)

	go SourceChanger(taskChan, done)

	tasks := ProcessDiff(lines, taskChan)
	for _, task := range tasks {
		log.WithField("task", task.String()).Debug("New task")
	}
	err = WriteStagedTasks(tasks)
	if err != nil {
		return err
	}
	<-done
	for _, task := range tasks {
		err := RestageTasks(task)
		if err != nil {
			log.WithError(err).Error("could not restage task after tagging")
		}
	}

	log.WithField("No. of tasks", len(tasks)).Info("Staged new tasks")
	return nil
}

func RestageTasks(task Task) error {
	cmd := exec.Command("git", "add", task.FileName)
	if _, err := cmd.Output(); err != nil {
		return err
	}
	return nil
}

// TODO: Should todoReg be a global variable?
// todoReg is a compiled regex to match the TODO comments
var todoReg *regexp.Regexp = regexp.MustCompile(
	`(?:[[:space:]]|)//(?:[[:space:]]|)TODO(?:.*):[[:space:]](.*)`)

const tag string = " <GITDO>"

// ProcessFileDiff Takes a diff section for a file and extracts TODO comments
// TODO: Handle multi line todo messages
func ProcessDiff(lines []diffparse.SourceLine, taskChan chan<- Task) []Task {
	var stagedTasks []Task
	for _, line := range lines {
		if line.Mode == diffparse.REMOVED {
			continue
		}
		task, found := CheckTask(line)
		if found {
			stagedTasks = append(stagedTasks, task)
			taskChan <- task
		}
	}
	close(taskChan)
	return stagedTasks
}

// SourceChanger waits for tasks on the given taskChan, and runs MarkSourceLines
// on them. When all tasks have been sent and the channel is closed it finishes
// it's write and sends a done signal
func SourceChanger(taskChan <-chan Task, done chan<- bool) {
	for {
		task, open := <-taskChan
		if open {
			err := MarkSourceLines(task)
			if err != nil {
				log.Errorf("error tagging source: %v", err)
				continue
			}
		} else {
			done <- true
			return
		}
	}
}

// MarkSourceLines takes a task, opens it's original file and replaces the
// corresponding comments file line with the same line plus a tag in the form "<GITDO>"
func MarkSourceLines(task Task) error {
	fileCont, err := ioutil.ReadFile(task.FileName)
	if err != nil {
		log.WithError(err).Error("Could not mark source code as extracted")
		return err
	}
	lines := strings.Split(string(fileCont), "\n")

	taskIndex := task.FileLine - 1
	lines[taskIndex] += " <GITDO>"

	err = ioutil.WriteFile(task.FileName, []byte(strings.Join(lines, "\n")), 0644)
	if err != nil {
		log.WithError(err).Error("Could not mark source code as extracted")
		return err
	}
	return nil
}

// CheckTask takes the given source line and checks for a match against the TODO regex.
// If a match is found a task is created and returned, along with a found bool
func CheckTask(line diffparse.SourceLine) (Task, bool) {
	if strings.HasSuffix(line.Content, tag) {
		return Task{}, false
	}

	match := todoReg.FindStringSubmatch(line.Content)
	if len(match) > 0 { // if match was found
		t := Task{
			line.FileTo,
			match[1],
			line.Position,
			config.Author,
			"",
		}
		return t, true
	}
	return Task{}, false
}
