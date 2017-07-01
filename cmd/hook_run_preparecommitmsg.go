package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var runPrepareCommitMsgHookCmd = &cobra.Command{
	Use:   "prepare-commit-msg",
	Short: "Run the prepare-commit-msg hook for the ide project",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Println("prepare-commit-msg")

		return nil
	},
}

func init() {
	runHookCmd.AddCommand(runPrepareCommitMsgHookCmd)
}
