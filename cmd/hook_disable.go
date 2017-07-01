package cmd

import (
	"github.com/spf13/cobra"
	"log"
)

// hookCmd represents the status command
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
