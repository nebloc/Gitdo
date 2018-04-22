package main

import (
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/nebloc/gitdo/app/diffparse"
	"github.com/urfave/cli"
	"fmt"
	"github.com/nebloc/gitdo/app/versioncontrol"
	"github.com/nebloc/gitdo/app/utils"
)

// CommitTasks gets existing tasks, removes them from the task file if deleted, adds new tasks, and runs the done plugin
// where applicable
// TODO: CommitTasks should be tested <G8F6PYby>
func CommitTasks(newTasks map[string]Task, deleted map[string]bool) error {
	if len(newTasks) == 0 && len(deleted) == 0 {
		return nil
	}

	tasks, err := getTasksFile()
	if err != nil {
		utils.Warnf("Could not read existing tasks: %v", err)
	}
	for id := range deleted {
		if _, exists := tasks.NewTasks[id]; exists {
			tasks.RemoveTask(id)
		}
		tasks.DoneTasks = append(tasks.DoneTasks, id)
	}

	tasks.StageNewTasks(newTasks)

	return writeTasksFile(tasks)
}

// Commit is called when commit mode. It gathers the git diff, parses it in to
// source lines and starts the processing for tasks and writing of staged tasks.
func Commit(_ *cli.Context) error {
	rawDiff, err := config.vc.GetDiff()

	if err == versioncontrol.ErrNoDiff {
		utils.Warn("Empty diff")
		return nil
	}
	if err != nil {
		return err
	}

	// Parse diff output
	lines, err := diffparse.ParseGitDiff(rawDiff)
	if err != nil {
		utils.Dangerf("Error processing diff: %v", err)
		return err
	}

	taskChan := make(chan Task, 2)
	done := make(chan struct{})

	go SourceChanger(taskChan, done)

	changes := processDiff(lines, taskChan)
	for _, task := range changes.New {
		utils.Highlightf("new task: %v", task.String())
	}
	err = CommitTasks(changes.New, changes.Deleted)
	if err != nil {
		return err
	}
	<-done
	for _, task := range changes.New {
		err := config.vc.RestageTasks(task.FileName)
		if err != nil {
			utils.Warnf("could not restage after tagging: %v", err)
		}
	}

	utils.Highlight(changes.String())

	return nil
}

var (
	// TODO: Create a library of regex's for use with other languages. <OaTSrQjZ>
	// todoReg is a compiled regex to match the TODO comments
	todoReg = regexp.MustCompile(
		`^[[:space:]]*(?://|#)[[:space:]]*TODO(?:.*):[[:space:]]*(.*)`)
	taggedReg = regexp.MustCompile(
		`^[[:space:]]*(?://|#)[[:space:]]*TODO(?:.*):[[:space:]]*(?:.*)<(.*)>`)
)

// processDiff Takes a diff section for a file and extracts TODO comments
// TODO: Be able to support multi line todo messages. <zyWHSPaM>
func processDiff(lines []diffparse.SourceLine, taskChan chan<- Task) taskChanges {
	changes := taskChanges{
		New:     make(map[string]Task),
		Moved:   make([]string, 0),
		Deleted: make(map[string]bool, 0),
	}
	for _, line := range lines {
		id, tagged := CheckTagged(line)
		switch {
		case line.Mode == diffparse.REMOVED && tagged:
			changes.Deleted[id] = true
		case line.Mode == diffparse.ADDED && tagged:
			changes.Moved = append(changes.Moved, id)
		case line.Mode == diffparse.ADDED && !tagged:
			task, found := CheckTask(line)
			if found {
				changes.New[task.id] = task
				taskChan <- task
			}
		}
	}
	close(taskChan)
	// Remove tasks from the deleted list that were just moved
	for _, id := range changes.Moved {
		delete(changes.Deleted, id)
	}
	return changes
}

// CheckTagged runs the tagged regex and returns the ID and whether it was a match or not
func CheckTagged(line diffparse.SourceLine) (string, bool) {
	match := taggedReg.FindStringSubmatch(line.Content)
	if len(match) != 2 {
		return "", false
	}
	return match[1], true
}

// SourceChanger waits for tasks on the given taskChan, and runs MarkSourceLines
// on them. When all tasks have been sent and the channel is closed it finishes
// it's write and sends a done signal
func SourceChanger(taskChan <-chan Task, done chan<- struct{}) {
	for {
		task, open := <-taskChan
		if open {
			err := MarkSourceLines(task)
			if err != nil {
				utils.Warnf("Error tagging source: %v", err)
				continue
			}
		} else {
			close(done)
			return
		}
	}
}

// MarkSourceLines takes a task, opens it's original file and replaces the
// corresponding comments file line with the same line plus a tag in the form "<GITDO>"
func MarkSourceLines(task Task) error {
	fileCont, err := ioutil.ReadFile(task.FileName)
	if err != nil {
		utils.Dangerf("Could not mark source code as extracted: %v", err)
		return err
	}

	sep := "\n"
	lines := strings.Split(string(fileCont), sep)

	if isCRLF(lines[0]) {
		for i, line := range lines {
			lines[i] = strings.TrimSuffix(line, "\r")
			sep = "\r\n"
		}
	}

	taskIndex := task.FileLine - 1

	//Short id is used to improve readability, and file line / name helps tie short id to long
	lines[taskIndex] += " <" + task.id + ">"
	err = ioutil.WriteFile(task.FileName, []byte(strings.Join(lines, sep)), 0644)
	if err != nil {
		utils.Dangerf("could not mark source code as extracted: %v", err)
		return err
	}
	return nil
}

// isCRLF returns true if the string contains a CR at the end (LF already stripped)
func isCRLF(line string) bool {
	if strings.HasSuffix(line, "\r") {
		return true
	}
	return false
}

// CheckTaskRegex checks the line given against the todoReg and returns an array
// of the submatches
func CheckTaskRegex(line string) []string {
	return todoReg.FindStringSubmatch(line)
}

// CheckTask takes the given source line and checks for a match against the TODO regex.
// If a match is found a task is created and returned, along with a found bool
func CheckTask(line diffparse.SourceLine) (Task, bool) {
	match := CheckTaskRegex(line.Content)
	if len(match) > 0 { // if match was found

		// Create Task
		t := Task{
			id:       "",
			FileName: strings.TrimSpace(line.FileTo),
			TaskName: match[1],
			FileLine: line.Position,
			Author:   config.Author,
			Hash:     "",
			Branch:   "",
		}

		// Get ID for task
		resp, err := RunPlugin(GETID, t)
		if err != nil {
			utils.Dangerf("couldn't get ID for task in plugin: %s, %v", resp, err)
			panic("Couldn't get id for task in plugin - " + err.Error() + resp)
		}
		t.id = resp
		return t, true
	}
	return Task{}, false
}

type taskChanges struct {
	New     map[string]Task
	Deleted map[string]bool
	Moved   []string
}

func (ch *taskChanges) String() string {
	return fmt.Sprintf(
		"Tasks Added: %d, Moved: %d, Done: %d",
		len(ch.New), len(ch.Moved), len(ch.Deleted))
}
