package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
)

var (
	// Plugin directory
	pluginDir = filepath.Join(GitdoDir, "plugins")

	//Plugin commands
	GETID  plugcommand = "getid"
	CREATE plugcommand = "create"
	DONE   plugcommand = "done"
	SETUP  plugcommand = "setup"
)

type plugcommand string

func RunPlugin(command plugcommand, elem interface{}) (string, error) {
	cmd := exec.Command(config.PluginInterpreter)     // i.e. 'python'
	cmd.Dir = filepath.Join(pluginDir, config.Plugin) // move to plugin working dir
	cmd.Args = append(cmd.Args, string(command))      // command to run
	switch {
	case command == GETID:
		if task, ok := elem.(Task); ok {
			bT, err := MarshalTask(task)
			if err != nil {
				return "", err
			}
			cmd.Args = append(cmd.Args, string(bT))
		} else {
			return "", fmt.Errorf("Passed interface not a task")
		}
	case command == CREATE:
		if task, ok := elem.(Task); ok {
			bT, err := MarshalTask(task)
			if err != nil {
				return "", err
			}
			cmd.Args = append(cmd.Args, task.id)
			cmd.Args = append(cmd.Args, string(bT))
		} else {
			return "", fmt.Errorf("Passed interface not a task")
		}
	case command == DONE:
		if id, ok := elem.(string); ok {
			cmd.Args = append(cmd.Args, id)
		} else {
			return "", fmt.Errorf("Passed interface not a string")
		}

	}
	resp, err := cmd.CombinedOutput()
	if err != nil {
		return stripNewlineChar(resp), err
	}
	return stripNewlineChar(resp), nil
}

func MarshalTask(task Task) ([]byte, error) {
	bT, err := json.MarshalIndent(task, "", "\t")
	if err != nil {
		return nil, err
	}
	return bT, nil
}

func GetPlugins() ([]string, error) {
	dirs, err := ioutil.ReadDir(pluginDir)
	if err != nil {
		return nil, err
	}
	var plugins []string

	for _, dir := range dirs {
		if dir.IsDir() {
			plugins = append(plugins, dir.Name())
		}
	}
	return plugins, nil
}
