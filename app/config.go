package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"github.com/urfave/cli"
	"github.com/nebloc/gitdo/app/versioncontrol"
	"github.com/nebloc/gitdo/app/utils"
	"strings"
)

type Config struct {
	vc versioncontrol.VersionControl
	// Author to attach to task in task manager.
	Author string `json:"author"`
	// Plugin to use at push time
	Plugin string `json:"plugin_name"`
	// The command to run for plugin files
	PluginInterpreter string `json:"plugin_interpreter"`

	// Example of plugin: "test" and plugin_interpreter: "python"
	// Will run 'python .git/gitdo/plugins/reserve_test'
}

// String returns a human readable format of the Config struct
func (c *Config) String() string {
	return fmt.Sprintf(
		"Author: %s\nPlugin: %s\nInterpreter: %s",
		c.Author, c.Plugin, c.PluginInterpreter)
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
func LoadConfig(_ *cli.Context) error {
	bConfig, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		utils.Warn("Could not find configuration file for gitdo. Have you ran \"gitdo init\"?")
		return err
	}

	err = json.Unmarshal(bConfig, config)
	if err != nil {
		return err
	}

	return nil
}

// WriteConfig saves the current config to be loaded in after setting
func WriteConfig() error {
	bConf, err := json.MarshalIndent(config, "", "\t")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(configFilePath, bConf, os.ModePerm)
	return err
}
