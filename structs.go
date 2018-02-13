package main

// TODO: What metadata should be in a task - priority, owner, etc
type Task struct {
	FileName string `json: file_name`
	TaskName string `json: task_name`
	FileLine int    `json: FileLine`
}
