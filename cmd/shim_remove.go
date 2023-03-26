package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var removeShimCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a shim for this ide project",
	Long:  "Remove a shim for this ide project",
	Args:  cobra.MinimumNArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		shims := []string{}
		if len(args) == 0 {
			if err := loadProject(cmd, args); err == nil {
				for shim := range project.ListShims() {
					shims = append(shims, shim)
				}
			}
		}
		return shims, cobra.ShellCompDirectiveNoFileComp
	},
	PreRunE: loadProject,
	RunE: func(cmd *cobra.Command, args []string) error {
		for _, shim := range args {
			err := project.RemoveShim(shim)
			if err != nil {
				log.Println(err)
			} else {
				log.Printf("Shim %s removed\n", shim)
			}
		}

		return nil
	},
}

func init() {
	shimCmd.AddCommand(removeShimCmd)
}
