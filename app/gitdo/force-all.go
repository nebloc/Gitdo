package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"github.com/nebloc/gitdo/app/utils"
	"github.com/nebloc/gitdo/app/versioncontrol"
	"github.com/urfave/cli"
)

type key int

const (
	keyHash key = iota
	keyBranch

	numberOfFileCrawlers = 5
)

var (
	mu            sync.Mutex
	pluginFailure = false
)

func canForceAll() bool {
	clean := config.vc.CheckClean()
	if !clean {
		utils.Warnf("Please start with a clean repository directory.")
		return false
	}

	confirmed := ConfirmWithUser("If a lot of tasks are found you may hit the rate limit of your task manager\nAre you sure you want to run this?")
	if !confirmed {
		return false
	}
	return true
}
func ForceAll(cliCon *cli.Context) error {
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
	errors := []error{}

	finished := 0
	for finished < numberOfFileCrawlers {
		select {
		case task := <-taskc:
			tasks[task.id] = task
		case err := <-errorc:
			cancel()
			errors = append(errors, err)
			fmt.Println(err)
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
			utils.Highlightf("Found: %s#L%d - %s", filename, ind+1, taskname)

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
				fmt.Println("Context canceled")
				break
			default:
			}

			// Get ID for task
			resp, err := RunPlugin(CREATE, t)
			if err != nil {
				latestError = fmt.Errorf("error creating task: %v", err)
				break
			}

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
				fmt.Printf("context done thrown")
				close(filec)
				return
			default:
			}
			filec <- file
			fileCount++
		}
		fmt.Printf("Files Checked: %d\n", fileCount)
		close(filec)
		ctx.Done()
	}()
	return filec
}

/**
func TagFile(filec chan string, hash string, branch string, taskc chan<- Task, donec chan<- struct{}) {
	for {
		fileName := <-filec
		if fileName == "" {
			donec <- struct{}{}
			return
		}

		cont, err := ioutil.ReadFile(fileName)
		if err != nil {
			utils.Dangerf("Could not read: %s to tag", fileName)
			donec <- struct{}{}
			return
		}

		sep := "\n"
		lines := strings.Split(string(cont), sep)
		if isCRLF(lines[0]) {
			sep = "\r\n"
		}

		changed := false
		for ind, line := range lines {
			line = utils.StripNewlineString(line)
			taskname, isTask := CheckRegex(looseTODOReg, line)
			if isTask {
				// Ignore tagged tasks
				if _, isTagged := CheckRegex(taggedReg, line); isTagged {
					continue
				}
				utils.Highlightf("Found: %s#L%d - %s", fileName, ind+1, taskname)

				// Create Task
				t := Task{
					id:       "",
					FileName: fileName,
					TaskName: taskname,
					FileLine: ind + 1,
					Author:   config.Author,
					Hash:     hash,
					Branch:   branch,
				}

				mu.Lock()
				if pluginFailure {
					mu.Unlock()
					break
				}
				// Get ID for task
				resp, err := RunPlugin(CREATE, t)
				if err != nil {
					pluginFailure = true
					utils.Dangerf("couldn't get ID for task in plugin: %s, %v", resp, err)
					mu.Unlock()
					break
				}
				mu.Unlock()

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

		if changed {
			err := ioutil.WriteFile(fileName, []byte(strings.Join(lines, sep)), os.ModePerm)
			if err != nil {
				utils.Dangerf("Could not tag %s", fileName)
			}
		}
	}
	donec <- struct{}{}
}

*/
