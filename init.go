package main

import (
	"os/exec"

	log "github.com/sirupsen/logrus"
	cli "github.com/urfave/cli"
)

func Init(ctx *cli.Context) error {
	cmd := exec.Command("git", "config", "core.hooksPath", "~/Dev/Go/src/github.com/nebbers1111/gitdo/hooks")
	if _, err := cmd.Output(); err != nil {
		log.WithError(err).Error("could not set hooks path")
		return err
	}

	return nil
}

func InitGit(ctx *cli.Context) error {
	cmd := exec.Command("git", "init")
	if _, err := cmd.Output(); err != nil {
		return err
	}
	return nil
}
