// +build !windows

package main

import (
	"errors"
	"os/user"
	"path/filepath"
)

func GetTemplateDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", errors.New("Could not determine user and therefore gitdo install directory")
	}
	gitdoPath := filepath.Join(usr.HomeDir, ".gitdo")
	return gitdoPath, nil
}
