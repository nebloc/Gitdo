package main

import (
	"encoding/json"
	"flag"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"time"
)

var (
	config *Config

	// Flags
	cachedFlag     *bool
	verboseLogFlag *bool
	commitMode     *bool
	pushMode       *bool
)

const (
	GitdoDir = ".git/gitdo/"

	// File name for writing and reading staged tasks from (between commit
	// and push)
	StagedTasksFile = GitdoDir + "staged_tasks.json"
)

func main() {
	startTime := time.Now() // To Benchmark

	HandleFlags()
	HandleLog()
	CheckFolder()

	log.Info("Gitdo started")

	err := LoadConfig()
	if err != nil {
		err = LoadDefaultConfig()
		if err != nil {
			log.WithError(err).Fatal("Could not set config")
		}
	}

	switch {
	case *commitMode:
		log.Debug("Starting in commit mode")
		Commit()
	case *pushMode:
		log.Debug("Starting in push mode")
		// Push()
	default:
		log.Warn("No mode given (--commit or --push)")
		PrintStaged()
	}

	log.WithField("time", time.Now().Sub(startTime)).Info("Gitdo finished")
}

// HandleFlags sets up the command line flag options and parses them
func HandleFlags() {
	verboseLogFlag = flag.Bool("v", false, "Verbose output")
	cachedFlag = flag.Bool("c", false, "Git diff ran with cached flag")
	commitMode = flag.Bool("commit", false, "Tool runs in commit mode")
	pushMode = flag.Bool("push", false, "Tool runs in push mode")
	flag.Parse()
}

// HandleLog sets up the logging level dependent on the -v (verbose) flag
func HandleLog() {
	log.SetLevel(log.InfoLevel)
	if *verboseLogFlag {
		log.SetLevel(log.DebugLevel)
	}
}

func PrintStaged() {
	bJson, err := ioutil.ReadFile(StagedTasksFile)
	if err != nil {
		log.WithError(err).Warn("Could not print staged tasks")
		return
	}
	var tasks []Task
	err = json.Unmarshal(bJson, &tasks)
	if err != nil {
		log.WithError(err).Warn("Could't print staged tasks")
	}
	log.Print(tasks)
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