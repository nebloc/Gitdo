package main

import (
	"errors"
	"io/ioutil"
	"os/exec"
	"regexp"
	"strings"

	"github.com/nebbers1111/gitdo/diffparse"
	"github.com/urfave/cli"
)

var (
	ErrNotGitDir = errors.New("directory is not a git repo")
	ErrNoDiff    = errors.New("diff is empty")
)

// TODO: Change diff method to be io.reader and pass file reader or exec reader instead

// GetDiffFromCmd runs the git diff command on the OS and returns a string of
// the result or the error that the cmd produced.
func GetDiffFromCmd() (string, error) {
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
			Dangerf("failed to exit git diff: %v, %v", err, stripNewlineChar(resp))
			return "", err
		} else {
			Danger("git diff couldn't be ran")
			return "", err
		}
	}
	diff := string(resp)
	if diff == "" {
		return "", ErrNoDiff
	}

	return diff, nil
}

// CommitTasks gets existing tasks, removes them from the task file if deleted, adds new tasks, and runs the done plugin
// where applicable
func CommitTasks(newTasks []Task, deleted []string) error {
	if len(newTasks) == 0 && len(deleted) == 0 {
		return nil
	}

	tasks, err := getTasksFile()
	if err != nil {
		Warnf("Could not read existing tasks: %v", err)
	}
	for _, id := range deleted {
		if _, exists := tasks.Staged[id]; exists {
			tasks.RemoveTask(id)
		} else {
			RunDonePlugin(id)
		}
	}

	tasks.StageNewTasks(newTasks)

	return writeTasksFile(tasks)
}

// Commit is called when commit mode. It gathers the git diff, parses it in to
// source lines and starts the processing for tasks and writing of staged tasks.
func Commit(_ *cli.Context) error {
	rawDiff, err := GetDiffFromCmd()
	if err != nil {
		return err
	}

	// Parse diff output
	lines, err := diffparse.ParseGitDiff(rawDiff)
	if err != nil {
		Dangerf("Error processing diff: %v", err)
		return err
	}

	taskChan := make(chan Task, 2)
	done := make(chan bool)

	go SourceChanger(taskChan, done)

	tasks, deleted := ProcessDiff(lines, taskChan)
	for _, task := range tasks {
		Highlightf("new task: %v", task.String())
	}
	err = CommitTasks(tasks, deleted)
	if err != nil {
		return err
	}
	<-done
	for _, task := range tasks {
		err := RestageTasks(task)
		if err != nil {
			Warnf("could not restage after tagging: %v", err)
		}
	}

	Highlightf("No. of tasks added: %d", len(tasks))
	return nil
}

func RestageTasks(task Task) error {
	cmd := exec.Command("git", "add", task.FileName)
	if _, err := cmd.Output(); err != nil {
		return err
	}
	return nil
}

var (
	// TODO: Should todoReg be a global variable?
	// todoReg is a compiled regex to match the TODO comments
	todoReg *regexp.Regexp = regexp.MustCompile(
		`(?:[[:space:]]|)//(?:[[:space:]]|)TODO(?:.*):[[:space:]](.*)`)
	taggedReg *regexp.Regexp = regexp.MustCompile(
		`(?:[[:space:]]|)//(?:[[:space:]]|)TODO(?:.*):[[:space:]](?:.*)<(.*)>`)
)

// ProcessFileDiff Takes a diff section for a file and extracts TODO comments
// TODO: Handle multi line todo messages
func ProcessDiff(lines []diffparse.SourceLine, taskChan chan<- Task) ([]Task, []string) {
	var stagedTasks []Task
	var deleted []string
	for _, line := range lines {
		if line.Mode == diffparse.REMOVED {
			id, found := CheckTagged(line)
			if !found {
				continue
			}
			deleted = append(deleted, id)
		}
		task, found := CheckTask(line)
		if found {
			stagedTasks = append(stagedTasks, task)
			taskChan <- task
		}
	}
	close(taskChan)
	return stagedTasks, deleted
}

func CheckTagged(line diffparse.SourceLine) (string, bool) {
	match := taggedReg.FindStringSubmatch(line.Content)
	if len(match) != 2 {
		return "", false
	}
	return match[1], true
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
				Warnf("Error tagging source: %v", err)
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
		Dangerf("Could not mark source code as extracted: %v", err)
		return err
	}
	lines := strings.Split(string(fileCont), "\n")

	taskIndex := task.FileLine - 1

	//Short id is used to improve readability, and file line / name helps tie short id to long
	lines[taskIndex] += " <" + task.id + ">"
	err = ioutil.WriteFile(task.FileName, []byte(strings.Join(lines, "\n")), 0644)
	if err != nil {
		Dangerf("could not mark source code as extracted: %v", err)
		return err
	}
	return nil
}

// CheckTask takes the given source line and checks for a match against the TODO regex.
// If a match is found a task is created and returned, along with a found bool
func CheckTask(line diffparse.SourceLine) (Task, bool) {
	tagged := taggedReg.MatchString(line.Content)
	if tagged {
		return Task{}, false
	}
	match := todoReg.FindStringSubmatch(line.Content)
	if len(match) > 0 { // if match was found
		t := Task{
			id:       "",
			FileName: line.FileTo,
			TaskName: match[1],
			FileLine: line.Position,
			Author:   config.Author,
			Hash:     "",
			Branch:   "",
		}

		id, err := RunGetIDPlugin(t)
		if err != nil {
			Dangerf("couldn't get ID for task in plugin: %v", err)
			panic("Not continuing")
		}
		t.id = id
		return t, true
	}
	return Task{}, false
}
