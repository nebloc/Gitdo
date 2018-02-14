package main

import (
	"encoding/json"
	"fmt"
	"github.com/waigani/diffparser"
	"log"
	"os"
	"os/exec"
	"regexp"
	"time"
)

var TODOReg *regexp.Regexp
var pluginfile string = "gitdo_trello.py"

// GetDiff runs the git diff command on the OS and returns a string of the result or the error that the cmd produced.
func GetDiff() (string, error) {
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

func main() {
	startTime := time.Now() // To Benchmark
	log.Print("Gitdo started")

	rawDiff, err := GetDiff()
	if err != nil {
		os.Exit(1)
	} else if rawDiff == "" {
		log.Print("No git diff output - exiting")
		os.Exit(1)
	}

	TODOReg = regexp.MustCompile(`(?:[[:space:]]|)//(?:[[:space:]]|)TODO:[[:space:]](.*)`)

	// Parse diff output
	diff, err := diffparser.Parse(rawDiff)
	if err != nil {
		log.Fatalf("Error processing diff: %v", err)
		os.Exit(1)
	}

	// Create channel for task arrays to be returned from goroutines processing files
	taskChan := make(chan []Task)
	// Task array for all tasks to be added to
	allTasks := make([]Task, 0)

	// Loop over files and run go routines for each file changed
	for _, file := range diff.Files {
		go ProcessFileDiff(file, taskChan)
	}

	// Capture all tasks sent back and add them to the full list
	for range diff.Files {
		allTasks = append(allTasks, <-taskChan...)
	}

	RunPlugin(allTasks)

	log.Print("Gitdo finished in ", time.Now().Sub(startTime))
}

func RunPlugin(allTasks []Task) {
	// JSONify all tasks for plugins
	b, err := json.Marshal(allTasks)
	if err != nil {
		log.Fatalf("Error marshalling task array to json: %v", err)
		os.Exit(1)
	}

	// Run Plugin
	plugin := exec.Command("python3", ".git/gitdo/"+pluginfile, fmt.Sprintf("%s", b))
	resp, err := plugin.CombinedOutput()
	log.Printf("Plugin output:\n%s", resp)
	if err != nil {
		log.Fatalf("Gitdo plugin failed: %v", err)
		os.Exit(1)
	}
}

// ProcessFileDiff Takes a diff section for a file and extracts TODO comments
func ProcessFileDiff(file *diffparser.DiffFile, taskChan chan<- []Task) {
	stagedTasks := make([]Task, 0)

	// TODO: Clean up this spaghetti code
	for _, hunk := range file.Hunks { // Loop through diff hunks
		for _, line := range hunk.NewRange.Lines { // Loop over line changes
			if line.Mode == 0 { // if line was added
				task, found := CheckTask(line, file.NewName)
				if found {
					stagedTasks = append(stagedTasks, task)
				}
			}
		}
	}

	taskChan <- stagedTasks
}

//TODO: Create test function for task reg
func CheckTask(line *diffparser.DiffLine, fileName string) (Task, bool) {
	match := TODOReg.FindStringSubmatch(line.Content)
	if len(match) > 0 { // if match was found
		t := Task{
			fileName,
			match[1],
			line.Number,
		}
		return t, true
	}
	return Task{}, false
}

type Task struct {
	FileName string `json:"file_name"`
	TaskName string `json:"task_name"`
	FileLine int    `json:"file_line"`
}
