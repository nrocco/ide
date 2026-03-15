package main

import (
	"os"

	"github.com/nrocco/ide/pkg/ide/linters"
	"github.com/spf13/cobra"
)

var testLintCmd = &cobra.Command{
	Use:   "tester",
	Short: "Test linter code",
	Long:  "Test linter code",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		linterResult, err := linters.NewLinterResult(os.Stdin, "<noname>", "<nofile>", linters.NewRegexMatcher(args[0]))
		if err != nil {
			return err
		}

		return linterResult.ForEachViolation(linters.PrintViolation)
	},
}

func init() {
	lintCmd.AddCommand(testLintCmd)
}
