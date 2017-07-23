package cmd

import (
	"github.com/spf13/cobra"
)

var linkCmd = &cobra.Command{
	Use:   "link",
	Short: "Manage link to executables for this ide project",
	Long:  "Manage link to executables for this ide project",
}

func init() {
	rootCmd.AddCommand(linkCmd)
}
