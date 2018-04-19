package main

import (
	"os/exec"
	"fmt"
	"errors"
	"path/filepath"
)

type Git struct {
	name     string
	dir      string
	topLevel string
}

// NewGit returns a pointer to a new git implementation of the VersionControl interface.
func NewGit() *Git {
	git := new(Git)
	git.name = "Git"
	git.dir = ".git"
	return git
}

// SetHooks creates a
func (g *Git) SetHooks(homeDir string) error {
	srcHooks := filepath.Join(homeDir, "hooks")
	dstHooks := filepath.Join(g.dir, "hooks")

	err := copyFolder(srcHooks, dstHooks)
	return err
}

func (g *Git) SetTopLevel(topLevel string) {
	g.topLevel = topLevel
}

func (g *Git) GetTopLevel() string {
	return g.topLevel
}

func (g *Git) NameOfDir() string {
	return g.dir
}

func (g *Git) NameOfVC() string {
	return g.name
}

func (*Git) GetDiff() (string, error) {
	// Run a git diff to look for changes --cached to be added for
	// pre-commit hook
	cmd := exec.Command("git", "diff", "--cached")
	resp, err := cmd.CombinedOutput()

	// If error running git diff abort all
	if err != nil {
		if err.Error() == "exit status 129" {
			return "", errNotVCDir
		}
		if err, ok := err.(*exec.ExitError); ok {
			Dangerf("failed to exit git diff: %v, %v", err, stripNewlineChar(resp))
			return "", err
		}
		Danger("git diff couldn't be ran")
		return "", err
	}
	diff := stripNewlineChar(resp)
	if diff == "" {
		return "", errNoDiff
	}

	return diff, nil
}

func (*Git) RestageTasks(task Task) error {
	cmd := exec.Command("git", "add", task.FileName)
	if _, err := cmd.Output(); err != nil {
		return err
	}
	return nil
}

func (*Git) GetEmail() (string, error) {
	cmd := exec.Command("git", "config", "user.email")
	resp, err := cmd.Output()
	if err != nil {
		Warn("Please set your git email address for this repo. git config user.email example@email.com")
		return "", fmt.Errorf("Could not get user.email from git: %v", err)
	}
	return stripNewlineChar(resp), nil
}

func (*Git) Init() error {
	fmt.Println("Initializing git...")
	cmd := exec.Command("git", "init")
	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	fmt.Println("Git initialized")
	return nil
}

func (*Git) GetBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	resp, err := cmd.Output()
	if err != nil {
		return "", errors.New("could not get branch of last commit")
	}

	branch := stripNewlineChar(resp)
	return branch, nil
}

func (*Git) GetHash() (string, error) {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	resp, err := cmd.Output()
	if err != nil {
		return "", errors.New("could not get hash of last commit")
	}
	hash := stripNewlineChar(resp)
	return hash, nil
}
