package main

import (
	"fmt"
	"github.com/waigani/diffparser"
	"log"
	"os"
	"os/exec"
	"regexp"
	"sync"
)

func main() {
	fmt.Println("Gitdo running...")

	// Run a git diff to look for changes --cached to be added for precommit hook
	cmd := exec.Command("git", "diff")
	resp, err := cmd.Output()

	// If error running git diff abort all
	if err, ok := err.(*exec.ExitError); ok {
		log.Fatalf("Error getting diff:\n\n%s\n\nAborting commit", string(err.Stderr))
		os.Exit(1)
	}

	// Save output as string
	cmdOutput := fmt.Sprintf("\n%s", resp)

	// Parse diff output
	diff, err := diffparser.Parse(cmdOutput)
	if err != nil {
		log.Fatalf("Error processing diff: %v", err)
	}

	// Create waitgroup to sync handling of all files
	var wg sync.WaitGroup
	wg.Add(len(diff.Files))

	// Loop over files and run go routines for each file changed
	for _, file := range diff.Files {
		go ProcessFileDiff(file, &wg)
	}
	wg.Wait()

	fmt.Println("Gitdo stopping...")
}

func ProcessFileDiff(file *diffparser.DiffFile, wg *sync.WaitGroup) {
	defer wg.Done()

	re := regexp.MustCompile(`(?:[[:space:]]|)//(?:[[:space:]]|)TODO:[[:space:]](.*)`)

	output := fmt.Sprintf("%s\n", file.NewName)
	for _, hunk := range file.Hunks {
		for _, line := range hunk.NewRange.Lines {
			if line.Mode == 0 {
				match := re.FindStringSubmatch(line.Content)
				if len(match) > 0 {
					t := Task{
						file.NewName,
						match[1],
						line.Position,
					}
					output += t.ToString() + "\n"
				}
			}
		}
	}
	fmt.Println(output)
}

type Task struct {
	FileName string
	TaskName string
	Position int
}

func (t *Task) ToString() string {
	return fmt.Sprintf("File: %s, Task: %s, Pos: %d", t.FileName, t.TaskName, t.Position)
}
