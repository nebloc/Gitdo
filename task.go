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

func (t *Task) toString() string {
	return fmt.Sprintf("Author: %s, Task: %s, File: %s, Position: %d",
		t.Author, t.TaskName, t.FileName, t.FileLine)
}
