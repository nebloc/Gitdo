// +build windows

package main

import (
	"errors"
	"os/user"
	"path/filepath"
)

// GetHomeDir gets the home directory of the current user, and gets the Gitdo folder in the AppData directory - windows.
func GetHomeDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", errors.New("Could not determine user and therefore gitdo install directory")
	}
	gitdoPath := filepath.Join(usr.HomeDir, "AppData", "roaming", "Gitdo")
	return gitdoPath, nil
}
