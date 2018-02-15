package main

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Author     string `json:"author"`
	PluginFile string `json:"plugin_file"`
	PluginCmd  string `json:"plugin_cmd"`
}

//TODO: load config from file test
func LoadConfig() error {
	bConfig, err := ioutil.ReadFile("config.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(bConfig, &config)
	if err != nil {
		return err
	}
	return nil
}
