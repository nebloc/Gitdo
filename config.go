package main

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
)

type Config struct {
	Author     string `json:"author"`
	PluginName string `json:"plugin_name"`
	PluginCmd  string `json:"plugin_cmd"`
	DiffFrom   string `json:"diff_from"`
}

//TODO: load config from file test
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
