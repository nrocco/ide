package cmd

import (
	"github.com/spf13/cobra"
)

// hookCmd represents the status command
var runHookCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a git hook against an ide project",
	Long:  ``,
}

func init() {
	hookCmd.AddCommand(runHookCmd)
}