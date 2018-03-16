package main

import (
	cli "github.com/urfave/cli"
)

var gitdo *cli.App

func init() {
	gitdo = AppBuilder()
	cachedFlag = true
	verboseLogFlag = true
}

/**
func TestSetup(t *testing.T) {
	args := []string{"not_a_command"}
	t.Logf("Args: %v", args)

	err := gitdo.Run(args)
	if err != nil {
		t.Error(err)
	}
}
*/
