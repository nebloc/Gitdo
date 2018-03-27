package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/mattn/go-colorable"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	config *Config

	// Flags
	cachedFlag     bool
	verboseLogFlag bool
)

const (
	GitdoDir = ".git/gitdo/"

	// File name for writing and reading staged tasks from (between commit
	// and push)
	StagedTasksFile = GitdoDir + "tasks.json"
)

// Current app version
var version string

func main() {
	gitdo := AppBuilder()
	err := gitdo.Run(os.Args)
	if err != nil {
		log.WithError(err).Fatal("Gitdo Failed.")
	}
}

// AppBuilder returns a urfave/cli app for directing commands and running setup
func AppBuilder() *cli.App {
	gitdo := cli.NewApp()
	gitdo.Name = "gitdo"
	gitdo.Usage = "track source code TODO comments"
	gitdo.Version = "0.0.0-a1"
	if version != "" {
		gitdo.Version = version
	}
	cli.VersionFlag = cli.BoolFlag{
		Name:  "version, V",
		Usage: "print the app version",
	}
	gitdo.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "verbose, v",
			Usage:       "sets logging to debug level",
			Destination: &verboseLogFlag,
		},
	}
	gitdo.Before = Setup
	gitdo.Commands = []cli.Command{
		{
			Name:   "list",
			Usage:  "prints the json of staged tasks",
			Flags:  []cli.Flag{cli.BoolFlag{Name: "config", Usage: "prints the current configuration"}},
			Action: List,
		},
		{
			Name:   "commit",
			Usage:  "gets git diff and stages any new tasks - normally ran from pre-commit hook",
			Action: Commit,
			Flags:  []cli.Flag{cli.BoolFlag{Name: "cached, c", Usage: "Diff is executed with --cached flag in commit mode", Destination: &cachedFlag}},
			After:  NotifyFinished,
		},
		{
			Name:   "init",
			Usage:  "sets the gitdo configuration for the current git repo",
			Flags:  []cli.Flag{cli.BoolFlag{Name: "with-git", Usage: "Initialises a git repo first, then gitdo"}},
			Action: Init,
		},
		{
			Name:   "post-commit",
			Usage:  "adds the commit hash that has just been committed to tasks with empty hash fields",
			Action: PostCommit,
			After:  NotifyFinished,
		},
		{
			Name:   "push",
			Usage:  "starts the plugin to move staged tasks into your task manager - normally ran from pre-push hook",
			Action: Push,
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
	log.Infof("Gitdo finished %s", ctx.Command.Name)
	return nil
}

// Setup sets the log level and makes sure that the config is set
func Setup(ctx *cli.Context) error {
	if ok, err := CheckInGit(); !ok {
		return fmt.Errorf("Could not verify gitdo is being ran from home of repository: %v\n", err)
	}
	HandleLog()
	config = &Config{}

	err := LoadConfig()
	if err == nil && config.IsSet() {
		return nil
	}
	SetConfig()
	return nil
}

// CheckInGit returns true if gitdo is being ran from the root of the git repo
func CheckInGit() (bool, error) {
	_, err := os.Stat(".git/")
	if err != nil {
		return false, err
	}
	return true, nil
}

// HandleLog sets up the logging level dependent on the -v (verbose) flag
func HandleLog() {
	if runtime.GOOS == "windows" {
		log.SetFormatter(&log.TextFormatter{ForceColors: true})
		log.SetOutput(colorable.NewColorableStdout())
	}
	log.SetLevel(log.InfoLevel)
	if verboseLogFlag {
		log.SetLevel(log.DebugLevel)
	}
}

// List pretty prints the tasks that are in file
func List(ctx *cli.Context) {
	if ctx.Bool("config") {
		fmt.Println(config.String())
		return
	}
	tasks, _ := getTasksFile()

	fmt.Println(tasks.String())
	return
}

// CheckFolder checks that the gitdo folder exists and calls Mkdir if not
func CheckFolder() error {
	if _, err := os.Stat(GitdoDir); os.IsNotExist(err) {
		err = os.Mkdir(GitdoDir, os.ModePerm|os.ModeDir)
		if err != nil {
			return err
		}
	}
	return nil
}

// stripNewLineChar takes a byte array (usually from an exec.Command run) and strips the newline characters, returning
// a string
func stripNewlineChar(orig []byte) string {
	var newStr string
	if strings.HasSuffix(string(orig), "\n") {
		newStr = string(orig)[:len(orig)-1]
	}
	if strings.HasSuffix(newStr, "\r") {
		newStr = newStr[:len(newStr)-1]
	}
	return newStr
}
