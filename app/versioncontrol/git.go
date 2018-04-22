package versioncontrol

import (
	"errors"
	"fmt"
	"github.com/nebloc/gitdo/app/utils"
	"os/exec"
	"path/filepath"
)

// Git is an implementation of the VersionControl interface for the Git version control system.
type Git struct {
	name     string
	dir      string
	TopLevel string
}

// NewGit returns a pointer to a new Git implementation of the VersionControl interface.
func NewGit() *Git {
	git := new(Git)
	git.name = "Git"
	git.dir = ".git"
	return git
}

// SetHooks copies the files inside the hooks subdirectory of the given homeDir
func (g *Git) SetHooks(homeDir string) error {
	srcHooks := filepath.Join(homeDir, "hooks")
	dstHooks := filepath.Join(g.dir, "hooks")

	err := utils.CopyFolder(srcHooks, dstHooks)
	return err
}

// NameOfDir returns the hidden directory name where git stores data. Should always be ".git"
func (g *Git) NameOfDir() string {
	return g.dir
}

// NameOfVC returns the name of the version control system for printing to the user. Should always be "Git"
func (g *Git) NameOfVC() string {
	return g.name
}

// PathOfTopLevel returns the value of topLevel where the path to the root of the project is kept (e.g. dir with ".git")
func (g *Git) PathOfTopLevel() string {
	return g.TopLevel
}

// GetDiff executes a "git diff --cached" command to return the difference between files that are staged for a commit.
// Returns with an ErrNoDiff if the returned diff was empty. Results from the diff cmd are striped of ending new line
// character and returned as a string.
func (*Git) GetDiff() (string, error) {
	// Run a git diff to look for changes --cached to be added for
	// pre-commit hook
	cmd := exec.Command("git", "diff", "--cached")
	resp, err := cmd.CombinedOutput()

	// If error running git diff abort all
	if err != nil {
		// Not a Git directory
		if err.Error() == "exit status 129" {
			return "", ErrNotVCDir
		}
		return "", err
	}
	diff := utils.StripNewlineChar(resp)
	if diff == "" {
		return "", ErrNoDiff
	}

	return diff, nil
}

// RestageTasks runs a "git add" on a new task's file name to re-stage it so that the ID is in the immediate commit.
func (*Git) RestageTasks(fileName string) error {
	cmd := exec.Command("git", "add", fileName)
	if _, err := cmd.Output(); err != nil {
		return err
	}
	return nil
}

// GetEmail probes git's user.email config and returns it as a string.
func (*Git) GetEmail() (string, error) {
	cmd := exec.Command("git", "config", "user.email")
	resp, err := cmd.Output()
	if err != nil {
		utils.Warn("Please set your git email address for this repo. git config user.email example@email.com")
		return "", fmt.Errorf("Could not get user.email from git: %v", err)
	}
	return utils.StripNewlineChar(resp), nil
}

// Init Initialises a Git repository in the current directory.
func (*Git) Init() error {
	cmd := exec.Command("git", "init")
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

// GetBranch retrieves the current git branch being used.
func (*Git) GetBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	resp, err := cmd.Output()
	if err != nil {
		return "", errors.New("could not get branch of last commit")
	}

	branch := utils.StripNewlineChar(resp)
	return branch, nil
}

// GetHash retrieves the long hash of the current HEAD.
func (*Git) GetHash() (string, error) {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	resp, err := cmd.Output()
	if err != nil {
		return "", errors.New("could not get hash of last commit")
	}
	hash := utils.StripNewlineChar(resp)
	return hash, nil
}
