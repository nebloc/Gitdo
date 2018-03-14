package diffparse

import (
	"regexp"
	"strconv"
	"strings"
)

const fromFilePrefix = "--- a/"
const toFilePrefix = "+++ b/"
const newFilePrefix = "--- /dev/null"
const delFilePrefix = "+++ /dev/null"

// ParseGitDiff loops over the given diff string and maps it to an array of
// SourceLine structs
func ParseGitDiff(rawDiff string) ([]SourceLine, error) {
	diffLines := strings.Split(rawDiff, "\n")

	isFirstFile := true
	inHeader := true

	var sourceLines []SourceLine

	var fromFileName string
	var toFileName string
	var linePos int

	hunkHeadReg := regexp.MustCompile(`@@ \-(\d+),?(\d+)? \+(\d+),?(\d+)? @@`)

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

		case strings.HasPrefix(line, delFilePrefix):
			toFileName = ""
		case strings.HasPrefix(line, newFilePrefix):
			fromFileName = ""

		case strings.HasPrefix(line, "@@ "):
			inHeader = false
			match := hunkHeadReg.FindStringSubmatch(line)

			newHunkLine, err := strconv.Atoi(match[3])
			if err != nil {
				return nil, err
			}

			linePos = newHunkLine - 1
		case !inHeader:
			linePos++
			if line == `\ No newline at end of file` {
				break
			}
			if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "++") {
				l := SourceLine{
					fromFileName,
					toFileName,
					strings.TrimPrefix(line, "+"),
					linePos,
					ADDED,
				}
				sourceLines = append(sourceLines, l)
			} else if strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "--") {
				l := SourceLine{
					fromFileName,
					toFileName,
					strings.TrimPrefix(line, "-"),
					linePos,
					REMOVED,
				}
				linePos--
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
	ADDED LineMode = iota
	REMOVED
)
