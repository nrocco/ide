package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
)

var unlinkEnvCmd = &cobra.Command{
	Use:   "unlink",
	Short: "Unlink and executable for this ide project",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return nil
		}

		executable := filepath.Join(".git", "bin", args[0])

		log.Printf("Unlinking %s", executable)

		Project.RemoveExecutable(args[0])

		return os.Remove(executable)
	},
}

func init() {
	envCmd.AddCommand(unlinkEnvCmd)
}
