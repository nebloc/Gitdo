package versioncontrol

import (
	"path"
	"os"
	"testing"
	"fmt"
	"path/filepath"
)

func (tvc *TestVC) moveToDir(t *testing.T) {
	t.Log("Here")
	err := os.Chdir(tvc.tmpDir)
	if err != nil {
		t.Fatalf("Failed to move to %s", tvc.tmpDir)
	}
	path, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current path: %v", err)
	}
	t.Logf("Currently in dir: %s", path)
}

var HomeDir string

type TestVC struct {
	VersionControl
	tmpDir string
}

var VCMap map[string]*TestVC

func init() {
	dir, _ := os.Getwd()
	HomeDir = filepath.Dir(filepath.Dir(dir))

	VCMap = make(map[string]*TestVC)
	VCMap[GIT_NAME] = &TestVC{NewGit(), ""}
	VCMap[MERCURIAL_NAME] = &TestVC{NewHg(), ""}

	for i, vc := range VCMap {
		VCMap[i].tmpDir = path.Join(os.TempDir(), "Gitdo_versioncontrol_"+vc.NameOfVC())

		_ = os.RemoveAll(vc.tmpDir)
		err := os.Mkdir(vc.tmpDir, os.ModePerm)
		if err != nil {
			panic("could not create test dir for " + vc.NameOfVC())
		}
		err = os.Chdir(vc.tmpDir)
		if err != nil {
			panic("could not move to new dir for " + vc.NameOfVC())
		}

		err = vc.Init()
		if err != nil {
			panic("Init function failed for " + vc.NameOfVC())
		}
	}

	for name, val := range VCMap {
		fmt.Printf("%s: %s\n", name, val.tmpDir)
	}

}
