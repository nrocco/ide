package cmd

import (
	"github.com/spf13/cobra"
	"log"
)

var unlinkLinkCmd = &cobra.Command{
	Use:   "rm [name]",
	Short: "Remove a linked program from this ide project",
	Long:  "Remove a linked program from this ide project",
	RunE: func(cmd *cobra.Command, args []string) error {
		for _, name := range args {
			err := Project.RemoveExecutable(name)
			if err != nil {
				log.Println(err)
			} else {
				log.Printf("Executable %s removed\n", name)
			}
		}

		return nil
	},
}

func init() {
	linkCmd.AddCommand(unlinkLinkCmd)
}
