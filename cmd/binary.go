package cmd

import (
	"github.com/spf13/cobra"
)

var binaryCmd = &cobra.Command{
	Use:   "binary",
	Short: "Manage binaries for an ide project",
	Long:  "Manage binaries for an ide project",
}

func init() {
	rootCmd.AddCommand(binaryCmd)
}
