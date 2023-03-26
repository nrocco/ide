package cmd

import (
	"github.com/spf13/cobra"
)

var shimCmd = &cobra.Command{
	Use:   "shim",
	Short: "Manage shims for an ide project",
	Long:  "Manage shims for an ide project",
}

func init() {
	rootCmd.AddCommand(shimCmd)
}
