package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/nebloc/gitdo/diffparse"
	"github.com/urfave/cli"
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
		VC:                GIT,
		Author:            "benjamin.coleman@me.com",
		Plugin:            "Test",
		PluginInterpreter: "python3",
	}
	SetVCType()

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
	if err != errNotGitDir {
		t.Errorf("Expected: %v, got: %v", errNotGitDir, err)
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
	tasks := testTaskFileCorrect(t)
	if len(tasks.NewTasks) != 1 {
		t.Errorf("Should have 1 task in the new task area. Have: %d", len(tasks.NewTasks))
	}

	testCommitHelper(t)
	testDeleteTaskCommentHelper(t, fileName)

	err = Commit(ctx)
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}
	tasks = testTaskFileCorrect(t)
	if len(tasks.DoneTasks) != 1 {
		t.Errorf("Should have 1 task in the done task area. Have: %d", len(tasks.DoneTasks))
	}
}

func testTaskFileCorrect(t *testing.T) *Tasks {
	tasks, err := getTasksFile()
	if err != nil {
		t.Errorf("Could not load tasks file to check: %v", err)
	}
	return tasks
}

func testDeleteTaskCommentHelper(t *testing.T, fileName string) {
	fileCont, err := ioutil.ReadFile(fileName)
	if err != nil {
		t.Errorf("Could not get mock file: %v", err)
	}
	newCont := strings.Replace(string(fileCont), "// TODO: Test <1234>\n", "", 1)

	err = ioutil.WriteFile(fileName, []byte(newCont), os.ModePerm)
	if err != nil {
		t.Errorf("Could not write new file with todo removed: %v", err)
	}
	testAddToGitHelper(t, fileName)
}

func TestCheckTagged(t *testing.T) {
	t.Log(config.String())
	line := diffparse.SourceLine{
		FileFrom: "main.go",
		FileTo:   "main.go",
		Content:  "",
		Position: 32,
		Mode:     diffparse.ADDED,
	}
	testData := []struct {
		LineContent string
		ExpFound    bool
		ExpID       string
	}{
		{"//TODO: Hello <08238>", true, "08238"},
		{"//TODO:Hello <08238>", true, "08238"},
		{"// TODO: Hello <08238>", true, "08238"},
		{"// TODO:Hello <08238>", true, "08238"},
		{"//TODO: Hello", false, ""},
		{"+// TODO: Test <fhsiufh>", false, ""},
		{"#TODO: Hello <08238>", true, "08238"},
		{"#TODO:Hello <08238>", true, "08238"},
		{"# TODO: Hello <08238>", true, "08238"},
		{"# TODO:Hello <08238>", true, "08238"},
		{"#TODO: Hello", false, ""},
		{"+# TODO: Test <fhsiufh>", false, ""},
	}

	for _, data := range testData {
		line.Content = data.LineContent
		id, found := CheckTagged(line)
		if found != data.ExpFound {
			t.Errorf("Line: %s\nExpected: %v, Got: %v", data.LineContent, data.ExpFound, found)
			continue
		}
		if id != data.ExpID {
			t.Errorf("Line: %s\nExpected: %v, Got: %v", data.LineContent, data.ExpID, id)
			continue
		}
	}
}

func TestCheckTaskRegex(t *testing.T) {
	t.Log(config.String())
	testData := []struct {
		LineContent string
		ExpTask     string
	}{
		{"//TODO: Hello", "Hello"},
		{"//TODO:Hello", "Hello"},
		{"// TODO: Hello", "Hello"},
		{"// TODO:Hello", "Hello"},
		{"+// TODO: Hello", ""},
		{"#TODO: Hello", "Hello"},
		{"#TODO:Hello", "Hello"},
		{"# TODO: Hello", "Hello"},
		{"+# TODO: Hello", ""},
	}

	for _, data := range testData {
		match := CheckTaskRegex(data.LineContent)
		if len(match) == 0 && data.ExpTask != "" {
			t.Errorf("Expected to not match: %v", data)
		} else if len(match) != 0 && match[1] != data.ExpTask {
			t.Errorf("Expected: %s, Got: %s", data.ExpTask, match[1])
		}
	}
}

func TestGetDiffFromCmd(t *testing.T) {
	_, closeDir := setupForTest(t)
	defer closeDir()
	t.Log(config.String())

	_, err := GetDiffFromGit()
	if err != errNotGitDir {
		t.Errorf("Expected not a git repo, got: %v", err)
	}
	testStartRepoHelper(t)
	diff, err := GetDiffFromGit()
	if err != errNoDiff {
		t.Errorf("expected diff to be empty: %v: %v", err, diff)
	}
	file := testMockFileHelper(t)
	testAddToGitHelper(t, file)
	diff, err = GetDiffFromGit()
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
	t.Logf("Added file: %s", fileName)
}

func testCommitHelper(t *testing.T) {
	cmd := exec.Command("git", "commit", "-am", "new file")
	err := cmd.Run()
	if err != nil {
		t.Fatalf("error commiting file: %v", err)
	}
	t.Log("Committed staged files")
}

// testStartRepoHelper runs git init and adds the gitdo folder, which is otherwise done in the Setup() func
func testStartRepoHelper(t *testing.T) {
	t.Helper()
	t.Log("Started repo")
	cmd := exec.Command("git", "init")
	if err := cmd.Run(); err != nil {
		t.Fatal("could not create repo")
	}
	if err := os.MkdirAll(".git/gitdo/plugins/Test", os.ModePerm); err != nil {
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
+// TODO: Test
+func main(){
+	fmt.Println("Hello Ben")
+}
\ No newline at end of file`
