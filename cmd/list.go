package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// List pretty prints the tasks that are in file
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "lists data stored by Gitdo",
}

var listTasksCmd = &cobra.Command{
	Use:   "tasks",
	Short: "lists the program tasks",
	Run: func(cmd *cobra.Command, args []string) {
		if err := Setup(); err != nil {
			fmt.Fprintf(os.Stderr, "Could not print tasks: %v", err)
			return
		}
		tasks, err := getTasksFile()
		if err != nil && !os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Could not get task file: %v", err)
			return
		}

		fmt.Println(tasks.String())
	},
}

var listConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "lists the current configuration",
	Run: func(cmd *cobra.Command, args []string) {
		if err := Setup(); err != nil {
			fmt.Fprintf(os.Stderr, "Could not print configuration: %v", err)
			return
		}
		fmt.Println(config.String())
	},
}
