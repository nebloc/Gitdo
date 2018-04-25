package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/nebloc/gitdo/app/utils"
	"github.com/nebloc/gitdo/app/versioncontrol"
	"github.com/urfave/cli"
)

type key int

const (
	keyHash key = iota
	keyBranch
)

var (
	mu            sync.Mutex
	pluginFailure = false

	throttle <-chan time.Time
)

func canForceAll() bool {
	clean := config.vc.CheckClean()
	if !clean {
		utils.Warnf("Please start with a clean repository directory.")
		return false
	}

	utils.Danger("This is an unstable feature.")
	confirmed := ConfirmWithUser("If a lot of tasks are found you may hit the rate limit of your task manager.\nAre you sure you want to run this?")
	if !confirmed {
		return false
	}
	return true
}

func ForceAll(cliCon *cli.Context) error {
	reqsPerSec := cliCon.Int("requests-per-second")
	if reqsPerSec <= 0 {
		return errors.New("Not a valid rate of requests")
	}

	rate := time.Second / time.Duration(reqsPerSec)
	throttle = time.Tick(rate)

	numberOfFileCrawlers := cliCon.Int("number-of-crawlers")
	if numberOfFileCrawlers <= 0 {
		return errors.New("Not a valid number of crawlers")
	}

	fmt.Printf("%d requests per second\n", reqsPerSec)
	fmt.Printf("%d crawlers\n", numberOfFileCrawlers)

	if !canForceAll() {
		return nil
	}

	filesToCheck, err := config.vc.GetTrackedFiles()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	hash, err := config.vc.GetHash()
	if err != nil {
		return err
	}
	ctx = context.WithValue(ctx, keyHash, hash)

	branch, err := config.vc.GetBranch()
	if err != nil {
		return err
	}
	ctx = context.WithValue(ctx, keyBranch, branch)

	utils.Highlightf("switching to new branch to make changes on - %s", versioncontrol.NewBranchName)
	config.vc.CreateBranch()
	if err := config.vc.SwitchBranch(); err != nil {
		utils.Danger("Could not switch branch")
		return err
	}

	// Start sending files
	filec := fileSender(ctx, filesToCheck)
	taskc := make(chan Task)
	errorc := make(chan error)
	donec := make(chan struct{})

	for x := 0; x < numberOfFileCrawlers; x++ {
		go crawlFiles(ctx, filec, taskc, errorc, donec)
	}

	tasks := make(map[string]Task)
	errorThrown := false

	finished := 0
	for finished < numberOfFileCrawlers {
		select {
		case task := <-taskc:
			tasks[task.id] = task
		case err := <-errorc:
			if !errorThrown {
				cancel()
				utils.Dangerf("Recieved error processing files, stopping...\nFirst error: %v", err)
				errorThrown = true

			}
		case <-donec:
			finished++
		}
	}

	if len(tasks) == 0 {
		utils.Highlight("No tasks found.")
		return nil
	}
	utils.Highlightf("Found %d tasks", len(tasks))

	err = CommitTasks(tasks, nil)
	if err != nil {
		return err
	}

	err = config.vc.RestageTasks(".")
	if err != nil {
		utils.Dangerf("Could not re-stage files: %v", err)
	}
	err = config.vc.NewCommit("gitdo tagged files")
	if err != nil {
		utils.Warnf("Could not commit changes: %v", err)
	}
	utils.Highlight("Please run any unit tests to ensure the code's working\nCheck the diff with 'git diff HEAD~1 -U0', before merging")
	return nil
}

func crawlFiles(ctx context.Context, filec <-chan string, taskc chan<- Task, errorc chan<- error, done chan<- struct{}) {
	processed := 0
	for {
		select {
		case <-ctx.Done():
			done <- struct{}{}
			return
		case filename, open := <-filec:
			if !open {
				done <- struct{}{}
				return
			}
			processed++
			if err := processFile(ctx, filename, taskc); err != nil {
				errorc <- err
				done <- struct{}{}
				return
			}
		default:
		}
	}
}

func processFile(ctx context.Context, filename string, taskc chan<- Task) error {
	cont, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("could not read file %s: %v", filename, err)
	}

	sep := "\n"
	lines := strings.Split(string(cont), sep)
	if isCRLF(lines[0]) {
		sep = "\r\n"
	}

	var latestError error
	changed := false

	for ind, line := range lines {
		line = utils.StripNewlineString(line)
		taskname, isTask := CheckRegex(looseTODOReg, line)
		if isTask {
			// Ignore tagged tasks
			if _, isTagged := CheckRegex(taggedReg, line); isTagged {
				continue
			}
			// Create Task
			t := Task{
				id:       "",
				FileName: filename,
				TaskName: taskname,
				FileLine: ind + 1,
				Author:   config.Author,
				Hash:     ctx.Value(keyHash).(string),
				Branch:   ctx.Value(keyBranch).(string),
			}
			select {
			case <-ctx.Done():
				break
			case <-throttle:
				// Get ID for task
				resp, err := RunPlugin(CREATE, t)
				if err != nil {
					latestError = fmt.Errorf("error creating task: %v: %s", err, resp)
					break
				}
				fmt.Printf("Found: %s#L%d - %s\n", filename, ind+1, taskname)

				t.id = resp
				taskc <- t

				fStr := "%s <%s>"
				if sep == "\r\n" {
					fStr = "%s <%s>\r"
				}

				lines[ind] = fmt.Sprintf(fStr, line, t.id)
				changed = true
			}

		}
	}

	if changed {
		err := ioutil.WriteFile(filename, []byte(strings.Join(lines, sep)), os.ModePerm)
		if err != nil {
			utils.Dangerf("Could not tag %s", filename)
		}
	}

	return latestError
}

func fileSender(ctx context.Context, files []string) chan string {
	filec := make(chan string, 5)
	go func() {
		fileCount := 0
		for _, file := range files {
			if strings.TrimSpace(file) == "" {
				continue
			}
			select {
			case <-ctx.Done():
				close(filec)
				return
			default:
			}
			filec <- file
			fileCount++
		}
		close(filec)
	}()
	return filec
}
