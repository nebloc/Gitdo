package versioncontrol

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const GIT_NAME string = "Git"

func TestGit_NameOfDir(t *testing.T) {
	result := VCMap[GIT_NAME].NameOfDir()
	expected := ".git"
	if result != expected {
		t.Errorf("Expected NameOfDir to return %s, got %s", expected, result)
	}
}

func TestGit_NameOfVC(t *testing.T) {
	result := VCMap[GIT_NAME].NameOfVC()
	expected := "Git"
	if result != expected {
		t.Errorf("Expected NameOfVC to return %s, got %s", expected, result)
	}
}

func TestGit_GetDiff(t *testing.T) {
	VCMap[GIT_NAME].moveToDir(t)

	fileName := "new.txt"
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		t.Fatalf("failed to open a new file: %v", err)
	}

	_, err = file.Write([]byte("test string"))
	if err != nil {
		t.Fatalf("failed to write to new file: %v", err)
	}

	cmd := exec.Command("git", "add", fileName)
	err = cmd.Run()
	if err != nil {
		t.Fatalf("failed to add %s to git: %v", fileName, err)
	}

	diff, err := VCMap[GIT_NAME].GetDiff()
	if err != nil {
		t.Errorf("didn't expect an error in GetDiff: %v", err)
	}
	if diff != expectedDiff {
		t.Logf("Expected:\n%s\n\nGot:\n%s\n", expectedDiff, diff)
	}
}

var expectedDiff = `diff --git a/new.txt b/new.txt
new file mode 100755
index 0000000..f500b14
--- /dev/null
+++ b/new.txt
@@ -0,0 +1 @@
+test string
\ No newline at end of file`

func TestGit_SetHooks(t *testing.T) {
	Hooks := []string{"pre-commit", "post-commit", "pre-push"}

	VCMap[GIT_NAME].moveToDir(t)
	err := VCMap[GIT_NAME].SetHooks(HomeDir)
	if err != nil {
		t.Errorf("Didn't expect error: %v", err)
	}
	for _, fileName := range Hooks {
		filePath := filepath.Join(VCMap[GIT_NAME].NameOfDir(), "hooks", fileName)

		fileCont, err := ioutil.ReadFile(filePath)
		if err != nil {
			t.Errorf("couldn't read new %s: %v", filePath, err)
		}

		strings.Contains(string(fileCont), "gitdo")
	}

}
