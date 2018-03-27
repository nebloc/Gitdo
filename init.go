package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	cli "github.com/urfave/cli"
)

// Init initialises the gitdo project by scaffolding the gitdo folder
func Init(ctx *cli.Context) error {
	if ctx.Bool("with-git") {
		if err := InitGit(); err != nil {
			return err
		}
	}

	cmd := exec.Command("git", "config", "core.hooksPath", "~/Dev/Go/src/github.com/nebbers1111/gitdo/hooks")
	if _, err := cmd.Output(); err != nil {
		log.WithError(err).Error("could not set hooks path")
		return err
	}

	os.Mkdir(".git/gitdo", os.ModePerm)

	cmd = exec.Command("cp", "-r", "/Users/bencoleman/Dev/Go/src/github.com/nebbers1111/gitdo/plugins", ".git/gitdo/")
	if res, err := cmd.CombinedOutput(); err != nil {
		log.Error(stripNewlineChar(res))
		return err
	}
	cmd = exec.Command("cp", "-r", "/Users/bencoleman/Dev/Go/src/github.com/nebbers1111/gitdo/.git/gitdo/secrets.json", ".git/gitdo/")
	if res, err := cmd.CombinedOutput(); err != nil {
		log.Error(stripNewlineChar(res))
		return err
	}
	CheckFolder()
	SetConfig()

	fmt.Println("Done - please check plugins are configured, some need secrets and ID's")
	return nil
}

// SetConfig checks the config is not set and asks the user relevant questions to set it
func SetConfig() {
	err := LoadConfig()
	if err == nil && config.IsSet() {
		return
	}

	if !config.authorIsSet() {
		author, err := AskAuthor()
		if err != nil {
			//TODO: Handle this
			return
		}
		config.Author = author
	}
	if !config.pluginIsSet() {
		plugin, err := AskPlugin()
		if err != nil {
			//TODO: Handle this
			return
		}
		config.Plugin = plugin
	}
	if !config.interpreterIsSet() {
		interp, err := AskInterpreter()
		if err != nil {
			//TODO: Handle this
			return
		}
		config.PluginInterpreter = interp
	}
	err = WriteConfig()
	if err != nil {
		log.WithError(err).Warn("Couldn't save config")
	}
}

// InitGit initialises a git repo before initialising gitdo
func InitGit() error {
	cmd := exec.Command("git", "init")
	if _, err := cmd.Output(); err != nil {
		return err
	}
	return nil
}

// AskAuthor notifies user of email address used
func AskAuthor() (string, error) {
	email, err := getGitEmail()
	if err != nil {
		return "", err
	}
	fmt.Printf("Using %s\n", email)
	return email, nil
}
// AskPlugin reads in plugins from the directory and gives the user a list of plugins, that have a "<name>_getid"
func AskPlugin() (string, error) {
	files, err := ioutil.ReadDir(pluginDir)
	if err != nil {
		return "", err
	}

	fmt.Println("Available plugins:")

	var plugins []string
	i := 0
	for _, f := range files {
		if strings.HasSuffix(f.Name(), getidSuffix) {
			i++
			plugin := strings.TrimSuffix(f.Name(), getidSuffix)
			plugins = append(plugins, plugin)
			fmt.Printf("%d: %s\n", i, plugin)
		}
	}

	fmt.Printf("What plugin would you like to use (1-%d): ", len(plugins))
	var choice string
	_, err = fmt.Scan(&choice)
	if err != nil {
		return "", err
	}
	pN, err := strconv.Atoi(strings.TrimSpace(choice))
	if err != nil || pN > len(plugins) {
		return "", err
	}
	return plugins[pN-1], nil
}

// AskInterpreter asks the user what command they want to use to run the plugin
func AskInterpreter() (string, error) {
	var interp string
	fmt.Printf("What interpreter for this plugin (i.e. python3/node/python): ")
	_, err := fmt.Scan(&interp)
	if err != nil {
		return "", err
	}
	return interp, nil
}
