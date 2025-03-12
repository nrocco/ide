package cmd

import (
	"log"
	"strings"

	"github.com/spf13/cobra"
)

var addShimCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a shim for this ide project",
	Long:  "Add a shim for this ide project",
	Args:  cobra.MinimumNArgs(2),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return nil, cobra.ShellCompDirectiveNoFileComp
	},
	PreRunE: loadProject,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := project.ShimAdd(args[0], strings.Join(args[1:], " "))
		if err != nil {
			return err
		}

		log.Printf("Shim %s add\n", args[0])

		return project.AddGitBinToPath()
	},
}

func init() {
	shimCmd.AddCommand(addShimCmd)
}
