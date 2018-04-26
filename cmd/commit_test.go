package cmd

import (
	"io/ioutil"
	"os"
	"testing"
)

// func TestCheckTagged(t *testing.T) {
// 	line := diffparse.SourceLine{
// 		FileFrom: "main.go",
// 		FileTo:   "main.go",
// 		Content:  "",
// 		Position: 32,
// 		Mode:     diffparse.ADDED,
// 	}
// 	testData := []struct {
// 		LineContent string
// 		ExpFound    bool
// 		ExpID       string
// 	}{
// 		{"//TODO: Hello <08238>", true, "08238"},
// 		{"//TODO:Hello <08238>", true, "08238"},
// 		{"// TODO: Hello <08238>", true, "08238"},
// 		{"// TODO:Hello <08238>", true, "08238"},
// 		{"//TODO: Hello", false, ""},
// 		{"+// TODO: Test <fhsiufh>", false, ""},
// 		{"#TODO: Hello <08238>", true, "08238"},
// 		{"#TODO:Hello <08238>", true, "08238"},
// 		{"# TODO: Hello <08238>", true, "08238"},
// 		{"# TODO:Hello <08238>", true, "08238"},
// 		{"#TODO: Hello", false, ""},
// 		{"+# TODO: Test <fhsiufh>", false, ""},
// 	}

// 	for _, data := range testData {
// 		line.Content = data.LineContent
// 		id, found := CheckTagged(line)
// 		if found != data.ExpFound {
// 			t.Errorf("Line: %s\nExpected: %v, Got: %v", data.LineContent, data.ExpFound, found)
// 			continue
// 		}
// 		if id != data.ExpID {
// 			t.Errorf("Line: %s\nExpected: %v, Got: %v", data.LineContent, data.ExpID, id)
// 			continue
// 		}
// 	}
// }

// func TestCheckTaskRegex(t *testing.T) {
// 	testData := []struct {
// 		LineContent string
// 		ExpTask     string
// 	}{
// 		{"//TODO: Hello", "Hello"},
// 		{"//TODO:Hello", "Hello"},
// 		{"// TODO: Hello", "Hello"},
// 		{"// TODO:Hello", "Hello"},
// 		{"+// TODO: Hello", ""},
// 		{"#TODO: Hello", "Hello"},
// 		{"#TODO:Hello", "Hello"},
// 		{"# TODO: Hello", "Hello"},
// 		{"+# TODO: Hello", ""},
// 	}

// 	for _, data := range testData {
// 		match := CheckTaskRegex(data.LineContent)
// 		if len(match) == 0 && data.ExpTask != "" {
// 			t.Errorf("Expected to not match: %v", data)
// 		} else if len(match) != 0 && match[1] != data.ExpTask {
// 			t.Errorf("Expected: %s, Got: %s", data.ExpTask, match[1])
// 		}
// 	}
// }

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
