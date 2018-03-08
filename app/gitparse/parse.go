package gitparse

import (
	"fmt"
	"os/exec"
	"strings"
)

func GetGitDiff() (string, error) {
	gitCmd := exec.Command("git", "diff", "--cached")
	bDiff, err := gitCmd.Output()
	if err != nil {
		if err, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("git diff failed to exit: %v", err.Stderr)
		} else {
			return "", fmt.Errorf("git diff couldn't be run: %v", err.Stderr)
		}
	}

	return fmt.Sprintf("%s", bDiff), nil
}

const fromFilePrefix = "--- a/"
const toFilePrefix = "+++ b/"
const newFilePrefix = "--- /dev/null"
const delFilePrefix = "+++ /dev/null"

type FileMode int

const (
	MODIFIED FileMode = iota
	NEW
	DELETED
)

func ParseGitDiff(rawDiff string) (Diff, error) {
	diffLines := strings.Split(rawDiff, "\n")

	var diff Diff
	var file *DiffFile
	var hunk *Hunk

	isFirstFile := true
	inHeader := true

	// Loop over diff
	for _, line := range diffLines {
		switch {
		case strings.HasPrefix(line, "diff "):
			inHeader = true
			if !isFirstFile {
				// Write File
			} else {
				isFirstFile = false
			}

			file = &DiffFile{}
			diff.Files = append(diff.Files, file)
		case strings.HasPrefix(line, fromFilePrefix):
			file.FromFileName = strings.TrimPrefix(line, fromFilePrefix)

		case strings.HasPrefix(line, toFilePrefix):
			file.ToFileName = strings.TrimPrefix(line, toFilePrefix)

		case strings.HasPrefix(line, newFilePrefix):
			file.Mode = NEW

		case strings.HasPrefix(line, delFilePrefix):
			file.Mode = DELETED

		case strings.HasPrefix(line, "@@ "):
			inHeader = false
			hunk = &Hunk{}
			file.Hunks = append(file.Hunks, hunk)
		case !inHeader:
			if line == `\ No newline at end of file` {
				break
			}
			if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "++") {
				hunk.Added += strings.TrimPrefix(line, "+") + "\n"
			} else if strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "--") {
				hunk.Removed += strings.TrimPrefix(line, "-") + "\n"

			}
		}
	}

	return diff, nil
}

type Diff struct {
	Files []*DiffFile
}

// MODE Modified, new, del
type DiffFile struct {
	FromFileName, ToFileName string
	Mode                     FileMode
	Hunks                    []*Hunk
}

type Hunk struct {
	Added, Removed string
}

type SourceLine struct {
	FileFrom, FileTo string
	Content          string
	Position         int
	Mode             LineMode
}

type LineMode int

const (
	ADDED = iota
	REMOVED
)
