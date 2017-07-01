package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Get the version of ide",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println(Version)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
