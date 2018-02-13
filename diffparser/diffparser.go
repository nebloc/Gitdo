package diffparser

import (
	"strings"
)

// TODO: Do I need to develop my own diffparser?
// SplitToFiles takes the given diff and splits it in to an array of strings, each a different file section of the diff
func SplitToFiles(diff string) []string {
	diffLines := strings.Split(diff, "\n")
	var files []string

	currentFile := ""
	for _, line := range diffLines {
		if strings.HasPrefix(line, "diff") && currentFile != "" {
			files = append(files, currentFile)
			currentFile = ""
		}
		currentFile += line + "\n"
	}
	files = append(files, currentFile)

	return files
}

// ParseDiff takes a given string diff and converts it in to useable structs
func ParseDiff(diff string) {

}

type File struct {
	OrigFile string
	NewFile  string
	Hunks    []hunk
}
type hunk struct {
	origStartInd int
	origRange    int
	newStartInd  int
	newRange     int
	content      string
}

func SplitToChunk(file string) {}
