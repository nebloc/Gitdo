package main

import (
	"diffparse"
	"encoding/json"
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"time"
)

var (
	TODOReg        *regexp.Regexp
	config         Config
	cachedFlag     *bool
	verboseLogFlag *bool
)

// GetDiffFromCmd runs the git diff command on the OS and returns a string of the result or the error that the cmd produced.
func GetDiffFromCmd() (string, error) {
	log.WithFields(log.Fields{
		"cached": *cachedFlag,
	}).Debug("Running Git diff")

	// Run a git diff to look for changes --cached to be added for precommit hook
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

// GetDiffFromFile reads in the filepath specified in the config and returns a string of the contents and any read errors
func GetDiffFromFile() (string, error) {
	bDiff, err := ioutil.ReadFile(config.DiffFrom)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s", bDiff), nil
}

func HandleLog() {
	log.SetLevel(log.InfoLevel)
	if *verboseLogFlag {
		log.SetLevel(log.DebugLevel)
	}
}

func HandleFlags() {
	verboseLogFlag = flag.Bool("v", false, "verbose output")
	cachedFlag = flag.Bool("c", false, "a flag that adds --cached to the git diff command")
	flag.Parse()
}

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
		log.Fatal("No git diff output - exiting")
		os.Exit(1)
	}
	return rawDiff
}

func WriteStagedTasks(tasks []Task) {
	file, err := os.OpenFile("staged_tasks.json", os.O_RDWR|os.O_CREATE, 0644)
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

// TODO: Refactor main in to smaller functions
func main() {
	startTime := time.Now() // To Benchmark

	HandleFlags()
	HandleLog()

	log.Info("Gitdo started")

	err := LoadConfig()
	if err != nil {
		log.WithError(err).Fatal("couldn't load config")
		os.Exit(1)
	}

	rawDiff := HandleDiffSource()

	// TODO: Load from config for XXX HACK FIXME and Custom annotation
	TODOReg = regexp.MustCompile(`(?:[[:space:]]|)//(?:[[:space:]]|)TODO(?:.*):[[:space:]](.*)`)

	// Parse diff output
	lines, err := diffparse.ParseGitDiff(rawDiff)
	if err != nil {
		log.Fatalf("Error processing diff: %v", err)
		os.Exit(1)
	}

	// Loop over files and run go routines for each file changed
	tasks := ProcessDiff(lines)
	for _, task := range tasks {
		log.WithField("task", task.toString()).Debug("New task")
	}

	WriteStagedTasks(tasks)

	log.WithField("time", time.Now().Sub(startTime)).Info("Gitdo finished")
}

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
	match := TODOReg.FindStringSubmatch(line.Content)
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

type Task struct {
	FileName string `json:"file_name"`
	TaskName string `json:"task_name"`
	FileLine int    `json:"file_line"`
	Author   string `json:"author"`
}

func (t *Task) toString() string {
	return fmt.Sprintf("Author: %s, Task: %s, File: %s, Position: %d",
		t.Author, t.TaskName, t.FileName, t.FileLine)
}
