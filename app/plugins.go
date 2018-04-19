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
	"github.com/nebloc/gitdo/app/utils"
)

var (
	//GETID is the mode that runs the getid file in the plugin dir
	GETID plugcommand = "getid" // Needs task
	//CREATE is the mode that runs the create file in the plugin dir
	CREATE plugcommand = "create" // Needs task with ID
	//DONE is the mode that runs the done file in the plugin dir
	DONE plugcommand = "done" // Needs ID
	//SETUP is the mode that runs the setup file in the plugin dir
	SETUP plugcommand = "setup" // Needs nothing
)

type plugcommand string

var (
	errNotTask   = errors.New("could not cast interface to task")
	errNotString = errors.New("could not cast interface to string")
)

// RunPlugin will run the plugins functions depending on the mode given. It
// moves the working dir to a sub folder in .git and calls the plugin in the
// users home directory
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
	os.MkdirAll(filepath.Join(pluginDirPath, config.Plugin), os.ModePerm) // Create plugin working dir if not exist
	cmd.Dir = filepath.Join(pluginDirPath, config.Plugin)                 // move to plugin working dir

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
			return "", errNotTask
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
			return "", errNotTask
		}
	case DONE:
		if id, ok := elem.(string); ok {
			cmd.Args = append(cmd.Args, id)
		} else {
			return "", errNotString
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
		utils.Warn(utils.StripNewlineChar(resp))
		return utils.StripNewlineChar(resp), err
	}
	return utils.StripNewlineChar(resp), nil
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
