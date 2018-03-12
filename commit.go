package main

import (
	"encoding/json"
	"fmt"
	"github.com/nebbers1111/gitdo/diffparse"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
)

// GetDiffFromCmd runs the git diff command on the OS and returns a string of
// the result or the error that the cmd produced.
func GetDiffFromCmd() (string, error) {
	log.WithFields(log.Fields{
		"cached": *cachedFlag,
	}).Debug("Running Git diff")

	// Run a git diff to look for changes --cached to be added for
	// precommit hook
	var cmd *exec.Cmd
	if *cachedFlag {
		cmd = exec.Command("git", "diff", "--cached")
	} else {
		cmd = exec.Command("git", "diff")
	}
	resp, err := cmd.CombinedOutput()

	// If error running git diff abort all
	if err != nil {
		if err, ok := err.(*exec.ExitError); ok {
			log.WithFields(log.Fields{
				"exit code": err,
				"stderr":    fmt.Sprintf("%s", resp),
			}).Fatal("Git diff failed to exit")
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

// GetDiffFromFile reads in the filepath specified in the config and returns a
// string of the contents and any read errors
func GetDiffFromFile() (string, error) {
	bDiff, err := ioutil.ReadFile(config.DiffFrom)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s", bDiff), nil
}

// HandleDiffSource checks the current config and gets the diff from the
// specified source (command or file)
func HandleDiffSource() string {
	GetDiff := GetDiffFromFile
	if config.DiffFrom == "cmd" {
		GetDiff = GetDiffFromCmd
	}
	rawDiff, err := GetDiff()
	if err != nil {
		log.Fatal("error getting diff: ", err.Error())
		os.Exit(1)
	} else if rawDiff == "" {
		log.Warn("No git diff output - exiting")
		os.Exit(1)
	}
	return rawDiff
}

// WriteStagedTasks writes the given task array to a staged tasks file
func WriteStagedTasks(tasks []Task) {
	if len(tasks) < 1 {
		return
	}

	// BUG: Currently overwriting the already staged tasks rather than
	// appending
	file, err := os.OpenFile(StagedTasksFile, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}
	defer file.Close()

	btask, err := json.Marshal(tasks)
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}
	_, err = file.Write(btask)
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}
}

// Commit is called when commit mode. It gathers the git diff, parses it in to
// source lines and starts the processing for tasks and writing of staged tasks.
func Commit() {
	rawDiff := HandleDiffSource()

	// Parse diff output
	lines, err := diffparse.ParseGitDiff(rawDiff)
	if err != nil {
		log.Fatalf("Error processing diff: %v", err)
		os.Exit(1)
	}

	// Loop over files and run go routines for each file changed
	tasks := ProcessDiff(lines)
	for _, task := range tasks {
		log.WithField("task", task.String()).Debug("New task")
	}

	WriteStagedTasks(tasks)

}

// TODO: Should todoReg be a global variable?
// todoReg is a compiled regex to match the TODO comments
var todoReg *regexp.Regexp = regexp.MustCompile(
	`(?:[[:space:]]|)//(?:[[:space:]]|)TODO(?:.*):[[:space:]](.*)`)

// ProcessFileDiff Takes a diff section for a file and extracts TODO comments
func ProcessDiff(lines []diffparse.SourceLine) []Task {

	var stagedTasks []Task
	for _, line := range lines {
		if line.Mode == diffparse.REMOVED {
			continue
		}
		task, found := CheckTask(line)
		if found {
			stagedTasks = append(stagedTasks, task)
		}
	}
	return stagedTasks
}

func CheckTask(line diffparse.SourceLine) (Task, bool) {
	match := todoReg.FindStringSubmatch(line.Content)
	if len(match) > 0 { // if match was found
		t := Task{
			line.FileTo,
			match[1],
			line.Position,
			config.Author,
		}
		return t, true
	}
	return Task{}, false
}
