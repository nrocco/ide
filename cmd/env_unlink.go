package cmd

import (
	"github.com/spf13/cobra"
	"log"
)

var unlinkEnvCmd = &cobra.Command{
	Use:   "unlink",
	Short: "Unlink and executable for this ide project",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO remove link from .git/bin
		// TODO remove git config option from ide.executables
		log.Println("unlink")
		return nil
	},
}

func init() {
	envCmd.AddCommand(unlinkEnvCmd)
}
