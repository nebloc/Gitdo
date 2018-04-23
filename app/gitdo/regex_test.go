package main

import (
	"testing"
)

func TestCheckTaskRegex(t *testing.T) {
	testData := []struct {
		LineContent string
		ExpTask     string
	}{
		{"//TODO: Hello", "Hello"},
		{"//TODO:Hello", "Hello"},
		{"// TODO: Hello", "Hello"},
		{"// TODO:Hello", "Hello"},
		{"+// TODO: Hello", ""},
		{"#TODO: Hello", "Hello"},
		{"#TODO:Hello", "Hello"},
		{"# TODO: Hello", "Hello"},
		{"+# TODO: Hello", ""},
	}

	for _, data := range testData {
		taskName, isTask := CheckRegex(todoReg, data.LineContent)
		if !isTask && data.ExpTask != "" {
			t.Errorf("Expected to not match: %v", data)
		} else if taskName != data.ExpTask {
			t.Errorf("Expected: %s, Got: %s", data.ExpTask, taskName)
		}
	}
}

func TestCheckTagged(t *testing.T) {
	testData := []struct {
		LineContent string
		ExpFound    bool
		ExpID       string
	}{
		{"//TODO: Hello <08238>", true, "08238"},
		{"//TODO:Hello <08238>", true, "08238"},
		{"// TODO: Hello <08238>", true, "08238"},
		{"// TODO:Hello <08238>", true, "08238"},
		{"//TODO: Hello", false, ""},
		{"+// TODO: Test <fhsiufh>", false, ""},
		{"#TODO: Hello <08238>", true, "08238"},
		{"#TODO:Hello <08238>", true, "08238"},
		{"# TODO: Hello <08238>", true, "08238"},
		{"# TODO:Hello <08238>", true, "08238"},
		{"#TODO: Hello", false, ""},
		{"+# TODO: Test <fhsiufh>", false, ""},
	}

	for _, data := range testData {
		id, found := CheckRegex(taggedReg, data.LineContent)
		if found != data.ExpFound {
			t.Errorf("Line: %s\nExpected: %v, Got: %v", data.LineContent, data.ExpFound, found)
			continue
		}
		if id != data.ExpID {
			t.Errorf("Line: %s\nExpected: %v, Got: %v", data.LineContent, data.ExpID, id)
			continue
		}
	}
}
