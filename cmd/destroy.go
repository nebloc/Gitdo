/**
Not really maintained - not a very useful command to users in it's current state, should make it uninstall all from the
project in future
*/

package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// Destroy deletes the staged tasks file if you need it to
func Destroy(cmd *cobra.Command, args []string) error {
	if ConfirmWithUser("Are you sure you want to purge the tasks.json?") {
		pWarning("Deleting")
		return os.Remove(stagedTasksFile)
	}
	pWarning("Cancelling")
	return nil
}

// ConfirmWithUser asks the user a message with Y/N and returns true if their answer is yes or Y
func ConfirmWithUser(message string) bool {
	var ans string
	pNormal("%s %s", message, "(y/n): ")
	_, err := fmt.Scan(&ans)
	if err != nil {
		return false
	}
	ans = strings.TrimSpace(ans)
	ans = strings.ToLower(ans)

	if ans == "y" || ans == "yes" {
		return true
	}
	return false
}
