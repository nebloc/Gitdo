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
	return fmt.Sprintf("Author: %s, Task: %s, File: %s, Position: %d",
		t.Author, t.TaskName, t.FileName, t.FileLine)
}

func (t1 *Task) CheckEqual(t2 Task) bool {
	return false
}
