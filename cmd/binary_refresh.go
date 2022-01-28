package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var refreshBinaryCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Refresh a binary for this ide project",
	Long:  "Refresh a binary for this ide project",
	Args:  cobra.ExactArgs(0),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return nil, cobra.ShellCompDirectiveNoFileComp
	},
	PreRunE: loadProject,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := project.RefreshBinaries(); err != nil {
			return err
		}

		log.Printf("Binaries refreshed\n")

		return nil
	},
}

func init() {
	binaryCmd.AddCommand(refreshBinaryCmd)
}
