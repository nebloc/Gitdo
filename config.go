package main

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
)

type Config struct {
	// Author to attach to task in task manager.
	Author string `json:"author"`
	// Plugin to use at push time
	PluginName string `json:"plugin_name"`
	// Plugin interpreter to use
	PluginCmd string `json:"plugin_cmd"`
	// Where to load the diff from. Currently for debugging only.
	DiffFrom string `json:"diff_from"`
}

// LoadConfig opens a configuration file and reads it in to the Config struct
func LoadConfig() error {
	bConfig, err := ioutil.ReadFile(".git/gitdo/config.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(bConfig, config)
	if err != nil {
		if _, ok := err.(*os.PathError); ok {
			log.Warn("No config.json file, using defaults")
		}
		return err
	}
	log.WithFields(log.Fields{
		"author":             config.Author,
		"plugin_name":        config.PluginName,
		"plugin_interpreter": config.PluginCmd,
		"diff_source":        config.DiffFrom,
	}).Debug("Config")
	return nil
}

func LoadDefaultConfig() error {
	email, err := getGitEmail()
	if err != nil {
		return fmt.Errorf("could not get default git email: %s", err)
	}
	config = &Config{
		email,
		"",
		"",
		"cmd",
	}
	return nil
}

func getGitEmail() (string, error) {
	cmd := exec.Command("git", "config", "user.email")
	resp, err := cmd.Output()
	if err != nil {
		return "", err
	}
	email := fmt.Sprintf("%s", resp)
	email = email[:len(email)-1]
	return email, nil
}
