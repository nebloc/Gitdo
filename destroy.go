package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli"
)

// Destroy deletes the staged tasks file if you need it to
func Destroy(ctx *cli.Context) error {
	if !ctx.Bool("yes") {
		return nil
	}
	return os.Remove(StagedTasksFile)
}

// ConfirmUser asks if the user really wants to delete the file, if yes it sets the yes flag
func ConfirmUser(ctx *cli.Context) error {
	if ctx.Bool("yes") {
		return nil
	}
	var ans string
	Warn("Are you sure you want to purge the task file? (y/n)")
	_, err := fmt.Scan(&ans)
	if err != nil {
		Warnf("Not purging: %v", err)
		return nil
	}
	ans = strings.TrimSpace(ans)
	ans = strings.ToLower(ans)

	if ans == "y" || ans == "yes" {
		ctx.Set("yes", "true")
		return nil
	}

	Warn("Not Purging")
	return nil
}
