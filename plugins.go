package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	// Plugin directory
	internPluginDir = filepath.Join(gitdoDir, "plugins")

	//Plugin commands
	GETID  plugcommand = "getid"  // Needs task
	CREATE plugcommand = "create" // Needs task with ID
	DONE   plugcommand = "done"   // Needs ID
	SETUP  plugcommand = "setup"  // Needs nothing
)

type plugcommand string

var (
	ErrNotTask   = errors.New("could not cast interface to task")
	ErrNotString = errors.New("could not cast interface to task")
)

func RunPlugin(command plugcommand, elem interface{}) (string, error) {
	homeDir, err := GetHomeDir()
	if err != nil {
		return "", err
	}

	interp := strings.Split(config.PluginInterpreter, " ")
	var cmd *exec.Cmd
	if len(interp) == 1 {
		cmd = exec.Command(interp[0]) // i.e. 'python'
	} else {
		cmd = exec.Command(interp[0], interp[1:]...) // i.e. 'osascript -l JavaScript'
	}
	os.MkdirAll(filepath.Join(internPluginDir, config.Plugin), os.ModePerm) // Create plugin working dir if not exist
	cmd.Dir = filepath.Join(internPluginDir, config.Plugin)                 // move to plugin working dir

	out := bytes.Buffer{}

	cmd.Stdout = &out
	cmd.Stderr = &out

	var resp []byte

	plugin := filepath.Join(homeDir, "plugins", config.Plugin, string(command))

	cmd.Args = append(cmd.Args, plugin) // command to run
	switch command {
	case GETID:
		if task, ok := elem.(Task); ok {
			bT, err := MarshalTask(task)
			if err != nil {
				return "", err
			}
			cmd.Args = append(cmd.Args, string(bT))
		} else {
			return "", ErrNotTask
		}
	case CREATE:
		if task, ok := elem.(Task); ok {
			bT, err := MarshalTask(task)
			if err != nil {
				return "", err
			}
			cmd.Args = append(cmd.Args, task.id)
			cmd.Args = append(cmd.Args, string(bT))
		} else {
			return "", ErrNotTask
		}
	case DONE:
		if id, ok := elem.(string); ok {
			cmd.Args = append(cmd.Args, id)
		} else {
			return "", ErrNotString
		}
	case SETUP:
		// Allow cmd to have console
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	err = cmd.Run()
	resp = out.Bytes()
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
	homeDir, err := GetHomeDir()
	if err != nil {
		return nil, err
	}

	dirs, err := ioutil.ReadDir(filepath.Join(homeDir, "plugins"))
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
