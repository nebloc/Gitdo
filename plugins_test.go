package main

import (
	"testing"
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
	config = &Config{
		Author:            "benjamin.coleman@me.com",
		Plugin:            "test",
		PluginInterpreter: "python3",
	}
	resp, err := RunPlugin(SETUP, "")
	if err != nil {
		t.Errorf("%s: %v", resp, err)
	}
	t.Log(resp)
}
