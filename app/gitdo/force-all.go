package main

import (
	"github.com/nebloc/gitdo/app/utils"
	"github.com/urfave/cli"
	"io/ioutil"
	"strings"
	"github.com/nebloc/gitdo/app/versioncontrol"
	"fmt"
)

func ForceAll(ctx *cli.Context) error {
	ConfirmWithUser("If a lot of tasks are found you may hit the rate limit of your task manager\nAre you sure you want to run this?")

	utils.Highlightf("Creating new branch to make changes on - %s", versioncontrol.NewBranchName)
	if err := config.vc.CreateBranch(); err != nil {
		utils.Danger("Could not Create branch")
		return err
	}
	if err := config.vc.SwitchBranch(); err != nil {
		utils.Danger("Could not Switch branch")
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
		go TagFile(file, taskc, donec)
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

	err = CommitTasks(tasks, nil)
	if err != nil {
		return err
	}

	utils.Highlight("Please run any unit tests to ensure the code's working, and check the diff, before merging back")
	return nil
}
func TagFile(fileName string, taskc chan<- Task, donec chan<- struct{}) {
	if fileName == "" {
		donec <- struct{}{}
		return
	}

	fmt.Printf("Checking file: %s\n", fileName)
	cont, err := ioutil.ReadFile(fileName)
	if err != nil {
		utils.Dangerf("Could not read: %s to tag", fileName)
		donec <- struct{}{}
		return
	}
	lines := strings.Split(string(cont), "\n")
	for ind, line := range lines {
		taskname, isTask := CheckRegex(looseTODOReg, line)
		if isTask {
			fmt.Printf("Found task: %s#L%d - %s\n", fileName, ind, taskname)
			// Create Task
			t := Task{
				id:       "",
				FileName: fileName,
				TaskName: taskname,
				FileLine: ind,
				Author:   config.Author,
				Hash:     "",
				Branch:   "",
			}

			// Get ID for task
			resp, err := RunPlugin(GETID, t)
			if err != nil {
				utils.Dangerf("couldn't get ID for task in plugin: %s, %v", resp, err)
				panic("Couldn't get id for task in plugin - " + err.Error() + resp)
			}
			t.id = resp
			taskc <- t
		}
	}
	fmt.Printf("Checking file: %s\n", fileName)
	donec <- struct{}{}
	return
}
