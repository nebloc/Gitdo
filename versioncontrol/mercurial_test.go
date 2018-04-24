package versioncontrol

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const MercurialName string = "Mercurial"

func TestMercurialNameOfDir(t *testing.T) {
	result := VCMap[MercurialName].NameOfDir()
	expected := ".hg"
	if result != expected {
		t.Errorf("Expected NameOfDir to return %s, got %s", expected, result)
	}
}

func TestMercurialNameOfVC(t *testing.T) {
	result := VCMap[MercurialName].NameOfVC()
	expected := "Mercurial"
	if result != expected {
		t.Errorf("Expected NameOfVC to return %s, got %s", expected, result)
	}
}

func TestMercurial_GetDiff(t *testing.T) {
	VCMap[MercurialName].moveToDir(t)

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

	diff, err := VCMap[MercurialName].GetDiff()
	if err != nil {
		t.Errorf("didn't expect an error in GetDiff: %v", err)
	}
	if len(strings.Split(diff, "\n")) != 6 {
		t.Errorf("Expected diff of length:\n%d\n\nGot:\n%s\n", 6, diff)
	}
}

func TestMercurial_SetHooks(t *testing.T) {
	VCMap[MercurialName].moveToDir(t)
	err := VCMap[MercurialName].SetHooks(HomeDir)
	if err != nil {
		t.Errorf("Didn't expect error setting hooks: %v", err)
		return
	}
	hgrc := filepath.Join(VCMap[MercurialName].NameOfDir(), "hgrc")
	contents, err := ioutil.ReadFile(hgrc)
	if !strings.Contains(string(contents), "gitdo commit") {
		t.Errorf("Expected .hgrc to contain 'gitdo commit' command, instead: %s", contents)
	}
}
