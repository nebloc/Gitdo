/**
Not really maintained - not a very useful command to users in it's current state, should make it uninstall all from the
project in future
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/nebloc/gitdo/utils"
	"github.com/spf13/cobra"
)

// Destroy deletes the staged tasks file if you need it to
func Destroy(cmd *cobra.Command, args []string) error {
	if ConfirmUser(cmd, args) {
		return os.Remove(stagedTasksFile)
	}
	return nil
}

// ConfirmUser asks if the user really wants to delete the file, if yes it sets the yes flag
func ConfirmUser(cmd *cobra.Command, args []string) bool {
	var ans string
	utils.Warn("Are you sure you want to purge the task file? (y/n)")
	_, err := fmt.Scan(&ans)
	if err != nil {
		utils.Warnf("Not purging: %v", err)
		return false
	}
	ans = strings.TrimSpace(ans)
	ans = strings.ToLower(ans)

	if ans == "y" || ans == "yes" {
		return true
	}

	utils.Warn("Not Purging")
	return false
}
