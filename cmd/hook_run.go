package cmd

import (
	"github.com/spf13/cobra"
)

var runHookCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a git hook against an ide project",
	Long:  "Run a git hook against an ide project",
}

func init() {
	hookCmd.AddCommand(runHookCmd)
}
