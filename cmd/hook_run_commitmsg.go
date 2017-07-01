package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var runCommitMsgHookCmd = &cobra.Command{
	Use:   "commit-msg",
	Short: "Run the commit-msg hook for the ide project",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Println("commit-msg")

		return nil
	},
}

func init() {
	runHookCmd.AddCommand(runCommitMsgHookCmd)
}
