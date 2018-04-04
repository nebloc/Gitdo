package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"

	"github.com/urfave/cli"
	"github.com/nebloc/gitdo/diffparse"
)

func TestRegexs(t *testing.T) {
	testStrings := []string{
		"//TODO: Test that taggedIds with 13 works <22031810224912>",
		"//TODO: Test that taggedIds with 14 works <125861826182573>",
		"//TODO: Test that taggedIds with email works <benjamin.coleman@me.com:14861826182573>",
	}
	for _, str := range testStrings {
		m := taggedReg.MatchString(str)
		if !m {
			t.Errorf("Expected match for tagged line: %s", str)
		}
	}
}

func setupForTest(t *testing.T) (*cli.Context, func()) {
	config = &Config{
		Author:            "benjamin.coleman@me.com",
		Plugin:            "test",
		PluginInterpreter: "python",
	}

	cDir, closeDir := testDirHelper(t)
	t.Logf("working in dir: %s", cDir)
	ctx := cli.NewContext(gitdo, nil, nil)

	return ctx, closeDir
}

func TestCommit(t *testing.T) {
	ctx, closeDir := setupForTest(t)
	defer closeDir()

	t.Log(config.String())

	err := Commit(ctx)
	if err != ErrNotGitDir {
		t.Errorf("Expected: %v, got: %v", ErrNotGitDir, err)
	}

	testStartRepoHelper(t)

	err = Commit(ctx)
	if err != nil {
		t.Errorf("Expected commit to return with no error, got: %v", err)
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

func TestCheckTagged(t *testing.T) {
	t.Log(config.String())
	line := diffparse.SourceLine{"main.go", "main.go", "", 32, diffparse.ADDED}
	testData := []struct {
		LineContent string
		ExpFound    bool
		ExpID       string
	}{
		{"//TODO: Hello <08238>", true, "08238"},
		{"//TODO: Hello", false, ""}, <B6UbQF7D>
		{"+// TODO: Test <fhsiufh>", false, ""},
	}

	for _, data := range testData {
		line.Content = data.LineContent
		id, found := CheckTagged(line)
		if found != data.ExpFound {
			t.Errorf("Line: %s\nExpected: %v, Got: %v", data.LineContent, data.ExpFound, found)
		}
		if id != data.ExpID {
			t.Errorf("Line: %s\nExpected: %v, Got: %v", data.LineContent, data.ExpID, id)
		}
	}
}

func TestGetDiffFromCmd(t *testing.T) {
	_, closeDir := setupForTest(t)
	defer closeDir()
	t.Log(config.String())

	_, err := GetDiffFromCmd()
	if err != ErrNotGitDir {
		t.Errorf("Expected not a git repo, got: %v", err)
	}
	testStartRepoHelper(t)
	diff, err := GetDiffFromCmd()
	if err != ErrNoDiff {
		t.Errorf("expected diff to be empty: %v: %v", err, diff)
	}
	file := testMockFileHelper(t)
	testAddToGitHelper(t, file)
	diff, err = GetDiffFromCmd()
	if err != nil {
		t.Errorf("Unexpected error getting diff: %v", err)
	}
	if diff != mockDiffExample {
		t.Errorf("Expected diff to be \n%v\nGot: \n%v", mockDiffExample, diff)
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
	if err := os.MkdirAll(".git/gitdo/plugins/test", os.ModePerm); err != nil {
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

// TODO: Test <1234>
func main(){
	fmt.Println("Hello Ben")
}`

const mockDiffExample string = `diff --git a/main.go b/main.go
new file mode 100755
index 0000000..a30278c
--- /dev/null
+++ b/main.go
@@ -0,0 +1,8 @@
+package main
+
+import "fmt"
+
+// TODO: Test <9ypvkCD1>
+func main(){
+	fmt.Println("Hello Ben")
+}
\ No newline at end of file`
