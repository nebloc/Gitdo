package gitparse

import (
	"testing"
)

func TestGetGitDiff(t *testing.T) {
	diff, err := GetGitDiff()
	if err != nil {
		t.Errorf("GetGitDiff returned error: %v", err)
	}
	t.Logf("GetGitDiff returned: %s", diff)
}

func TestParseGitDiff(t *testing.T) {
	diff, _ := GetGitDiff()
	diffFiles, err := ParseGitDiff(diff)
	if err != nil {
		t.Errorf("parse diff returned error: %v", err)
	}
	t.Logf("parse returned: %v", diffFiles.Files[0].Hunks[0].Added)
}
