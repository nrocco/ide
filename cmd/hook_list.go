package cmd

import (
	"github.com/spf13/cobra"
	"log"
)

// hookCmd represents the status command
var listHookCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available hooks",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Println("list")
		return nil
	},
}

func init() {
	hookCmd.AddCommand(listHookCmd)
}
