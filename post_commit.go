package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
	cli "github.com/urfave/cli"
)

func PostCommit(ctx *cli.Context) error {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	resp, err := cmd.Output()
	if err != nil {
		return errors.New("could not get hash of last commit")
	}
	hash := stripNewlineChar(resp)

	bFile, err := ioutil.ReadFile(StagedTasksFile)
	if err != nil {
		log.WithError(err).Info("No staged tasks file")
		return nil
	}

	var tasks Tasks
	err = json.Unmarshal(bFile, &tasks)
	if err != nil {
		log.WithError(err).Error("poor formatted json")
		return err
	}
	for i, task := range tasks.Staged {
		if task.Hash == "" {
			tasks.Staged[i].Hash = hash
		}
	}
	bUpdated, err := json.MarshalIndent(tasks, "", "\t")
	if err != nil {
		log.Error("couldn't marcshal tasks with added hash")
		return err
	}
	err = ioutil.WriteFile(StagedTasksFile, bUpdated, os.ModePerm)
	if err != nil {
		log.Error("couldn't write tasks with hash back to tasks.json")
		return err
	}
	return nil
}
