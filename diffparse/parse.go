package diffparse

import (
	"strings"
)

const fromFilePrefix = "--- a/"
const toFilePrefix = "+++ b/"
const newFilePrefix = "--- /dev/null"
const delFilePrefix = "+++ /dev/null"

type FileMode int

const (
	// Type of change to file in git diff
	MODIFIED FileMode = iota // File contains a change
	NEW                      // File is new to git
	DELETED                  // File has been deleted
)

// ParseGitDiff loops over the given diff string and maps it to an array of
// SourceLine structs
func ParseGitDiff(rawDiff string) ([]SourceLine, error) {
	diffLines := strings.Split(rawDiff, "\n")

	isFirstFile := true
	inHeader := true

	var sourceLines []SourceLine

	var fromFileName string
	var toFileName string

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
		case strings.HasPrefix(line, fromFilePrefix):
			fromFileName = strings.TrimPrefix(line, fromFilePrefix)

		case strings.HasPrefix(line, toFilePrefix):
			toFileName = strings.TrimPrefix(line, toFilePrefix)

		case strings.HasPrefix(line, "@@ "):
			inHeader = false

		case !inHeader:
			if line == `\ No newline at end of file` {
				break
			}
			if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "++") {
				l := SourceLine{
					fromFileName,
					toFileName,
					strings.TrimPrefix(line, "+"),
					0,
					ADDED,
				}
				sourceLines = append(sourceLines, l)
			} else if strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "--") {
				l := SourceLine{
					fromFileName,
					toFileName,
					strings.TrimPrefix(line, "-"),
					0,
					REMOVED,
				}
				sourceLines = append(sourceLines, l)
			}
		}
	}

	return sourceLines, nil
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
