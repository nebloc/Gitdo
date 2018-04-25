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
	config = &Config{
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

	rootCmd := &cobra.Command{
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
	rootCmd.AddCommand(versionCmd)

	// INIT
	rootCmd.AddCommand(initCmd)

	// LIST
	listCmd.AddCommand(listConfigCmd)
	listCmd.AddCommand(listTasksCmd)
	rootCmd.AddCommand(listCmd)

	// COMMIT
	rootCmd.AddCommand(commitCmd)

	// POST COMMIT
	rootCmd.AddCommand(postCommitCmd)

	// PUSH
	rootCmd.AddCommand(pushCmd)

	return rootCmd
}

func versionString(version string) string {
	if version == "" {
		pWarning("No version number set on this build\n")
	}
	return fmt.Sprintf("Version: %s\nBuild: %s_%s", version, runtime.GOOS, runtime.GOARCH)
}

/**
// AppBuilder returns a urfave/cli app for directing commands and running setup
func AppBuilder() *cli.App {
	gitdo := cli.NewApp()
	gitdo.Name = "gitdo"
	gitdo.Usage = "track source code TODO comments - https://github.com/nebloc/Gitdo"
	gitdo.Version = "0.0.0-A5"
	if version != "" {
		gitdo.Version = fmt.Sprintf("App: %s, Build: %s_%s", version, runtime.GOOS, runtime.GOARCH)
	}
	cli.VersionFlag = cli.BoolFlag{
		Name:  "version, V",
		Usage: "print the app version",
	}
	gitdo.Commands = []cli.Command{
		{
			Before: ChangeToVCRoot,
			Name:   "list",
			Usage:  "prints the json of staged tasks",
			Flags:  []cli.Flag{cli.BoolFlag{Name: "config", Usage: "prints the current configuration"}},
			Action: List,
		},
		{

			Name:   "commit",
			Usage:  "gets git diff and stages any new tasks - normally ran from pre-commit hook",
			Action: Commit,
			Before: LoadConfig,
			After:  NotifyFinished,
		},
		{
			Name:  "init",
			Usage: "sets the gitdo configuration for the current git repo",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "with-vc",
					Usage: "Must be 'Mercurial' or 'Git'. Initialises a repo first",
				},
			},
			Action: Init,
		},
		{
			Before: LoadConfig,
			Name:   "post-commit",
			Usage:  "adds the commit hash that has just been committed to tasks with empty hash fields",
			Action: PostCommit,
			After:  NotifyFinished,
		},
		{
			Name:   "push",
			Usage:  "starts the plugin to move staged tasks into your task manager - normally ran from pre-push hook",
			Action: Push,
			Before: LoadConfig,
			After:  NotifyFinished,
		},
		{
			Name:   "destroy",
			Usage:  "deletes all of the stored tasks",
			Flags:  []cli.Flag{cli.BoolFlag{Name: "yes", Usage: "confirms purge of task file"}},
			Before: ConfirmUser,
			Action: Destroy,
		},
	}
	return gitdo
}
*/
