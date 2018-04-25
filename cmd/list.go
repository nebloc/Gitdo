package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists data stored by Gitdo",
}

var listTasksCmd = &cobra.Command{
	Use:   "tasks",
	Short: "Lists the program tasks",
	Run: func(cmd *cobra.Command, args []string) {
		if err := setup(); err != nil {
			fmt.Fprintf(os.Stderr, "Could not print tasks: %v\n", err)
			return
		}
		tasks, err := getTasksFile()
		if err != nil && !os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Could not get task file: %v\n", err)
			return
		}

		fmt.Println(tasks.String())
	},
}

var listConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Lists the current configuration",
	Run: func(cmd *cobra.Command, args []string) {
		if err := setup(); err != nil {
			fmt.Fprintf(os.Stderr, "Could not print configuration: %v\n", err)
			return
		}
		fmt.Println(app.String())
	},
}
