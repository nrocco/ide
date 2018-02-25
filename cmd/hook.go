package cmd

import (
	"github.com/spf13/cobra"
)

var hookCmd = &cobra.Command{
	Use:   "hook",
	Short: "Manage git hooks for an ide project",
	Long:  "Manage git hooks for an ide project",
}

func init() {
	rootCmd.AddCommand(hookCmd)
}
