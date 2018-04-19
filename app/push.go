package main

import (
	"fmt"

	"github.com/urfave/cli"
	"github.com/nebloc/gitdo/app/utils"
)

// Push reads in tasks that are staged to be added, gives them to the create plugin and notifies the user that they were
// uploaded. Then moves them in to committed tasks and saves the task file. If the plugin fails, then the tasks are left
// and should be retried next 'git push'
func Push(_ *cli.Context) error {
	tasks, err := getTasksFile()
	if err != nil {
		return err
	}

	if len(tasks.NewTasks) == 0 && len(tasks.DoneTasks) == 0 {
		utils.Warn("No new tasks or done tasks")
		return nil
	}

	for id, task := range tasks.NewTasks {
		_, err = RunPlugin(CREATE, task)
		if err != nil {
			utils.Warnf("Failed to add task '%s': %v", task.String(), err)
			continue
		}
		fmt.Printf("Task %s added to %s\n", id, config.Plugin)
		tasks.RemoveTask(id)
	}

	failedIds := []string{}
	for _, id := range tasks.DoneTasks {
		_, err = RunPlugin(DONE, id)
		if err != nil {
			utils.Warnf("Failed to mark %s as done", id)
			failedIds = append(failedIds, id)
			continue
		}
		utils.Highlightf("Task %s marked as done", id)
	}
	tasks.DoneTasks = failedIds

	err = writeTasksFile(tasks)
	if err != nil {
		utils.Danger("could not save updated tasks list")
		return err
	}

	return nil
}
