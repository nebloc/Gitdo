package main

import (
	"fmt"
	"github.com/waigani/diffparser"
	"log"
	"os"
	"os/exec"
	"sync"
)

var DIFF string = ""

func main() {
	fmt.Println("Gitdo running...")

	cmd := exec.Command("git", "diff", "-U1", "--cached")
	resp, err := cmd.Output()

	if err, ok := err.(*exec.ExitError); ok {
		log.Fatalf("Error getting diff:\n\n%s\n\nAborting commit", string(err.Stderr))
		os.Exit(1)
	}
	DIFF = fmt.Sprintf("\n%s", resp)

	diff, err := diffparser.Parse(DIFF)
	if err != nil {
		log.Fatalf("Error processing diff: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(len(diff.Files))

	for _, file := range diff.Files {
		go ProcessFile(file, &wg)
	}
	wg.Wait()

	fmt.Println("Gitdo stopping...")
}

func ProcessFile(file *diffparser.DiffFile, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println(file.Mode)
}
