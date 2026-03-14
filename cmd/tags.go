package cmd

import (
	"github.com/spf13/cobra"
)

var tagsCmd = &cobra.Command{
	Use:   "tags",
	Short: "Manage tags",
	Long:  "Manage tags",
}

func init() {
	rootCmd.AddCommand(tagsCmd)
}
