package cmd

import (
	"testing"
)

var task = Task{
	id:       "1234",
	TaskName: "Test plugins",
	FileName: "main.go",
	FileLine: 7,
	Author:   "benjamin.coleman@me.com",
	Hash:     "8749387nvjnv347jnveiu703",
	Branch:   "master",
}

type ID string

func TestRunPlugin(t *testing.T) {

	testData := []struct {
		Command   plugcommand
		Arg       interface{}
		ExpResult string
	}{
		{
			GETID,
			task,
			"1234",
		},
		{
			CREATE,
			task,
			"Creating: 1234",
		},
		{
			DONE,
			"1234",
			"Marking 1234 as done",
		},
	}
	config = &Config{
		Author:            "benjamin.coleman@me.com",
		Plugin:            "test",
		PluginInterpreter: "python",
	}
	for _, data := range testData {
		resp, err := RunPlugin(data.Command, data.Arg)
		if err != nil {
			t.Errorf("Failed %v passed to %s: %v", data.Arg, data.Command, err)
		}
		if resp != data.ExpResult {
			t.Errorf("%s: Expected: %s Got: %s", data.Command, data.ExpResult, resp)
		}
	}
}
