package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"

	colorable "github.com/mattn/go-colorable"
	log "github.com/sirupsen/logrus"
	cli "github.com/urfave/cli"
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

func main() {
	gitdo := AppBuilder()
	err := gitdo.Run(os.Args)
	if err != nil {
		log.WithError(err).Fatal("Gitdo Failed.")
	}
}

func AppBuilder() *cli.App {
	gitdo := cli.NewApp()
	gitdo.Name = "gitdo"
	gitdo.Usage = "track source code TODO comments"
	gitdo.Version = "0.0.0-a1"
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
			Action: List,
		},
		{
			Name:   "commit",
			Usage:  "gets git diff and stages any new tasks - normally ran from pre-commit hook",
			Action: Commit,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:        "cached, c",
					Usage:       "Diff is executed with --cached flag in commit mode",
					Destination: &cachedFlag,
				},
			},
			After: After,
		},
		{
			Name:   "init",
			Usage:  "sets the gitdo configuration for the current git repo",
			Action: Init,
		},
		{
			Name:   "post-commit",
			Usage:  "adds the commit hash that has just been committed to tasks with empty hash fields",
			Action: PostCommit,
			After:  After,
		},
		{
			Name:   "push",
			Usage:  "starts the plugin to move staged tasks into your task manager - normally ran from pre-push hook",
			Action: Push,
			After:  After,
		},
	}
	return gitdo
}

func After(ctx *cli.Context) error {
	log.Infof("Gitdo finished %s", ctx.Command.Name)
	return nil
}

func Setup(ctx *cli.Context) error {
	CheckFolder()
	HandleLog()
	err := LoadConfig()
	if err != nil {
		err = LoadDefaultConfig()
		if err != nil {
			log.WithError(err).Fatal("Could not get config")
			cli.NewExitError("Could not get email address from git user.Email", 2)
		}
		err = WriteConfig()
		if err != nil {
			log.WithError(err).Warn("Couldn't save config")
		}
		return nil
	}
	return nil
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

func List(c *cli.Context) error {
	bJson, err := ioutil.ReadFile(StagedTasksFile)
	if err != nil {
		log.WithError(err).Info("No staged tasks")
		return err
	}

	var tasks Tasks
	err = json.Unmarshal(bJson, &tasks)
	if err != nil {
		log.WithError(err).Error("poor formatted json")
		return err
	}
	fmt.Println(tasks.String())
	return nil
}

func CheckFolder() error {
	if _, err := os.Stat(GitdoDir); os.IsNotExist(err) {
		err = os.Mkdir(GitdoDir, os.ModePerm|os.ModeDir)
		if err != nil {
			return err
		}
	}
	return nil
}
