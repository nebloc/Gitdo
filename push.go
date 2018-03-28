package main

import (
	"fmt"
	cli "github.com/urfave/cli"
)

// Push reads in tasks that are staged to be added, gives them to the create plugin and notifies the user that they were
// uploaded. Then moves them in to committed tasks and saves the task file. If the plugin fails, then the tasks are left
// and should be retried next 'git push'
func Push(_ *cli.Context) error {
	tasks, err := getTasksFile()
	if err != nil {
		return err
	}

	if len(tasks.Staged) == 0 {
		return nil
	}

	changed := false

	for id, task := range tasks.Staged {
		err := RunCreatePlugin(task)
		if err != nil {
			Warnf("Failed to add task '%s': %v", task.String(), err)
			continue
		}
		fmt.Printf("Task %s added to %s\n", id, config.Plugin)
		changed = true
		tasks.MoveTask(id)
	}
	if changed {
		err := writeTasksFile(tasks)
		if err != nil {
			Danger("could not save updated tasks list")
			return err
		}
	}

	return nil
}
