package main

import (
	log "github.com/sirupsen/logrus"
	cli "github.com/urfave/cli"
)

func Push(ctx *cli.Context) error {
	tasks, err := getTasksFile()
	if err != nil {
		return err
	}

	if len(tasks.Staged) == 0 {
		return nil
	}

	for _, task := range tasks.Staged {
		err := RunActivatePlugin(task)
		if err != nil {
			log.WithError(err).Errorf("Failed to add task: %s", task.String())
		}
	}

	return nil
}

func RunActivatePlugin(task Task) error {
	log.Infof("Activated: %s - %s", task.String(), task.ID)
	return nil
}