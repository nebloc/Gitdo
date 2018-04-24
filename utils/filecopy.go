package utils

import (
	"io/ioutil"
	"os"
	"fmt"
	"path/filepath"
	"io"
)

// AppendFile copies the contents of a file from src and appends to dst.
func AppendFile(src, dst string) error {
	from, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}
	to, err := os.OpenFile(dst, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer to.Close()

	if _, err = to.Write(from); err != nil {
		return err
	}
	return nil
}

// CopyFolder copies a folder from src to dst. It looks through the src folder and copies files one by one to the
// destination folder. It does not copy subdirectories
func CopyFolder(src, dst string) error {
	fmt.Printf("Copying from: %s to %s\n", src, dst)
	files, err := ioutil.ReadDir(src)
	if err != nil {
		return fmt.Errorf("could not get %s files: %v", src, err)
	}
	for _, file := range files {
		sf := filepath.Join(src, file.Name())
		df := filepath.Join(dst, file.Name())
		err = CopyFile(sf, df)
		if err != nil {
			fmt.Printf("could not copy %v - skipping: %v\n", file.Name(), err)
		}
	}

	return nil
}

// CopyFile copies a file from src to dst. If src and dst files exist, and are the same, then return success. Otherise,
// attempt to create a hard link between the two files. If that fail, copy the file contents from src to dst.
func CopyFile(src, dst string) (err error) {
	sfi, err := os.Stat(src)
	if err != nil {
		return
	}
	if !sfi.Mode().IsRegular() {
		// cannot copy non-regular files (e.g., directories,
		// symlinks, devices, etc.)
		return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
	}
	dfi, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
	} else {
		if !(dfi.Mode().IsRegular()) {
			return fmt.Errorf("CopyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
		}
		if os.SameFile(sfi, dfi) {
			return
		}
	}
	if err = os.Link(src, dst); err == nil {
		return
	}
	err = CopyFileContents(src, dst)
	return
}

// CopyFileContents copies the contents of the file named src to the file named by dst. The file will be created if it
// does not already exist. If the destination file exists, all it's contents will be replaced by the contents of the
// source file.
func CopyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}
