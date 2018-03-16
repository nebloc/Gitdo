package main

import (
	"io/ioutil"
	"os"
	"runtime"
	"time"

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
	StagedTasksFile = GitdoDir + "staged_tasks.json"
)

func main() {
	startTime := time.Now() // To Benchmark

	err := LoadConfig()
	if err != nil {
		err = LoadDefaultConfig()
		if err != nil {
			log.WithError(err).Fatal("Could not set config")
		}
		err = WriteConfig()
		if err != nil {
			log.WithError(err).Warn("Couldn't save config")
		}
	}

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

	gitdo.Commands = []cli.Command{
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "prints the json of staged tasks",
			Action:  PrintStaged,
		},
		{
			Name:    "commit",
			Aliases: []string{""},
			Usage:   "gets git diff and stages any new tasks - normally ran from pre-commit",
			Action:  Commit,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "cached, c",
					Usage: "Diff is executed with --cached flag in commit mode",
				},
			},
		},
	}

	err = gitdo.Run(os.Args)
	if err != nil {
		log.WithError(err).Fatal("uh oh")
	}

	log.WithField("time", time.Now().Sub(startTime)).Info("Gitdo finished")
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

func PrintStaged(c *cli.Context) error {
	bJson, err := ioutil.ReadFile(StagedTasksFile)
	if err != nil {
		log.WithError(err).Info("No staged tasks")
		return err
	}
	log.Print(string(bJson))
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
