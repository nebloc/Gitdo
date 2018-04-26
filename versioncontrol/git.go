package versioncontrol

import (
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/nebloc/gitdo/utils"
)

// Git is an implementation of the VersionControl interface for the Git version control system.
type Git struct {
	name     string
	dir      string
	TopLevel string
}

// CheckClean verifies that the current git repository is clean
func (*Git) CheckClean() bool {
	cmd := exec.Command("git", "diff-files", "--quiet")
	err := cmd.Run()
	if err != nil {
		return false
	}
	cmd = exec.Command("git", "diff", "--quiet", "--cached")
	err = cmd.Run()
	if err != nil {
		return false
	}
	return true
}

// NewCommit creates a new git commit with the given message
func (*Git) NewCommit(message string) error {
	cmd := exec.Command("git", "commit", "-m", message)
	return cmd.Run()
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
	srcHooks := filepath.Join(homeDir, "hooks", "git")
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
	diff := utils.StripNewlineByte(resp)
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
		return "", fmt.Errorf("Could not get user.email from git: %v", err)
	}
	return utils.StripNewlineByte(resp), nil
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

	branch := utils.StripNewlineByte(resp)
	return branch, nil
}

// GetHash retrieves the long hash of the current HEAD.
func (*Git) GetHash() (string, error) {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	resp, err := cmd.Output()
	if err != nil {
		return "", errors.New("could not get hash of last commit")
	}
	hash := utils.StripNewlineByte(resp)
	return hash, nil
}

// CreateBranch creates a new git branch for gitdo to tag files on
func (*Git) CreateBranch() error {
	cmd := exec.Command("git", "branch", NewBranchName)
	err := cmd.Run()
	return err
}

// SwitchBranch attempts to switch to the GITDO_FORCED branch to safely tag source code.
func (*Git) SwitchBranch() error {
	cmd := exec.Command("git", "checkout", NewBranchName)
	return cmd.Run()
}

// GetTrackedFiles returns a string of files that are being trakced by Git
func (*Git) GetTrackedFiles(branch string) ([]string, error) {
	cmd := exec.Command("git", "ls-tree", "-r", branch, "--name-only")
	raw, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	files := strings.Split(string(raw), "\n")
	if len(files) == 0 {
		return nil, err
	}
	if strings.HasSuffix(files[0], "\r") {
		for i, fileName := range files {
			files[i] = utils.StripNewlineString(fileName)
		}
	}

	return files, nil
}
