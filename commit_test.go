package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"

	cli "github.com/urfave/cli"
)

func TestCommit(t *testing.T) {
	cDir, closeDir := testDirHelper(t)
	defer closeDir()
	t.Logf("working in dir: %s", cDir)

	config = &Config{
		Author:     "benjamin.coleman@me.com",
		PluginName: "",
		PluginCmd:  "",
		DiffFrom:   "cmd",
	}

	ctx := cli.NewContext(gitdo, nil, nil)

	err := Commit(ctx)
	if err != ErrNotGitDir {
		t.Errorf("Expected: %v, got: %v", ErrNotGitDir, err)
	}

	testStartRepoHelper(t)

	err = Commit(ctx)
	if err != ErrNoDiff {
		t.Errorf("Expected: %v, got: %v", ErrNoDiff, err)
	}

	fileName := testMockFileHelper(t)

	err = Commit(ctx)
	if err != nil {
		t.Errorf("didn't expect an error: %v", err)
	}

	bMock, err := ioutil.ReadFile(fileName)
	if err != nil {
		t.Fatalf("couldn't read in file to against golden: %v", err)
	}
	if bytes.Compare([]byte(goldenFileContent), bMock) != 0 {
		t.Errorf("expected:\n%s \n\ngot:\n%s", goldenFileContent, bMock)
	}

}

// testMockFileHelper creates a file that has a TODO comment.
func testMockFileHelper(t *testing.T) string {
	t.Helper()
	fileName := "main.go"
	err := ioutil.WriteFile(fileName, []byte(mockFileContent), os.ModePerm)
	if err != nil {
		t.Fatal("could not create mock file")
	}
	testAddToGitHelper(t, fileName)
	return fileName
}

// testAddToGitHelper runs 'git add $fileName'
func testAddToGitHelper(t *testing.T, fileName string) {
	t.Helper()
	cmd := exec.Command("git", "add", fileName)
	err := cmd.Run()
	if err != nil {
		t.Fatalf("could not 'git add %s'", fileName)
	}
	cmd = exec.Command("git", "status")
	resp, err := cmd.Output()
	if err != nil {
		t.Fatal("could not 'git status'")
	}
	t.Logf("status:\n%s", resp)
}

// testStartRepoHelper runs git init and adds the gitdo folder, which is otherwise done in the Setup() func
func testStartRepoHelper(t *testing.T) {
	t.Helper()
	cmd := exec.Command("git", "init")
	if err := cmd.Run(); err != nil {
		t.Fatal("could not create repo")
	}
	if err := os.Mkdir(".git/gitdo", os.ModePerm); err != nil {
		t.Fatal("could not create gitdo folder")
	}
}

// testCommitHelper creates a new directory and moves in to it, returning a close function to be called to move back to the original dir
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
	}
}

const mockFileContent string = `package main

import "fmt"

// TODO: Test
func main(){
	fmt.Println("Hello Ben")
}`

const goldenFileContent string = `package main

import "fmt"

// TODO: Test <GITDO>
func main(){
	fmt.Println("Hello Ben")
}`
