package cmd

import (
	"github.com/spf13/cobra"
)

var toolCmd = &cobra.Command{
	Use:   "tool",
	Short: "Miscellaneous tools",
	Long:  "A collection of miscellaneous tools",
}

func init() {
	rootCmd.AddCommand(toolCmd)
}
