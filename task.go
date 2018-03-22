package main

import (
	"fmt"
)

type Task struct {
	ID       string `json:"id"`
	FileName string `json:"file_name"`
	TaskName string `json:"task_name"`
	FileLine int    `json:"file_line"`
	Author   string `json:"author"`
	Hash     string `json:"hash"`
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

func (ts *Tasks) String() (str string) {
	str = "===Staged Tasks===\n"
	for _, task := range ts.Staged {
		str += fmt.Sprintf("%s\n", task.String())
	}
	if len(ts.Staged) == 0 {
		str += "no staged tasks\n"
	}
	str += "===Commited Tasks===\n"
	for _, task := range ts.Committed {
		str += fmt.Sprintf("%s\n", task.String())
	}
	if len(ts.Committed) == 0 {
		str += "no committed tasks\n"
	}
	return
}

//TODO: Change function to remove in the least number of loops possible
func (ts *Tasks) RemoveTasks(ids []string) {
	for i := len(ts.Staged) - 1; i >= 0; i-- {
		task := ts.Staged[i]
		// Condition to decide if current element has to be deleted:
		if inArray(task.ID, ids) {
			ts.Staged = append(ts.Staged[:i], ts.Staged[i+1:]...)
		}
	}
}

func inArray(taskID string, arr []string) bool {
	for _, id := range arr {
		if taskID == id {
			fmt.Printf("Deleting: %s\n", taskID)
			return true
		}
	}
	return false
}
