// +build windows

package main

import (
	"errors"
	"fmt"
	"os/user"
	"path/filepath"
)

func GetTemplateDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", errors.New("Could not determine user and therefore gitdo install directory")
	}
	gitdoPath := filepath.Join(usr.HomeDir, "AppData", "roaming", "Gitdo")
	return gitdoPath, nil
}
