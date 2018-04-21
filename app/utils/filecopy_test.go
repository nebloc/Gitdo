package utils

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func init() {
	dir := filepath.Join(os.TempDir(), "gitdofilecopy")
	err := os.Mkdir(dir, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		panic("could not create directory for testing")
	}
	err = os.Chdir(dir)
	if err != nil {
		panic("could not move to a directory for testing")
	}
}

func TestAppendFile(t *testing.T) {

}

func TestCopyFile(t *testing.T) {
	firstFile := "test.txt"
	secondFile := "test2.txt"

	testData := []byte("Foo\n")

	err := ioutil.WriteFile(firstFile, testData, os.ModePerm)
	if err != nil {
		t.Errorf("could not create file for testing: %v", err)
	}

	err = CopyFile(firstFile, secondFile)
	if err != nil {
		t.Errorf("failed to copy file: %v", err)
	}

	result, err := ioutil.ReadFile(secondFile)
	if err != nil {
		t.Errorf("could not read new file in: %v", err)
	}

	for i, bt := range result {
		if bt != testData[i] {
			t.Errorf("file not the same as original: \n%v\nnew:\n%v", testData, result)
			break
		}
	}
}
