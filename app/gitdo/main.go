package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"path/filepath"

	"github.com/urfave/cli"
	"github.com/nebloc/gitdo/app/versioncontrol"
	"github.com/nebloc/gitdo/app/utils"
)

var (
	// Config needed for commit and push to use plugins and add author metadata
	// to task
	config = &Config{
		Author:            "",
		Plugin:            "",
		PluginInterpreter: "",
	}

	// Current version
	version string

	// Gitdo working directory (holds plugins, secrets, tasks, etc.)
	gitdoDir string

	// File name for writing and reading staged tasks from (between commit
	// and push)
	stagedTasksFile string
	configFilePath  string
	pluginDirPath   string
)

func main() {
	gitdo := AppBuilder()
	err := gitdo.Run(os.Args)
	if err != nil {
		utils.Warnf("Gitdo failed to run: %v", err)
	}
}

// AppBuilder returns a urfave/cli app for directing commands and running setup
func AppBuilder() *cli.App {
	gitdo := cli.NewApp()
	gitdo.Name = "gitdo"
	gitdo.Usage = "track source code TODO comments - https://github.com/nebloc/Gitdo"
	gitdo.Version = "0.0.0-A5"
	if version != "" {
		gitdo.Version = fmt.Sprintf("App: %s, Build: %s_%s", version, runtime.GOOS, runtime.GOARCH)
	}
	cli.VersionFlag = cli.BoolFlag{
		Name:  "version, V",
		Usage: "print the app version",
	}
	gitdo.Commands = []cli.Command{
		{
			Before: ChangeToVCRoot,
			Name:   "list",
			Usage:  "prints the json of staged tasks",
			Flags:  []cli.Flag{cli.BoolFlag{Name: "config", Usage: "prints the current configuration"}},
			Action: List,
		},
		{

			Name:   "commit",
			Usage:  "gets git diff and stages any new tasks - normally ran from pre-commit hook",
			Action: Commit,
			Before: LoadConfig,
			After:  NotifyFinished,
		},
		{
			Name:  "init",
			Usage: "sets the gitdo configuration for the current git repo",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "with-vc",
					Usage: "Must be 'Mercurial' or 'Git'. Initialises a repo first",
				},
			},
			Action: Init,
		},
		{
			Before: LoadConfig,
			Name:   "post-commit",
			Usage:  "adds the commit hash that has just been committed to tasks with empty hash fields",
			Action: PostCommit,
			After:  NotifyFinished,
		},
		{
			Name:   "push",
			Usage:  "starts the plugin to move staged tasks into your task manager - normally ran from pre-push hook",
			Action: Push,
			Before: LoadConfig,
			After:  NotifyFinished,
		},
		{
			Name:   "destroy",
			Usage:  "deletes all of the stored tasks",
			Flags:  []cli.Flag{cli.BoolFlag{Name: "yes", Usage: "confirms purge of task file"}},
			Before: ConfirmUser,
			Action: Destroy,
		},
	}
	return gitdo
}

// NotifyFinished prints that the process has finished and what command was ran
func NotifyFinished(ctx *cli.Context) error {
	log.Printf("Gitdo finished %s", ctx.Command.Name)
	return nil
}

// List pretty prints the tasks that are in file
func List(ctx *cli.Context) error {
	if ctx.Bool("config") {
		err := LoadConfig(ctx)
		if err != nil {
			return err
		}
		fmt.Println(config.String())
		return nil
	}
	tasks, _ := getTasksFile()

	fmt.Println(tasks.String())
	return nil
}

// ChangeToVCRoot allows the running of Gitdo from subdirectories by moving the working dir to the top level according
// to git or mercurial
func ChangeToVCRoot(_s *cli.Context) error {
	TryGitTopLevel()
	TryHgTopLevel()

	if config.vc == nil {
		return versioncontrol.ErrNotVCDir
	}
	if config.vc.PathOfTopLevel() == "" {
		return fmt.Errorf("could not determine root directory of project from %s", config.vc.NameOfVC())
	}
	SetVCPaths()
	err := os.Chdir(config.vc.PathOfTopLevel())
	return err
}

// TryGitTopLevel tries to get the root directory of the project from Git, if it can't we assume it is not a
// Git project.
func TryGitTopLevel() {
	if config.vc != nil {
		return
	}

	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	result, err := cmd.Output()
	if err != nil {
		return
	}
	vc := versioncontrol.NewGit()
	vc.TopLevel = utils.StripNewlineChar(result)
	config.vc = vc
}

// TryHGTopLevel tries to get the root directory of the project from mercuruial, if it can't we assume it is not a
// Git project.
func TryHgTopLevel() {
	if config.vc != nil {
		return
	}
	cmd := exec.Command("hg", "root")
	result, err := cmd.Output()
	if err != nil {
		return
	}
	vc := versioncontrol.NewHg()
	vc.TopLevel = utils.StripNewlineChar(result)
	config.vc = vc
}

func SetVCPaths() {
	gitdoDir = filepath.Join(config.vc.NameOfDir(), "gitdo")
	// File name for writing and reading staged tasks from (between commit
	// and push)
	stagedTasksFile = filepath.Join(gitdoDir, "tasks.json")
	configFilePath = filepath.Join(gitdoDir, "config.json")
	pluginDirPath = filepath.Join(gitdoDir, "plugins")
}
