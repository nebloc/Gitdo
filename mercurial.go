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

// Hg is an implementation of the VersionControl interface for the Mercurial version control system.
type Hg struct {
	topLevel string
	name     string
	dir      string
}

// NewHGreturns a pointer to a new Mercurial implementation of the VersionControl interface.
func NewHg() *Hg {
	hg := new(Hg)
	hg.dir = ".hg"
	hg.name = "Mercurial"
	return hg
}

// SetHooks attempts to append the hgrc file in the homeDir to the end of the .hg/hgrc file. If the file is missing,
// it will create one.
func (h *Hg) SetHooks(homeDir string) error {
	srcHook := filepath.Join(homeDir, "hgrc")
	dstHook := filepath.Join(h.dir, "hgrc")
	err := appendFile(srcHook, dstHook)
	if err != nil {
		return fmt.Errorf("could not move .hgrc to inside %s: %v", h.dir, err)
	}
	return nil
}

// NameOfDir returns the hidden directory name where mercurial stores data. Should always be ".hg"
func (h *Hg) NameOfDir() string {
	return h.dir
}

// NameOfVC returns the name of the version control system for printing to the user. Should always be "Mercurial"
func (h *Hg) NameOfVC() string {
	return h.name
}

// PathOfTopLevel returns the value of topLevel where the path to the root of the project is kept (e.g. dir with ".hg")
func (h *Hg) PathOfTopLevel() string {
	return h.topLevel
}

// GetDiff executes a "hg diff" command to return the difference between current files that are tracked since last
// commit. Returns with an ErrNoDiff if the returned diff was empty. Results from the diff cmd are striped of ending new
// line character and returned as a string.
func (*Hg) GetDiff() (string, error) {
	cmd := exec.Command("hg", "diff")
	resp, err := cmd.CombinedOutput()
	if err != nil {
		panic("mercurial failed to give diff")
	}
	diff := stripNewlineChar(resp)
	if diff == "" {
		return "", ErrNoDiff
	}
	return diff, nil
}

// RestageTasks returns nil as there is no need to re-stage in Mercurial
func (*Hg) RestageTasks(task Task) error {
	return nil
}

// GetEmail asks the user to type their email for the project.
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

// Init Initialises a Mercurial repository in the current directory.
func (*Hg) Init() error {
	cmd := exec.Command("hg", "init")
	_, err := cmd.CombinedOutput()
	return err
}

// GetBranch retrieves the current Mercurial branch being used.
func (*Hg) GetBranch() (string, error) {
	cmd := exec.Command("Hg", "branch")
	resp, err := cmd.Output()
	if err != nil {
		return "", errors.New("could not get branch of last commit")
	}

	branch := stripNewlineChar(resp)
	return branch, nil
}

// GetHash retrieves the short hash of the current HEAD.
func (*Hg) GetHash() (string, error) {
	cmd := exec.Command("Hg", "id", "-i")
	resp, err := cmd.Output()
	if err != nil {
		return "", errors.New("could not get hash of last commit")
	}
	hash := stripNewlineChar(resp)
	return hash, nil
}
