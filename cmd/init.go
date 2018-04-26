package cmd

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/nebloc/gitdo/utils"
	"github.com/nebloc/gitdo/versioncontrol"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialises Gitdo in the current repository",
	Run: func(cmd *cobra.Command, args []string) {
		if err := Init(cmd, args); err != nil {
			pDanger("Failed to run Gitdo initialisation: %v\n", err)
			return
		}

		pNormal("Gitdo finished initialising\n")
	},
}

// Init initialises the gitdo project by scaffolding the gitdo folder
func Init(cmd *cobra.Command, args []string) error {
	// Initialise repo
	withVC := strings.ToLower(withVC)

	if withVC != "" {
		pInfo("Initialising: %s\n", withVC)
		switch withVC {
		case "git":
			if err := versioncontrol.NewGit().Init(); err != nil {
				return fmt.Errorf("could not create a Git repo: %v", err)
			}
		case "mercurial":
			if err := versioncontrol.NewHg().Init(); err != nil {
				return fmt.Errorf("could not create a Mercurial repo: %v", err)
			}
		default:
			return fmt.Errorf("could not initialise version control for %s", withVC)
		}
	}

	if err := ChangeToVCRoot(); err != nil {
		return fmt.Errorf("could not change to root directory: %v", err)
	}
	setVCPaths()

	pInfo("Making %s/gitdo", app.vc.NameOfDir())
	if err := os.MkdirAll(gitdoDir, os.ModePerm); err != nil {
		return err
	}

	if err := setConfig(); err != nil {
		return err
	}

	if err := CreatePluginsDir(); err != nil {
		return err
	}

	pInfo("Running plugin's setup...\n")
	if _, err := RunPlugin(SETUP, ""); err != nil {
		return err
	}

	if err := createHooks(); err != nil {
		return err
	}

	if err := setInterpFile(); err != nil {
		return err
	}

	fmt.Println("Done")
	return nil
}

// CreatePluginsDir creates a directory structure inside the Gitdo folder for Plugins to use as working space.
func CreatePluginsDir() error {
	path := filepath.Join(pluginDirPath, app.Plugin)
	err := os.MkdirAll(path, os.ModePerm)
	return err
}

// SetConfig checks the config is not set and asks the user relevant questions to set it
func setConfig() error {
	if app.IsSet() {
		return nil
	}

	if !app.authorIsSet() {
		author, err := askAuthor()
		if err != nil {
			return err
		}
		app.Author = author
	}

	if !app.pluginIsSet() {
		plugin, err := askPlugin()
		if err != nil {
			return err
		}
		app.Plugin = plugin
	}

	if !app.interpreterIsSet() {
		interp, err := getInterp()
		if err != nil {
			pWarning("No interp file in %s dir\n", app.Plugin)
			interp, err = askInterpreter()
			if err != nil {
				return err
			}

		}
		app.PluginInterpreter = interp
	}
	err := writeConfig()
	if err != nil {
		return fmt.Errorf("couldn't write config: %v", err)
	}
	return nil
}

// AskAuthor notifies user of email address used
func askAuthor() (string, error) {
	email, err := app.vc.GetEmail()
	if err != nil {
		return "", err
	}
	pInfo("Using %s\n", email)
	return email, nil
}

// AskPlugin reads in plugins from the directory and gives the user a list of plugins, that have a "<name>_getid"
func askPlugin() (string, error) {
	fmt.Println("Available plugins:")

	plugins, err := getPlugins()
	if err != nil {
		return "", err
	}
	if len(plugins) < 1 {
		pWarning("No plugins found\n")
		return "", fmt.Errorf("no plugins")
	}
	for i, name := range plugins {
		fmt.Printf("%d: %s\n", i+1, name)
	}

	chosen := false
	pN := 0

	for !chosen {
		fmt.Printf("What plugin would you like to use (1-%d): ", len(plugins))
		var choice string
		_, err = fmt.Scan(&choice)
		if err != nil {
			return "", err
		}
		pN, err = strconv.Atoi(strings.TrimSpace(choice))
		if err != nil || pN > len(plugins) || pN < 1 {
			continue
		}
		chosen = true
	}
	plugin := plugins[pN-1]

	pInfo("Using %s\n", plugin)
	return plugin, nil
}

// AskInterpreter asks the user what command they want to use to run the plugin
func askInterpreter() (string, error) {
	pWarning("Currently all plugins made as an example need python 3 set up in path. Redesign of plugin language choice and use coming soon.\n")
	var interp string
	for interp == "" {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("What interpreter for this plugin (i.e. python3/node/python): ")
		var err error
		interp, err = reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		interp = strings.TrimSpace(interp)
	}
	pInfo("Using %s\n", interp)
	return interp, nil
}

func getInterp() (string, error) {
	homePath, err := GetHomeDir()
	if err != nil {
		return "", err
	}
	contents, err := ioutil.ReadFile(filepath.Join(homePath, "plugins", app.Plugin, "interp"))
	if err != nil {
		return "", err
	}
	interp := utils.StripNewlineByte(contents)
	pInfo("Using %s - found in interp file\n", interp)
	return interp, err
}

func setInterpFile() error {
	homePath, err := GetHomeDir()
	if err != nil {
		return err
	}
	path := filepath.Join(homePath, "plugins", app.Plugin, "interp")
	data := []byte(app.PluginInterpreter)
	return ioutil.WriteFile(path, data, os.ModePerm)
}

// CreateHooks gets the users main Gitdo directory and copies the hooks from it to the correct version control hidden
// folder
func createHooks() error {
	homeDir, err := GetHomeDir()
	if err != nil {
		return err
	}
	pInfo("Copying hooks...\n")
	return app.vc.SetHooks(homeDir)
}
