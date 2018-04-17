package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/urfave/cli"
	"io"
)

// Init initialises the gitdo project by scaffolding the gitdo folder
func Init(ctx *cli.Context) error {
	if ctx.Bool("with-git") {
		if err := InitGit(); err != nil {
			return err
		}
	}

	if err := os.MkdirAll(gitdoDir, os.ModePerm); err != nil {
		return err
	}

	if err := CreateHooks(); err != nil {
		return err
	}

	if err := SetConfig(); err != nil {
		return err
	}

	if err := CreatePlugins(); err != nil {
		return err
	}

	if _, err := RunPlugin(SETUP, ""); err != nil {
		return err
	}

	fmt.Println("Done")
	return nil
}

func CreatePlugins() error {
	path := filepath.Join(pluginDirPath, config.Plugin)
	Warn(path)
	err := os.MkdirAll(path, os.ModePerm)
	return err
}

// SetConfig checks the config is not set and asks the user relevant questions to set it
func SetConfig() error {
	if config.IsSet() {
		return nil
	}

	if !config.authorIsSet() {
		author, err := AskAuthor()
		if err != nil {
			return err
		}
		config.Author = author
	}

	if !config.pluginIsSet() {
		plugin, err := AskPlugin()
		if err != nil {
			return err
		}
		config.Plugin = plugin
	}

	if !config.interpreterIsSet() {
		interp, err := AskInterpreter()
		if err != nil {
			return err
		}
		config.PluginInterpreter = interp
	}

	err := WriteConfig()
	if err != nil {
		Dangerf("Couldn't save config: %v", err)
		return err
	}
	return nil
}

// InitGit initialises a git repo before initialising gitdo
func InitGit() error {
	fmt.Println("Initializing git...")
	cmd := exec.Command("git", "init")
	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	fmt.Println("Git initialized")
	return nil
}

// AskAuthor notifies user of email address used
func AskAuthor() (string, error) {
	email, err := getGitEmail()
	if err != nil {
		return "", err
	}
	Highlightf("Using %s", email)
	return email, nil
}

// AskPlugin reads in plugins from the directory and gives the user a list of plugins, that have a "<name>_getid"
func AskPlugin() (string, error) {
	fmt.Println("Available plugins:")

	plugins, err := GetPlugins()
	if err != nil {
		return "", err
	}
	if len(plugins) < 1 {
		Warn("No plugins found")
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

	Highlightf("Using %s", plugin)
	return plugin, nil
}

// AskInterpreter asks the user what command they want to use to run the plugin
func AskInterpreter() (string, error) {
	Warn("Currently all plugins made as an example need python 3 set up in path. Redesign of plugin language choice and use coming soon.")
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
	Highlightf("Using %s", interp)
	return interp, nil
}

func CreateHooks() error {
	homeDir, err := GetHomeDir()
	if err != nil {
		return err
	}

	switch config.VC {
	case GIT:
		srcHooks := filepath.Join(homeDir, "hooks")
		dstHooks := filepath.Join(".git", "hooks")

		err = os.MkdirAll(dstHooks, os.ModePerm)
		if err != nil {
			return fmt.Errorf("could not create hook dir inside .git/gitdo: %v", err)
		}

		err = copyFolder(srcHooks, dstHooks)
		if err != nil {
			return err
		}
	case HG:
		srcHook := filepath.Join(homeDir, "hgrc")
		dstHook := filepath.Join(".hg", "hgrc")
		err = copyFile(srcHook, dstHook)
		if err != nil {
			return fmt.Errorf("could not move .hgrc to inside .hgrc: %v", err)
		}
	}

	return nil
}

// copyFolder copies a folder from src to dst. It looks through the src folder and copies files one by one to the
// destination folder. It does not copy subdirectories
func copyFolder(src, dst string) error {
	fmt.Printf("Copying from: %s to %s\n", src, dst)
	files, err := ioutil.ReadDir(src)
	if err != nil {
		return fmt.Errorf("could not get %s files: %v", src, err)
	}
	for _, file := range files {
		sf := filepath.Join(src, file.Name())
		df := filepath.Join(dst, file.Name())
		err = copyFile(sf, df)
		if err != nil {
			fmt.Printf("could not copy %v - skipping: %v\n", file.Name(), err)
		}
	}

	return nil
}

// copyFile copies a file from src to dst. If src and dst files exist, and are the same, then return success. Otherise,
// attempt to create a hard link between the two files. If that fail, copy the file contents from src to dst.
func copyFile(src, dst string) (err error) {
	sfi, err := os.Stat(src)
	if err != nil {
		return
	}
	if !sfi.Mode().IsRegular() {
		// cannot copy non-regular files (e.g., directories,
		// symlinks, devices, etc.)
		return fmt.Errorf("copyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
	}
	dfi, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
	} else {
		if !(dfi.Mode().IsRegular()) {
			return fmt.Errorf("copyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
		}
		if os.SameFile(sfi, dfi) {
			return
		}
	}
	if err = os.Link(src, dst); err == nil {
		return
	}
	err = copyFileContents(src, dst)
	return
}

// copyFileContents copies the contents of the file named src to the file named by dst. The file will be created if it
// does not already exist. If the destination file exists, all it's contents will be replaced by the contents of the
// source file.
func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}
