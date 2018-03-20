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
	PluginName string `json:"plugin_name"`
	// Plugin interpreter to use
	PluginCmd string `json:"plugin_cmd"`
	// Where to load the diff from. Currently for debugging only.
	DiffFrom string `json:"diff_from"`
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
	return stripNewlineChar(resp), nil
}

func stripNewlineChar(orig []byte) string {
	new := string(orig)
	new = new[:len(new)-1]
	return new
}

func WriteConfig() error {
	bConf, err := json.MarshalIndent(config, "", "\t")
	if err != nil {
		return err
	}
	ioutil.WriteFile(ConfigFilePath, bConf, os.ModePerm)
	return nil
}
