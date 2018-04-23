/**
Not really maintained - not a very useful command to users in it's current state, should make it uninstall all from the
project in future
*/
package main

import (
	"fmt"
	"os"
	"strings"
	"github.com/nebloc/gitdo/app/utils"
	"github.com/urfave/cli"
)

// Destroy deletes the staged tasks file if you need it to
func Destroy(ctx *cli.Context) error {
	if !ctx.Bool("yes") {
		return nil
	}
	return os.Remove(stagedTasksFile)
}

// CheckPurge asks if the user really wants to delete the file, if yes it sets the yes flag
func CheckPurge(ctx *cli.Context) error {
	if ctx.Bool("yes") {
		return nil
	}

	confirmed := ConfirmWithUser("Are you sure you want to purge the task file?")
	if confirmed {
		ctx.Set("yes", "true")
		return nil
	}

	return nil
}

func ConfirmWithUser(message string) bool {
	var ans string
	utils.Warnf("%s (y/n)", message)
	_, err := fmt.Scan(&ans)
	if err != nil {
		utils.Warnf("Backing out: %v", err)
		return false
	}
	ans = strings.TrimSpace(ans)
	ans = strings.ToLower(ans)

	if ans == "y" || ans == "yes" {
		return true
	}

	utils.Warn("Backing out")
	return false
}
