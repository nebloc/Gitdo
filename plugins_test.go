package main

import (
	"testing"
	"fmt"
)

var task = Task{
	id:       "",
	TaskName: "Test plugins",
	FileName: "main.go",
	FileLine: 7,
	Author:   "benjamin.coleman@me.com",
	Hash:     "8749387nvjnv347jnveiu703",
	Branch:   "master",
}

func TestRunPlugin(t *testing.T) {
	extGitDir, err := GetHomeDir()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(extGitDir)
	config = &Config{
		Author:            "benjamin.coleman@me.com",
		Plugin:            "test",
		PluginInterpreter: "python3",
	}
	resp, err := RunPlugin(GETID, task)
	if err != nil {
		t.Errorf("%s: %v", resp, err)
	}
	t.Log(resp)
}

func TestRunGetIDPlugin(t *testing.T) {
	internPluginDir = "./plugins/"
	config = &Config{
		Author:            "benjamin.coleman@me.com",
		Plugin:            "trello",
		PluginInterpreter: "python3",
	}
	trelloID, err := RunPlugin(GETID, task)
	if err != nil {
		t.Errorf("Didn't get ID from trello correctly: %v", err)
	}
	task.id = trelloID
}

func TestRunCreatePlugin(t *testing.T) {
	internPluginDir = "./plugins/"
	config = &Config{
		Author:            "benjamin.coleman@me.com",
		Plugin:            "trello",
		PluginInterpreter: "python3",
	}
	task.id = "0q48OchJ"
	_, err := RunPlugin(CREATE, task)
	if err != nil {
		t.Errorf("create task failed: %v", err)
	}
}
