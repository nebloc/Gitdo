package cmd

import (
	"os"
	"testing"

	"github.com/nebloc/gitdo/versioncontrol"
)

func init() {
	app = &config{
		vc:                versioncontrol.NewGit(),
		Author:            "benjamin.coleman@me.com",
		Plugin:            "Test",
		PluginInterpreter: "python",
	}
}

// setupForTest creates a directory in the os.TMPDIR and moves in to it, sets the configuration and creates a new app
// context. Returns the app context and a function to chdir back.
func setupForTest(t *testing.T) func() {
	cDir, closeDir := testDirHelper(t)
	t.Logf("working in dir: %s", cDir)

	setVCPaths()

	return closeDir
}

// testCommitHelper creates a new directory and moves in to it, returning a close function to be called to move back to
// the original dir
func testDirHelper(t *testing.T) (string, func()) {
	t.Helper()
	origPath, err := os.Getwd()
	if err != nil {
		t.Fatal("couldn't get current path")
	}
	dirPath := os.TempDir() + "gitdotest/"

	if err = os.RemoveAll(dirPath); err != nil {
		t.Fatal("could not remove temp dir")
	}

	err = os.Mkdir(dirPath, os.ModePerm)
	if err != nil {
		t.Fatal("couldn't create temp dir")
	}
	err = os.Chdir(dirPath)
	if err != nil {
		t.Fatalf("couldn't move in to temp dir: %v", err)
	}
	return dirPath, func() {
		err := os.Chdir(origPath)
		if err != nil {
			t.Fatal("couldn't change back to default dir")
		}
		// Not deleting the temp dir currently, as the OS will eventually, however may want to for different VC tests
	}
}
