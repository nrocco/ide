package cmd

import (
	"bufio"
	"os"
	"regexp"

	"github.com/nrocco/ide/pkg/ide/tools"
	"github.com/spf13/cobra"
)

var testLintCmd = &cobra.Command{
	Use:   "tester",
	Short: "Test linter code",
	Long:  "Test linter code",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		matcher, err := regexp.Compile(args[0])
		if err != nil {
			return err
		}

		linterResult := tools.LinterResult{
			Scanner: bufio.NewScanner(os.Stdin),
			Name: "<noname>",
			File: "<nofile>",
			Matcher: matcher,
		}

		return linterResult.ForEachViolation(tools.PrintViolation)
	},
}

func init() {
	lintCmd.AddCommand(testLintCmd)
}
