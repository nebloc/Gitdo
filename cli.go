package main

import (
	"encoding/json"
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
	cmd := exec.Command("git", "diff", "--cached")
	resp, err := cmd.Output()

	// If error running git diff abort all
	if err, ok := err.(*exec.ExitError); ok {
		log.Fatalf("Error getting diff:\n\n%s\n\nAborting commit", string(err.Stderr))
		os.Exit(1)
	}

	// Save output as string
	cmdOutput := fmt.Sprintf("%s", resp)

	// Parse diff output
	diff, err := diffparser.Parse(cmdOutput)
	if err != nil {
		log.Fatalf("Error processing diff: %v", err)
		os.Exit(1)
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

// ProcessFileDiff Takes a diff section for a file and extracts TODO comments
func ProcessFileDiff(file *diffparser.DiffFile, wg *sync.WaitGroup) {
	defer wg.Done()

	re := regexp.MustCompile(`(?:[[:space:]]|)//(?:[[:space:]]|)TODO:[[:space:]](.*)`)

	stagedTasks := make([]Task, 0)

	// TODO: Clean up this spaghetti code
	for _, hunk := range file.Hunks { // Loop through diff hunks
		for _, line := range hunk.NewRange.Lines { // Loop over line changes
			if line.Mode == 0 { // if line was added
				match := re.FindStringSubmatch(line.Content)
				if len(match) > 0 { // if match was found
					t := Task{
						file.NewName,
						match[1],
						line.Number,
					}
					stagedTasks = append(stagedTasks, t)
				}
			}
		}
	}

	b, err := json.Marshal(stagedTasks)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", b)
}
