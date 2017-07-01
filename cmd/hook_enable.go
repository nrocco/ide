package cmd

import (
	"github.com/spf13/cobra"
	"log"
)

// hookCmd represents the status command
var enableHookCmd = &cobra.Command{
	Use:   "enable",
	Short: "Enable a git hook for this ide project",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Println("enable")
		return nil
	},
}

func init() {
	hookCmd.AddCommand(enableHookCmd)
}
