package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/tabwriter"

	log "github.com/sirupsen/logrus"
)

type Task struct {
	id       string
	FileName string `json:"file_name"`
	TaskName string `json:"task_name"`
	FileLine int    `json:"file_line"`
	Author   string `json:"author"`
	Hash     string `json:"hash"`
	Branch   string `json:"branch"`
}

// String prints the Task in a readable format
func (t *Task) String() string {
	return fmt.Sprintf("%s#%d:\t%s\tid#%s\t",
		t.FileName, t.FileLine, t.TaskName, t.id)
}

type Tasks struct {
	Staged    map[string]Task `json:"staged_task,omitempty"`
	Committed map[string]Task `json:"committed_tasks,omitempty"`
}

func (ts *Tasks) String() string {
	buf := bytes.NewBufferString("===Staged Tasks===\n")
	const padding = 2
	w := tabwriter.NewWriter(buf, 0, 0, padding, ' ', 0)

	// Print staged
	for _, task := range ts.Staged {
		fmt.Fprintln(w, task.String())
	}
	w.Flush()

	// If no staged
	if len(ts.Staged) == 0 {
		fmt.Fprintln(w, "no staged tasks")
	}

	// Print committed
	fmt.Fprintln(w, "===Commited Tasks===")
	for _, task := range ts.Committed {
		fmt.Fprintln(w, task.String())
	}
	w.Flush()

	// If no committed
	if len(ts.Committed) == 0 {
		fmt.Fprintln(w, "no committed tasks")
	}
	return strings.TrimSpace(buf.String())
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
		existingTasks.Committed[id] = task
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

func (ts *Tasks) RemoveTask(id string) {
	delete(ts.Staged, id)
}

func (ts *Tasks) StageNewTasks(newTasks []Task) {
	for _, task := range newTasks {
		ts.Staged[task.id] = task
	}
}

func (ts *Tasks) MoveTask(id string) {
	ts.Committed[id] = ts.Staged[id]
	delete(ts.Staged, id)
}
