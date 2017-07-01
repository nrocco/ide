package cmd

import (
	"github.com/spf13/cobra"
)

// execCmd represents the status command
var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Run programs against an ide project",
	Long:  ``,
}

func init() {
	RootCmd.AddCommand(execCmd)
}
