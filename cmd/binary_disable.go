package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var disableBinaryCmd = &cobra.Command{
	Use:   "disable",
	Short: "Disable a binary for this ide project",
	Long:  "Disable a binary for this ide project",
	Args:  cobra.MinimumNArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		binaries := []string{}
		for binary := range project.ListBinaries() {
			binaries = append(binaries, binary)
		}
		return binaries, cobra.ShellCompDirectiveNoFileComp // TODO project is not initialized at this point
	},
	PreRunE: loadProject,
	RunE: func(cmd *cobra.Command, args []string) error {
		for _, binary := range args {
			err := project.DisableBinary(binary)
			if err != nil {
				log.Println(err)
			} else {
				log.Printf("Binary %s disabled\n", binary)
			}
		}

		return nil
	},
}

func init() {
	binaryCmd.AddCommand(disableBinaryCmd)
}
