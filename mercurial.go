package main

import (
	"os/exec"
	"errors"
	"path/filepath"
	"fmt"
	"bufio"
	"os"
	"strings"
)

type Hg struct {
	topLevel string
	name     string
	dir      string
}

func NewHg() *Hg {
	hg := new(Hg)
	hg.dir = ".hg"
	hg.name = "Mercurial"
	return hg
}

func (h *Hg) SetHooks(homeDir string) error {
	srcHook := filepath.Join(homeDir, "hgrc")
	dstHook := filepath.Join(h.dir, "hgrc")
	err := appendFile(srcHook, dstHook)
	if err != nil {
		return fmt.Errorf("could not move .hgrc to inside %s: %v", h.dir, err)
	}
	return nil
}

func (h *Hg) SetTopLevel(topLevel string) {
	h.topLevel = topLevel
}

func (h *Hg) GetTopLevel() string {
	return h.topLevel
}

func (h *Hg) NameOfDir() string {
	return h.dir
}

func (h *Hg) NameOfVC() string {
	return h.name
}

func (*Hg) GetDiff() (string, error) {
	cmd := exec.Command("hg", "diff")
	resp, err := cmd.CombinedOutput()
	if err != nil {
		panic("mercurial failed to give diff")
	}
	diff := stripNewlineChar(resp)
	if diff == "" {
		return "", errNoDiff
	}
	return diff, nil
}

func (*Hg) RestageTasks(task Task) error {
	// No concept of staging that I can see in mercurial - all changes to files are already in the commit
	return nil
}

func (*Hg) GetEmail() (string, error) {
	// No easy way of getting email from mercurial, ask user instead
	var email string
	for email == "" {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("What email should be used: ")
		var err error
		email, err = reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		email = strings.TrimSpace(email)
	}
	return email, nil
}

func (*Hg) Init() error {
	cmd := exec.Command("Hg", "init")
	_, err := cmd.CombinedOutput()
	return err
}

func (*Hg) GetBranch() (string, error) {
	cmd := exec.Command("Hg", "branch")
	resp, err := cmd.Output()
	if err != nil {
		return "", errors.New("could not get branch of last commit")
	}

	branch := stripNewlineChar(resp)
	return branch, nil
}

func (*Hg) GetHash() (string, error) {
	cmd := exec.Command("Hg", "id", "-i")
	resp, err := cmd.Output()
	if err != nil {
		return "", errors.New("could not get hash of last commit")
	}
	hash := stripNewlineChar(resp)
	return hash, nil
}
