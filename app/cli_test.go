package main

import "encoding/json"
import "fmt"
import "testing"

func TestGetDiff(t *testing.T) {
	fmt.Println("Running Test")
	_, err := GetDiff()
	if err != nil {
		t.Fail()
	}
}

func TestVet(t *testing.T) {
	example := []byte(`{"file_name": "cli.go", "task_name":"task todo", "file_line":2}`)
	var task Task
	err := json.Unmarshal(example, &task)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(task)
}
