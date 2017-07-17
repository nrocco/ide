package cmd

import (
	"github.com/spf13/cobra"
	"log"
)

var destroyEnvCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Destroy an environment for this ide project",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Println("docker stop|rm envname")
		return nil
	},
}

func init() {
	envCmd.AddCommand(destroyEnvCmd)
}
