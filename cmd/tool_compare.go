package cmd

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var compareCmd = &cobra.Command{
	Use:   "compare",
	Short: "Compare a file in the current project with another project",
	Long:  "Compare a file in the current project with another project",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		file := args[len(args)-1]
		options := []string{"-d", file}

		for _, project := range args[0 : len(args)-1] {
			options = append(options, filepath.Join(project, file))
		}

		command := exec.Command("nvim", options...)
		command.Stdin = os.Stdin
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr

		return command.Run()
	},
}

func init() {
	toolCmd.AddCommand(compareCmd)
}
