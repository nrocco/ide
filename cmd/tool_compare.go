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

		vim := exec.Command("vim", options...)
		vim.Stdin = os.Stdin
		vim.Stdout = os.Stdout
		vim.Stderr = os.Stderr

		return vim.Run()
	},
}

func init() {
	toolCmd.AddCommand(compareCmd)
}
