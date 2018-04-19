package versioncontrol

import (
	"testing"
)

var git VersionControl

func init() {
	git = NewGit()
}

func TestNameOfDir(t *testing.T) {
	result := git.NameOfDir()
	expected := ".git"
	if result != expected {
		t.Errorf("Expected NameOfDir to return %s, got %s", expected, result)
	}
}

func TestNameOfVC(t *testing.T) {
	result := git.NameOfVC()
	expected := "Git"
	if result != expected {
		t.Errorf("Expected NameOfVC to return %s, got %s", expected, result)
	}
}
