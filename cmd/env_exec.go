package cmd

import (
	"github.com/spf13/cobra"
	"log"
)

var execEnvCmd = &cobra.Command{
	Use:   "exec",
	Short: "Execute a binary in this ide project's environment",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO get container name based on first arg
		// TODO docker exec binary in the container
		log.Println("exec")
		return nil
	},
}

func init() {
	envCmd.AddCommand(execEnvCmd)
}
