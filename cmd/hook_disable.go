package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var disableHookCmd = &cobra.Command{
	Use:   "disable",
	Short: "Disable a git hook for this ide project",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Println("disable")
		return nil
	},
}

func init() {
	hookCmd.AddCommand(disableHookCmd)
}
