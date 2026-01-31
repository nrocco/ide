package cmd

import (
	"github.com/spf13/cobra"
)

var toolsCmd = &cobra.Command{
	Use:   "tools",
	Short: "Miscellaneous tools",
	Long:  "A collection of miscellaneous tools",
}

func init() {
	rootCmd.AddCommand(toolsCmd)
}
