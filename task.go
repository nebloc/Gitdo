package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
)

type Task struct {
	id       string
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
	Staged    map[string]Task `json:"staged_task,omitempty"`
	Committed map[string]Task `json:"committed_tasks,omitempty"`
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

func getTasksFile() (*Tasks, error) {
	existingTasks := NewTaskMap()

	bExisting, err := ioutil.ReadFile(StagedTasksFile)
	if err != nil {
		return existingTasks, err
	}
	err = json.Unmarshal(bExisting, &existingTasks)
	if err != nil {
		log.Error("Poorly formatted staged JSON")
		return existingTasks, err
	}
	for id, task := range existingTasks.Staged {
		task.id = id
		existingTasks.Staged[id] = task
	}
	for id, task := range existingTasks.Committed {
		task.id = id
		existingTasks.Staged[id] = task
	}

	return existingTasks, nil
}

func NewTaskMap() *Tasks {
	return &Tasks{
		Staged:    make(map[string]Task),
		Committed: make(map[string]Task),
	}
}

func writeTasksFile(tasks *Tasks) error {
	btask, err := json.MarshalIndent(*tasks, "", "\t")
	if err != nil {
		log.Error("couldn't marshal tasks")
		return err
	}
	err = ioutil.WriteFile(StagedTasksFile, btask, os.ModePerm)
	if err != nil {
		log.Error("couldn't write new staged tasks")
		return err
	}
	return nil
}

func (ts *Tasks) RemoveStagedTasks(ids []string) {
	for _, id := range ids {
		delete(ts.Staged, id)
	}
}

func (ts *Tasks) StageNewTasks(newTasks []Task) {
	for _, task := range newTasks {
		ts.Staged[task.id] = task
	}
}
