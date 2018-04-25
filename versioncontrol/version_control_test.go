package versioncontrol

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"testing"
)

func (tvc *TestVC) moveToDir(t *testing.T) {
	t.Helper()
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
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Could not get current directory: %v", err)
	}
	HomeDir = filepath.Join(filepath.Dir(dir), "resources")
	log.Print(HomeDir)

	VCMap = make(map[string]*TestVC)
	VCMap[GitName] = &TestVC{NewGit(), ""}
	VCMap[MercurialName] = &TestVC{NewHg(), ""}

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
		log.Printf("%s: %s\n", name, val.tmpDir)
	}

}
