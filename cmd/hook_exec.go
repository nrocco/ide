package cmd

import (
	"github.com/spf13/cobra"
	"log"
)

// hookCmd represents the status command
var execHookCmd = &cobra.Command{
	Use:   "exec",
	Short: "Execute a git hook against an ide project",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			log.Fatalln("You need to provide a git hook name")
		}

		hook := args[0]

		log.Println("exec " + hook)
		return nil
	},
}

func init() {
	hookCmd.AddCommand(execHookCmd)
}
