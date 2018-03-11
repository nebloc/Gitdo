package main

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
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
	bConfig, err := ioutil.ReadFile("./config.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(bConfig, &config)
	if err != nil {
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
