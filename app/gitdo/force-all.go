package main

import (
	"github.com/nebloc/gitdo/app/utils"
	"github.com/urfave/cli"
	"io/ioutil"
	"strings"
	"github.com/nebloc/gitdo/app/versioncontrol"
	"fmt"
	"os"
	"sync"
)

var (
	mu            sync.Mutex
	pluginFailure = false
)

func ForceAll(ctx *cli.Context) error {
	clean := config.vc.CheckClean()
	if !clean {
		utils.Warnf("Please start with a clean repository directory.")
		return nil
	}

	confirmed := ConfirmWithUser("If a lot of tasks are found you may hit the rate limit of your task manager\nAre you sure you want to run this?")
	if !confirmed {
		return nil
	}

	hash, err := config.vc.GetHash()
	if err != nil {
		return err
	}

	branch, err := config.vc.GetBranch()
	if err != nil {
		return err
	}

	utils.Highlightf("switching to new branch to make changes on - %s", versioncontrol.NewBranchName)
	config.vc.CreateBranch()
	if err := config.vc.SwitchBranch(); err != nil {
		utils.Danger("Could not switch branch")
		return err
	}

	filesToCheck, err := config.vc.GetTrackedFiles()
	if err != nil {
		return err
	}

	taskc := make(chan Task)
	donec := make(chan struct{})
	tasks := make(map[string]Task)

	for _, file := range filesToCheck {
		go TagFile(file, hash, branch, taskc, donec)
	}

	finished := 0
	for finished < len(filesToCheck) {
		select {
		case task := <-taskc:
			tasks[task.id] = task
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
	utils.Highlight("Please run any unit tests to ensure the code's working\nCheck the diff with 'git diff HEAD~1', before merging")
	return nil
}
func TagFile(fileName string, hash string, branch string, taskc chan<- Task, donec chan<- struct{}) {
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

	donec <- struct{}{}
	return
}
