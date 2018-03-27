package main

import "testing"

var task = Task{
	id:       "",
	TaskName: "Test plugins",
	FileName: "main.go",
	FileLine: 7,
	Author:   "benjamin.coleman@me.com",
	Hash:     "8749387nvjnv347jnveiu703",
	Branch:   "master",
}

func TestRunGetIDPlugin(t *testing.T) {
	pluginDir = "./plugins/"
	config = &Config{
		Author:            "benjamin.coleman@me.com",
		Plugin:            "trello",
		PluginInterpreter: "python3",
	}
	trelloID, err := RunGetIDPlugin(task)
	if err != nil {
		t.Errorf("Didn't get ID from trello correctly: %v", err)
	}
	task.id = trelloID
}

func TestRunCreatePlugin(t *testing.T) {
	pluginDir = "./plugins/"
	config = &Config{
		Author:            "benjamin.coleman@me.com",
		Plugin:            "trello",
		PluginInterpreter: "python3",
	}
	task.id = "0q48OchJ"
	err := RunCreatePlugin(task)
	if err != nil {
		t.Errorf("create task failed: %v", err)
	}
}
