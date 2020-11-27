package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var disableHookCmd = &cobra.Command{
	Use:   "disable",
	Short: "Disable a git hook for this ide project",
	Long:  "Disable a git hook for this ide project",
	Args:  cobra.MinimumNArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return project.ListHooks(), cobra.ShellCompDirectiveNoFileComp // TODO project is not initialized at this point
	},
	PreRunE: loadProject,
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
