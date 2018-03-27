package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
	"strings"
)

const ConfigFilePath = GitdoDir + "config.json"

type Config struct {
	// Author to attach to task in task manager.
	Author string `json:"author"`
	// Plugin to use at push time
	Plugin string `json:"plugin_name"`
	// The command to run for plugin files
	PluginInterpreter string `plugin_interpreter`

	// Example of plugin: "test" and plugin_interpreter: "python"
	// Will run 'python .git/gitdo/plugins/reserve_test'
}

// String returns a human readable format of the Config struct
func (c *Config) String() string {
	return fmt.Sprintf("Author: %s\nPlugin: %s", c.Author, c.Plugin)
}

// Checks that the configuration has all the information needed
func (c *Config) IsSet() bool {
	if !c.pluginIsSet() {
		return false
	}
	if !c.authorIsSet() {
		return false
	}
	if !c.interpreterIsSet() {
		return false
	}
	return true
}

// pluginIsSet returns if the plugin in config is not empty
func (c *Config) pluginIsSet() bool {
	return strings.TrimSpace(c.Plugin) != ""
}

// authorIsSet returns if the author in config is not empty
func (c *Config) authorIsSet() bool {
	return strings.TrimSpace(c.Author) != ""
}

// interpreterIsSet returns if the plugin interpreter in config is not empty
func (c *Config) interpreterIsSet() bool {
	return strings.TrimSpace(c.PluginInterpreter) != ""
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

// getGitEmail runs the 'git config user.email' command to get the default email address of the user for the current dir
func getGitEmail() (string, error) {
	cmd := exec.Command("git", "config", "user.email")
	resp, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return stripNewlineChar(resp), nil
}

// WriteConfig saves the current config to be loaded in after setting
func WriteConfig() error {
	bConf, err := json.MarshalIndent(config, "", "\t")
	if err != nil {
		return err
	}
	ioutil.WriteFile(ConfigFilePath, bConf, os.ModePerm)
	return nil
}
