package cmd

import (
	"github.com/spf13/cobra"
)

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Manage environment and its executables for an ide project",
	Long:  ``,
}

func init() {
	rootCmd.AddCommand(envCmd)
}
