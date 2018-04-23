package main

import (
	"io/ioutil"
	"os"
	"testing"
)

var origFile = []byte("1\n2\n3\n4\n5\n6\n7\n8\n9\n10")
var newFile = []byte("1\n2\n3\n4\n5\n6\n7 <1234>\n8\n9\n10")

func TestMarkSourceLines(t *testing.T) {
	fileName := "test.txt"
	setupForTest(t)
	err := ioutil.WriteFile(fileName, origFile, os.ModePerm)
	if err != nil {
		t.Fatal("Could not create test file")
	}
	task := Task{
		id:       "1234",
		FileName: fileName,
		TaskName: "7",
		FileLine: 7,
		Author:   "example@email.com",
		Hash:     "",
		Branch:   "",
	}
	err = MarkSourceLines(task)
	if err != nil {
		t.Errorf("Failed to run mark lines: %v", err)
	}

	result, err := ioutil.ReadFile(fileName)
	if err != nil {
		t.Errorf("could not read newly marked file: %v", err)
	}
	if string(result) != string(newFile) {
		t.Errorf("Expected: \n%v\n, Got: \n%v\n", newFile, result)
	}
}
