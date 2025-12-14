package cmd

import (
	"github.com/spf13/cobra"
)

var lintCmd = &cobra.Command{
	Use:   "lint",
	Short: "Lint or fix source code",
	Long:  "Lint or fix source code",
}

func init() {
	rootCmd.AddCommand(lintCmd)
}
