package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var refreshShimCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Refresh a shim for this ide project",
	Long:  "Refresh a shim for this ide project",
	Args:  cobra.ExactArgs(0),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return nil, cobra.ShellCompDirectiveNoFileComp
	},
	PreRunE: loadProject,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := project.RefreshShims(); err != nil {
			return err
		}

		log.Printf("Shims refreshed\n")

		return nil
	},
}

func init() {
	shimCmd.AddCommand(refreshShimCmd)
}
