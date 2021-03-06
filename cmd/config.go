package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/nebloc/gitdo/versioncontrol"
)

type config struct {
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
func (c *config) String() string {
	return fmt.Sprintf(
		"Author: %s\nPlugin: %s\nInterpreter: %s",
		c.Author, c.Plugin, c.PluginInterpreter)
}

// Checks that the configuration has all the information needed
func (c *config) IsSet() bool {
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
func (c *config) pluginIsSet() bool {
	return strings.TrimSpace(c.Plugin) != ""
}

// authorIsSet returns if the author in config is not empty
func (c *config) authorIsSet() bool {
	return strings.TrimSpace(c.Author) != ""
}

// interpreterIsSet returns if the plugin interpreter in config is not empty
func (c *config) interpreterIsSet() bool {
	return strings.TrimSpace(c.PluginInterpreter) != ""
}

// loadConfig opens a configuration file and reads it in to the Config struct
func loadConfig() error {
	bConfig, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return fmt.Errorf("could not read in configuration file: %v", err)
	}

	err = json.Unmarshal(bConfig, app)
	if err != nil {
		return err
	}

	return nil
}

// writeConfig saves the current config to be loaded in after setting
func writeConfig() error {
	bConf, err := json.MarshalIndent(app, "", "\t")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(configFilePath, bConf, os.ModePerm)
	return err
}
