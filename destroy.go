package main

import (
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	cli "github.com/urfave/cli"
)

func Destroy(ctx *cli.Context) error {
	if !ctx.Bool("yes") {
		return nil
	}
	return os.Remove(StagedTasksFile)
}

func ConfirmUser(ctx *cli.Context) error {
	if ctx.Bool("yes") {
		return nil
	}
	var ans string
	fmt.Print("Are you sure you want to purge the task file? (y/n)")
	_, err := fmt.Scan(&ans)
	if err != nil {
		log.WithError(err).Error("not purging")
		return nil
	}
	ans = strings.TrimSpace(ans)
	ans = strings.ToLower(ans)

	if ans == "y" || ans == "yes" {
		ctx.Set("yes", "true")
		return nil
	}

	log.Info("Not Purging")
	return nil
}
