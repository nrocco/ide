package cmd

import (
	"github.com/spf13/cobra"
)

var fafaCmd = &cobra.Command{
	Use:     "fafa",
	Hidden:  true,
	PreRunE: loadProject,
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	rootCmd.AddCommand(fafaCmd)
}
