package cmd

import (
	"fmt"
	"runtime"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	// Coloured Outputs
	pDanger  = color.New(color.FgHiRed).PrintfFunc()
	pWarning = color.New(color.FgHiYellow).PrintfFunc()
	pInfo    = color.New(color.FgHiCyan).PrintfFunc()
	pNormal  = fmt.Printf
)

var (
	// Config needed for commit and push to use plugins and add author metadata
	// to task
	app = &config{
		Author:            "",
		Plugin:            "",
		PluginInterpreter: "",
	}

	// Gitdo working directory (holds plugins, secrets, tasks, etc.)
	gitdoDir string

	// File name for writing and reading staged tasks from (between commit
	// and push)
	stagedTasksFile string
	configFilePath  string
	pluginDirPath   string

	// FLAGS
	withVC string
)

// New creates a new base command for executing Gitdo
func New(version string) *cobra.Command {
	initCmd.PersistentFlags().StringVarP(&withVC, "with-vc", "w", "", "Initialises repository as well as gitdo. Supports 'Git' and 'Mercurial'")

	gitdoCmd := &cobra.Command{
		Use:   "gitdo",
		Short: "A tool for tracking task annotations using version control systems.",
		Long: fmt.Sprintf(`A tool for tracking task annotations using version control systems.

%s

Please run gitdo help to see a list of commands.
More information and documentation can be found at https://github.com/nebloc/gitdo`, versionString(version)),
	}

	// VERSION
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Prints the version number of the current Gitdo app.",
		Run: func(*cobra.Command, []string) {
			fmt.Println(versionString(version))
		},
	}
	gitdoCmd.AddCommand(versionCmd)

	// INIT
	gitdoCmd.AddCommand(initCmd)

	// LIST
	listCmd.AddCommand(listConfigCmd)
	listCmd.AddCommand(listTasksCmd)
	gitdoCmd.AddCommand(listCmd)

	// COMMIT
	gitdoCmd.AddCommand(commitCmd)

	// POST COMMIT
	gitdoCmd.AddCommand(postCommitCmd)

	// PUSH
	gitdoCmd.AddCommand(pushCmd)

	return gitdoCmd
}

func versionString(version string) string {
	if version == "" {
		pWarning("No version number set on this build\n")
	}
	return fmt.Sprintf("Version: %s\nBuild: %s_%s", version, runtime.GOOS, runtime.GOARCH)
}
