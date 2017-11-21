package cmd

import (
	"github.com/spf13/cobra"
)

var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Manage executables for this ide project",
	Long:  "Manage executables for this ide project",
}

func init() {
	rootCmd.AddCommand(execCmd)
}
