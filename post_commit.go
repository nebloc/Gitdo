package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/urfave/cli"
)

// PostCommit is ran from a git post-commit hook to set the hash values and branch values of any tasks that have just
// been committed
func PostCommit(_ *cli.Context) error {
	hash, err := getHash()
	if err != nil {
		return err
	}
	branch, err := getBranch()
	if err != nil {
		return err
	}

	tasks, err := getTasksFile()
	if err != nil {
		Warn("No tasks file")
		return nil
	}
	for id, task := range tasks.Staged {
		if task.Hash == "" {
			task.Hash = hash
			task.Branch = branch
			tasks.Staged[id] = task
		}
	}

	bUpdated, err := json.MarshalIndent(tasks, "", "\t")
	if err != nil {
		Danger("couldn't marshal tasks with added hash")
		return err
	}
	err = ioutil.WriteFile(stagedTasksFile, bUpdated, os.ModePerm)
	if err != nil {
		Danger("couldn't write tasks with hash back to tasks.json")
		return err
	}
	return nil
}

// getHash runs rev-parse on git HEAD to get the latest commit
func getHash() (string, error) {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	resp, err := cmd.Output()
	if err != nil {
		return "", errors.New("could not get hash of last commit")
	}

	hash := stripNewlineChar(resp)
	return hash, nil
}

// getBranch gets the latest branch post-commit
func getBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	resp, err := cmd.Output()
	if err != nil {
		return "", errors.New("could not get branch of last commit")
	}

	branch := stripNewlineChar(resp)
	return branch, nil
}
