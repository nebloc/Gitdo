package main

import (
	"fmt"
)

type Task struct {
	FileName string `json:"file_name"`
	TaskName string `json:"task_name"`
	FileLine int    `json:"file_line"`
	Author   string `json:"author"`
}

// String prints the Task in a readable format
func (t *Task) String() string {
	return fmt.Sprintf("%s#%d: %s",
		t.FileName, t.FileLine, t.TaskName)
}

type Tasks struct {
	Staged    []Task `json:"staged_task,omitempty"`
	Committed []Task `json:"committed_tasks,omitempty"`
}

func (t *Tasks) String() (str string) {
	str = "===Staged Tasks===\n"
	for _, task := range t.Staged {
		str += fmt.Sprintf("%s\n", task.String())
	}
	str += "===Commited Tasks===\n"
	for _, task := range t.Committed {
		str += fmt.Sprintf("%s\n", task.String())
	}
	return
}
