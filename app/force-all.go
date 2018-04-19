package main

import (
	"github.com/nebloc/gitdo/app/utils"
	"github.com/urfave/cli"
)

func ForceAll(ctx *cli.Context) error {
	if err := config.vc.CreateBranch(); err != nil {
		utils.Danger("Could not CreateBranch")
		return err
	}
	if err := config.vc.SwitchBranch(); err != nil {
		utils.Danger("Could not SwitchBranch")
		return err
	}

	return nil
}
