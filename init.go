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

	if err := os.Mkdir(".git/gitdo", os.ModePerm); err != nil {
		log.WithError(err).Error("could not create gitdo folder")
		return err
	}

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

	return nil
}

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
	err = WriteConfig()
	if err != nil {
		log.WithError(err).Warn("Couldn't save config")
	}
}

func InitGit() error {
	cmd := exec.Command("git", "init")
	if _, err := cmd.Output(); err != nil {
		return err
	}
	return nil
}

func AskAuthor() (string, error) {
	email, err := getGitEmail()
	if err != nil {
		return "", err
	}
	fmt.Printf("Using %s\n", email)
	return email, nil
}

func AskPlugin() (string, error) {
	files, err := ioutil.ReadDir(".git/gitdo/plugins")
	if err != nil {
		return "", err
	}

	fmt.Println("What plugin would you like to use:")

	var plugins []string
	prefix := "reserve_"
	i := 0
	for _, f := range files {
		if strings.HasPrefix(f.Name(), prefix) {
			i++
			plugin := strings.TrimPrefix(f.Name(), prefix)
			plugins = append(plugins, plugin)
			fmt.Printf("%d: %s\n", i, plugin)
		}
	}

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
