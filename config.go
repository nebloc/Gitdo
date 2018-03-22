package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

const ConfigFilePath = GitdoDir + "config.json"

type Config struct {
	// Author to attach to task in task manager.
	Author string `json:"author"`
	// Plugin to use at push time
	Plugin string `json:"plugin_name"`
}

func (c *Config) String() string {
	return fmt.Sprintf("Author: %s\nPlugin: %s", c.Author, c.Plugin)
}

func (c *Config) IsSet() bool {
	if !c.pluginIsSet() {
		return false
	}
	if !c.authorIsSet() {
		return false
	}
	return true
}

func (c *Config) pluginIsSet() bool {
	return c.Plugin != ""
}
func (c *Config) authorIsSet() bool {
	return c.Author != ""
}

// LoadConfig opens a configuration file and reads it in to the Config struct
func LoadConfig() error {
	bConfig, err := ioutil.ReadFile(ConfigFilePath)
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
		"author": config.Author,
		"plugin": config.Plugin,
	}).Debug("Config")
	return nil
}

func getGitEmail() (string, error) {
	cmd := exec.Command("git", "config", "user.email")
	resp, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return stripNewlineChar(resp), nil
}

func WriteConfig() error {
	bConf, err := json.MarshalIndent(config, "", "\t")
	if err != nil {
		return err
	}
	ioutil.WriteFile(ConfigFilePath, bConf, os.ModePerm)
	return nil
}
