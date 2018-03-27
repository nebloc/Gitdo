package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

var (
	// Plugin directory
	pluginDir = ".git/gitdo/plugins/"

	//Plugin suffixs
	getidSuffix  = "_getid"
	createSuffix = "_create"
	doneSuffix   = "_done"
)

// Called as diff is being analysed to get an id for the new task
func RunGetIDPlugin(task Task) (string, error) {
	bTask, err := json.Marshal(task)
	if err != nil {
		return "", fmt.Errorf("error marshalling task: %v\n", err)
	}
	cmd := exec.Command(config.PluginInterpreter, pluginDir+config.Plugin+getidSuffix, string(bTask))
	res, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error running getid plugin: %v: %v",
			stripNewlineChar(res), err.Error())
	}
	return stripNewlineChar(res), nil
}

// Called after diff has been analysed to delete any tasks that have been removed
func RunDonePlugin(id string) error {
	cmd := exec.Command(config.PluginInterpreter, pluginDir+config.Plugin+doneSuffix, id)
	res, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error running done plugin: %v: %v\n", string(res), err.Error())
	}
	return nil
}

// Called post-commit to create the task once hash and branch are known
func RunCreatePlugin(task Task) error {
	bTask, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("error marshalling task: %v\n", err)
	}
	cmd := exec.Command(config.PluginInterpreter, pluginDir+config.Plugin+createSuffix, task.id, string(bTask))
	res, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error running plugin: %v\n", stripNewlineChar(res))
	}
	return nil
}
