package main

import (
	"fmt"
	"os/exec"

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

	for id, task := range tasks.Staged {
		err := RunActivatePlugin(id)
		if err != nil {
			log.WithError(err).Errorf("Failed to add task: %s", task.String())
		}
	}

	return nil
}

func RunActivatePlugin(id string) error {
	cmd := exec.Command(".git/gitdo/plugins/activate_"+config.Plugin, id)
	res, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error running plugin: %v\n", stripNewlineChar(res))
	}
	return nil
}
