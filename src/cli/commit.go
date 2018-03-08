package main

import (
	"diffparse"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	// TODO: add logrus package to support verbose options
	"log"
	"os"
	"os/exec"
	"regexp"
	"time"
)

var (
	TODOReg    *regexp.Regexp
	config     Config
	cachedFlag *bool
)

// GetDiffFromCmd runs the git diff command on the OS and returns a string of the result or the error that the cmd produced.
func GetDiffFromCmd() (string, error) {
	log.Printf("Running git diff with cached set to %v", *cachedFlag)
	// Run a git diff to look for changes --cached to be added for precommit hook
	var cmd *exec.Cmd
	if *cachedFlag {
		cmd = exec.Command("git", "diff", "--cached")
	} else {
		cmd = exec.Command("git", "diff")
	}
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

// TODO: Refactor main in to smaller functions
func main() {
	startTime := time.Now() // To Benchmark
	log.Print("Gitdo started")

	cachedFlag = flag.Bool("c", false, "a flag that adds --cached to the git diff command")
	flag.Parse()

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
		log.Print(task)
	}

	file, err := os.OpenFile("staged_tasks.json", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Print(err.Error())
		os.Exit(1)
	}
	defer file.Close()

	btask, err := json.Marshal(tasks)
	if err != nil {
		log.Print(err.Error())
		os.Exit(1)
	}
	_, err = file.Write(btask)
	if err != nil {
		log.Print(err.Error())
		os.Exit(1)
	}
	log.Print("Gitdo finished in ", time.Now().Sub(startTime))
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
