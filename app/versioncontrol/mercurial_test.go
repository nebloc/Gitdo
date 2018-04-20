package versioncontrol

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const MERCURIAL_NAME string = "Mercurial"

func TestMercurial_NameOfDir(t *testing.T) {
	result := VCMap[MERCURIAL_NAME].NameOfDir()
	expected := ".hg"
	if result != expected {
		t.Errorf("Expected NameOfDir to return %s, got %s", expected, result)
	}
}

func TestMercurial_NameOfVC(t *testing.T) {
	result := VCMap[MERCURIAL_NAME].NameOfVC()
	expected := "Mercurial"
	if result != expected {
		t.Errorf("Expected NameOfVC to return %s, got %s", expected, result)
	}
}

func TestMercurial_GetDiff(t *testing.T) {
	VCMap[MERCURIAL_NAME].moveToDir(t)

	fileName := "new.txt"
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		t.Fatalf("failed to open a new file: %v", err)
	}

	_, err = file.Write([]byte("test string"))
	if err != nil {
		t.Fatalf("failed to write to new file: %v", err)
	}

	cmd := exec.Command("hg", "add", fileName)
	err = cmd.Run()
	if err != nil {
		t.Fatalf("failed to add %s to mercurial: %v", fileName, err)
	}

	diff, err := VCMap[MERCURIAL_NAME].GetDiff()
	if err != nil {
		t.Errorf("didn't expect an error in GetDiff: %v", err)
	}
	if len(strings.Split(diff, "\n")) != 6 {
		t.Errorf("Expected diff of length:\n%d\n\nGot:\n%s\n", 6, diff)
	}
}

func TestMercurial_SetHooks(t *testing.T) {
	VCMap[MERCURIAL_NAME].moveToDir(t)
	err := VCMap[MERCURIAL_NAME].SetHooks(HomeDir)
	if err != nil {
		t.Errorf("Didn't expect error setting hooks: %v", err)
	}
	hgrc := filepath.Join(VCMap[MERCURIAL_NAME].NameOfDir(), "hgrc")
	contents, err := ioutil.ReadFile(hgrc)
	if !strings.Contains(string(contents), "gitdo commit") {
		t.Errorf("Expected .hgrc to contain 'gitdo commit' command, instead: %s", contents)
	}
}
