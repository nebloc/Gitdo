package main

type Task struct {
	FileName string `json: file_name`
	TaskName string `json: task_name`
	FileLine int    `json: FileLine`
}
