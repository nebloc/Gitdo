package cmd

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/nebloc/gitdo/utils"
	"github.com/spf13/cobra"
)

var postCommitCmd = &cobra.Command{
	Use:   "post-commit",
	Short: "Tags all entries in the tasks.json file with current data",
	Long:  "Tags all entries in the tasks.json file with the current version control hash, and branch.",
	Run: func(cmd *cobra.Command, args []string) {
		if err := Setup(); err != nil {
			pDanger("Could not load gitdo: %v\n", err)
			return
		}
		if err := PostCommit(cmd, args); err != nil {
			pDanger("Failed to run post-commit: %v\n", err)
			return
		}

		pNormal("Gitdo finished post-commit process\n")
	},
}

// PostCommit is ran from a git post-commit hook to set the hash values and branch values of any tasks that have just
// been committed
func PostCommit(cmd *cobra.Command, args []string) error {
	hash, err := config.vc.GetHash()
	if err != nil {
		return err
	}
	branch, err := config.vc.GetHash()
	if err != nil {
		return err
	}

	tasks, err := getTasksFile()
	if err != nil {
		utils.Warn("No tasks file")
		return nil
	}
	for id, task := range tasks.NewTasks {
		if task.Hash == "" {
			task.Hash = hash
			task.Branch = branch
			tasks.NewTasks[id] = task
		}
	}

	bUpdated, err := json.MarshalIndent(tasks, "", "\t")
	if err != nil {
		utils.Danger("couldn't marshal tasks with added hash")
		return err
	}
	err = ioutil.WriteFile(stagedTasksFile, bUpdated, os.ModePerm)
	if err != nil {
		utils.Danger("couldn't write tasks with hash back to tasks.json")
		return err
	}
	return nil
}
