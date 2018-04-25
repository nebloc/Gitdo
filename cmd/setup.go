package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/nebloc/gitdo/utils"
	"github.com/nebloc/gitdo/versioncontrol"
)

func setup() error {
	if err := ChangeToVCRoot(); err != nil {
		return fmt.Errorf("could not change to the root of the VCS: %v", err)
	}
	setVCPaths()
	if err := loadConfig(); err != nil {
		return fmt.Errorf("could not load configuration: %v", err)
	}
	return nil
}

// TryGitTopLevel tries to get the root directory of the project from Git, if it can't we assume it is not a
// Git project.
func TryGitTopLevel() {
	if app.vc != nil {
		return
	}

	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	result, err := cmd.Output()
	if err != nil {
		return
	}
	vc := versioncontrol.NewGit()
	vc.TopLevel = utils.StripNewlineChar(result)
	app.vc = vc
}

// TryHgTopLevel tries to get the root directory of the project from mercuruial, if it can't we assume it is not a Mercurial project.
func TryHgTopLevel() {
	if app.vc != nil {
		return
	}
	cmd := exec.Command("hg", "root")
	result, err := cmd.Output()
	if err != nil {
		return
	}
	vc := versioncontrol.NewHg()
	vc.TopLevel = utils.StripNewlineChar(result)
	app.vc = vc
}

func setVCPaths() {
	gitdoDir = filepath.Join(app.vc.NameOfDir(), "gitdo")
	// File name for writing and reading staged tasks from (between commit
	// and push)
	stagedTasksFile = filepath.Join(gitdoDir, "tasks.json")
	configFilePath = filepath.Join(gitdoDir, "config.json")
	pluginDirPath = filepath.Join(gitdoDir, "plugins")
}

// ChangeToVCRoot allows the running of Gitdo from subdirectories by moving the working dir to the top level according
// to git or mercurial
func ChangeToVCRoot() error {
	TryGitTopLevel()
	TryHgTopLevel()

	if app.vc == nil {
		return versioncontrol.ErrNotVCDir
	}
	if app.vc.PathOfTopLevel() == "" {
		return fmt.Errorf("could not determine root directory of project from %s", app.vc.NameOfVC())
	}
	err := os.Chdir(app.vc.PathOfTopLevel())
	return err
}
