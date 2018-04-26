package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/nebloc/gitdo/utils"
	"github.com/nebloc/gitdo/versioncontrol"
	"github.com/spf13/cobra"
)

type key int

const (
	keyHash key = iota
	keyBranch
)

var (
	reqsPerSec           = 5
	numberOfFileCrawlers = 5
)

var forceAllCmd = &cobra.Command{
	Use:   "force-all",
	Short: "Checks and adds all files in git repository for task annotations. Does so on a new branch.",
	Run: func(cmd *cobra.Command, args []string) {
		if err := setup(); err != nil {
			pDanger("Could not load gitdo: %v\n", err)
			return
		}
		if reqsPerSec <= 0 {
			pWarning("not a valid rate of requests\n")
			return
		}
		throttle = time.Tick(time.Second / time.Duration(reqsPerSec))

		if numberOfFileCrawlers <= 0 {
			pWarning("Not a valid number of crawlers\n")
			return
		}

		fmt.Printf("%d requests per second\n", reqsPerSec)
		fmt.Printf("%d crawlers\n", numberOfFileCrawlers)

		if !canForceAll() {
			pInfo("Stopping force-all\n")
			return
		}

		if err := ForceAll(); err != nil {
			pDanger("Failed to run force-all: %v\n", err)
			return
		}

		pNormal("Gitdo finished force-all\n")
	},
}

var throttle <-chan time.Time

func canForceAll() bool {
	clean := app.vc.CheckClean()
	if !clean {
		pWarning("Please start with a clean repository directory.\n")
		return false
	}

	pDanger("This is an unstable feature.\n")
	confirmed := ConfirmWithUser("If a lot of tasks are found you may hit the rate limit of your task manager.\nAre you sure you want to run this?")
	if !confirmed {
		return false
	}
	return true
}

// ForceAll gets relevant information about the current version control state, moves to a new branch, and sets up file crawlers to find TODOs.
//The task files are then tagged, staged, and committed.
func ForceAll() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	hash, err := app.vc.GetHash()
	if err != nil {
		return err
	}
	ctx = context.WithValue(ctx, keyHash, hash)

	branch, err := app.vc.GetBranch()
	if err != nil {
		return err
	}
	ctx = context.WithValue(ctx, keyBranch, branch)

	filesToCheck, err := app.vc.GetTrackedFiles(branch)
	if err != nil {
		return err
	}

	pInfo("switching to new branch to make changes on - %s\n", versioncontrol.NewBranchName)
	app.vc.CreateBranch()
	if err := app.vc.SwitchBranch(); err != nil {
		pDanger("Could not switch branch\n")
		return err
	}

	// Start sending files
	filec := fileSender(ctx, filesToCheck)
	taskc := make(chan Task)
	errorc := make(chan error)
	donec := make(chan struct{})

	for x := 0; x < numberOfFileCrawlers; x++ {
		go crawlFiles(ctx, filec, taskc, errorc, donec)
	}

	tasks := make(map[string]Task)
	errorThrown := false

	finished := 0
	for finished < numberOfFileCrawlers {
		select {
		case task := <-taskc:
			tasks[task.id] = task
		case err := <-errorc:
			if !errorThrown {
				cancel()
				pDanger("Recieved error processing files, stopping...\nFirst error: %v\n", err)
				errorThrown = true

			}
		case <-donec:
			finished++
		}
	}

	if len(tasks) == 0 {
		pInfo("No tasks found.\n")
		return nil
	}
	pInfo("Found %d tasks\n", len(tasks))

	err = CommitTasks(tasks, nil)
	if err != nil {
		return err
	}

	err = app.vc.RestageTasks(".")
	if err != nil {
		pDanger("Could not re-stage files: %v\n", err)
	}
	err = app.vc.NewCommit("gitdo tagged files")
	if err != nil {
		pWarning("Could not commit changes: %v\n", err)
	}
	pNormal("Please run any unit tests to ensure the code's working\nCheck the diff with 'git diff HEAD~1 -U0', before merging\n")
	return nil
}

func crawlFiles(ctx context.Context, filec <-chan string, taskc chan<- Task, errorc chan<- error, done chan<- struct{}) {
	processed := 0
	for {
		select {
		case <-ctx.Done():
			done <- struct{}{}
			return
		case filename, open := <-filec:
			if !open {
				done <- struct{}{}
				return
			}
			processed++
			if err := processFile(ctx, filename, taskc); err != nil {
				errorc <- err
				done <- struct{}{}
				return
			}
		default:
		}
	}
}

func processFile(ctx context.Context, filename string, taskc chan<- Task) error {
	cont, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("could not read file %s: %v", filename, err)
	}

	sep := "\n"
	lines := strings.Split(string(cont), sep)
	if isCRLF(lines[0]) {
		sep = "\r\n"
	}

	var latestError error
	changed := false

	for ind, line := range lines {
		line = utils.StripNewlineString(line)
		taskname, isTask := CheckRegex(looseTODOReg, line)
		if isTask {
			// Ignore tagged tasks
			if _, isTagged := CheckRegex(taggedReg, line); isTagged {
				continue
			}
			// Create Task
			t := Task{
				id:       "",
				FileName: filename,
				TaskName: taskname,
				FileLine: ind + 1,
				Author:   app.Author,
				Hash:     ctx.Value(keyHash).(string),
				Branch:   ctx.Value(keyBranch).(string),
			}
			select {
			case <-ctx.Done():
				break
			case <-throttle:
				// Get ID for task
				resp, err := RunPlugin(CREATE, t)
				if err != nil {
					latestError = fmt.Errorf("error creating task: %v: %s", err, resp)
					break
				}
				fmt.Printf("Found: %s#L%d - %s\n", filename, ind+1, taskname)

				t.id = resp
				taskc <- t

				fStr := "%s <%s>"
				if sep == "\r\n" {
					fStr = "%s <%s>\r"
				}

				lines[ind] = fmt.Sprintf(fStr, line, t.id)
				changed = true
			}

		}
	}

	if changed {
		err := ioutil.WriteFile(filename, []byte(strings.Join(lines, sep)), os.ModePerm)
		if err != nil {
			pDanger("Could not tag %s", filename)
		}
	}

	return latestError
}

func fileSender(ctx context.Context, files []string) chan string {
	filec := make(chan string, 5)
	go func() {
		fileCount := 0
		for _, file := range files {
			if strings.TrimSpace(file) == "" {
				continue
			}
			select {
			case <-ctx.Done():
				close(filec)
				return
			default:
			}
			filec <- file
			fileCount++
		}
		close(filec)
	}()
	return filec
}
