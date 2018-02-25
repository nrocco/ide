package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var disableHookCmd = &cobra.Command{
	Use:   "disable",
	Short: "Disable a git hook for this ide project",
	Long:  "Disable a git hook for this ide project",
	RunE: func(cmd *cobra.Command, args []string) error {
		for _, hook := range args {
			err := project.DisableHook(hook)
			if err != nil {
				log.Println(err)
			} else {
				log.Printf("Hook %s disabled\n", hook)
			}
		}

		return nil
	},
}

func init() {
	hookCmd.AddCommand(disableHookCmd)
}
