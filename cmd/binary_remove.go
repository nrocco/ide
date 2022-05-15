package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var removeBinaryCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a binary for this ide project",
	Long:  "Remove a binary for this ide project",
	Args:  cobra.MinimumNArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		binaries := []string{}
		if len(args) == 0 {
			if err := loadProject(cmd, args); err == nil {
				for binary := range project.ListBinaries() {
					binaries = append(binaries, binary)
				}
			}
		}
		return binaries, cobra.ShellCompDirectiveNoFileComp
	},
	PreRunE: loadProject,
	RunE: func(cmd *cobra.Command, args []string) error {
		for _, binary := range args {
			err := project.RemoveBinary(binary)
			if err != nil {
				log.Println(err)
			} else {
				log.Printf("Binary %s removed\n", binary)
			}
		}

		return nil
	},
}

func init() {
	binaryCmd.AddCommand(removeBinaryCmd)
}
