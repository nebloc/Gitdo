package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"time"

	colorable "github.com/mattn/go-colorable"
	log "github.com/sirupsen/logrus"
)

var (
	config *Config

	// Flags
	cachedFlag     *bool
	verboseLogFlag *bool
)

const (
	GitdoDir = ".git/gitdo/"

	// File name for writing and reading staged tasks from (between commit
	// and push)
	StagedTasksFile = GitdoDir + "staged_tasks.json"
)

func main() {
	startTime := time.Now() // To Benchmark

	if !SetArgs() {
		os.Exit(1)
	}
	HandleLog()
	CheckFolder()

	log.Info("Gitdo started")

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

	switch {
	case flag.Arg(0) == "commit":
		log.Debug("Starting in commit mode")
		Commit()
	case flag.Arg(0) == "push":
		log.Debug("Starting in push mode")
		// Push()
	default:
		log.Errorf("%s is not a valid command", flag.Arg(0))
	}

	log.WithField("time", time.Now().Sub(startTime)).Info("Gitdo finished")
}

func PrintOptions() {
	fmt.Fprintln(os.Stderr, "Welcome to gitdo.\n\n==Options==\nlist: Prints a list of currently staged tasks\ncommit: To be invoked by the precommit hook to add staged tasks\npush: To be invoked by the push hook\n\nflags:")
	flag.PrintDefaults()
}

// SetArgs sets up the command line flag options and parses them
func SetArgs() bool {
	verboseLogFlag = flag.Bool("v", false, "Verbose output")
	cachedFlag = flag.Bool("c", false, "Git diff ran with cached flag")
	flag.Parse()

	if flag.Arg(0) == "" {
		PrintOptions()
		return false
	}
	return true
}

// HandleLog sets up the logging level dependent on the -v (verbose) flag
func HandleLog() {
	if runtime.GOOS == "windows" {
		log.SetFormatter(&log.TextFormatter{ForceColors: true})
		log.SetOutput(colorable.NewColorableStdout())
	}
	log.SetLevel(log.InfoLevel)
	if *verboseLogFlag {
		log.SetLevel(log.DebugLevel)
	}
}

func PrintStaged() {
	bJson, err := ioutil.ReadFile(StagedTasksFile)
	if err != nil {
		log.WithError(err).Info("No staged tasks")
		return
	}
	log.Print(string(bJson))
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
