package main

import (
	"fmt"
	"github.com/nebbers1111/diffparse"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"time"
)

var (
	TODOReg *regexp.Regexp
	config  Config
)

// GetDiffFromCmd runs the git diff command on the OS and returns a string of the result or the error that the cmd produced.
func GetDiffFromCmd() (string, error) {
	// Run a git diff to look for changes --cached to be added for precommit hook
	// cmd := exec.Command("git", "diff", "--cached")
	cmd := exec.Command("git", "diff")

	resp, err := cmd.Output()

	// If error running git diff abort all
	if err != nil {
		if err, ok := err.(*exec.ExitError); ok {
			log.Print("git diff failed to exit: ", string(err.Stderr))
			return "", err
		} else {
			log.Print("git diff couldn't be ran: ", err.Error())
			return "", err
		}
	}

	return fmt.Sprintf("%s", resp), nil
}

func GetDiffFromFile() (string, error) {
	bDiff, err := ioutil.ReadFile(config.DiffFrom)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s", bDiff), nil
}

func main() {
	startTime := time.Now() // To Benchmark
	log.Print("Gitdo started")

	err := LoadConfig()
	if err != nil {
		log.Print("couldn't load config: ", err)
		os.Exit(1)
	}

	GetDiff := GetDiffFromFile
	if config.DiffFrom == "cmd" {
		GetDiff = GetDiffFromCmd
	}

	rawDiff, err := GetDiff()
	if err != nil {
		log.Print("error getting diff: ", err.Error())
		os.Exit(1)
	} else if rawDiff == "" {
		log.Print("No git diff output - exiting")
		os.Exit(1)
	}

	// TODO: Load from config for XXX HACK FIXME and Custom annotation
	TODOReg = regexp.MustCompile(`(?:[[:space:]]|)//(?:[[:space:]]|)TODO:[[:space:]](.*)`)

	// Parse diff output
	lines, err := diffparse.ParseGitDiff(rawDiff)
	if err != nil {
		log.Fatalf("Error processing diff: %v", err)
		os.Exit(1)
	}

	// Loop over files and run go routines for each file changed
	tasks := ProcessDiff(lines)
	for _, task := range tasks {
		log.Print(task)
	}

	log.Print("Gitdo finished in ", time.Now().Sub(startTime))
}

// ProcessFileDiff Takes a diff section for a file and extracts TODO comments
func ProcessDiff(lines []diffparse.SourceLine) []Task {
	var stagedTasks []Task
	for _, line := range lines {
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
