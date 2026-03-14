package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nrocco/ide/pkg/ide/linters"
	"github.com/spf13/cobra"
)

var runLintCmd = &cobra.Command{
	Use:   "run",
	Short: "Lint source code and report errors",
	Long:  "Lint source code and report errors",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		for _, path := range args {
			fileInfo, err := os.Stat(path)
			if err != nil {
				return err
			} else if fileInfo.IsDir() {
				return fmt.Errorf("%s is a directory", path)
			}

			debug, _ := cmd.Flags().GetBool("debug")

			switch filepath.Ext(path) {
			case ".go":
				linters.LintWhitespace(path, true, true, false)
				linters.GovetLinter.Exec(path, debug).ForEachViolation(linters.PrintViolation)
				linters.GolintLinter.Exec(path, debug).ForEachViolation(linters.PrintViolation)
				linters.GobuildLinter.Exec(path, debug).ForEachViolation(linters.PrintViolation)
			case ".html":
				linters.LintWhitespace(path, true, true, true)
			case ".json":
				linters.LintWhitespace(path, true, true, true)
				linters.JqLinter.Exec(path, debug).ForEachViolation(linters.PrintViolation)
			case ".php":
				linters.LintWhitespace(path, true, true, true)
				linters.PhpLinter.Exec(path, debug).ForEachViolation(linters.PrintViolation)
				// TODO linters.LintPhpstan(path)
			case ".py":
				linters.LintWhitespace(path, true, true, true)
				linters.Flake8Linter.Exec(path, debug).ForEachViolation(linters.PrintViolation)
			case ".rb":
				linters.LintWhitespace(path, true, true, true)
				linters.CookstyleLinter.Exec(path, debug).ForEachViolation(linters.PrintViolation)
			case ".sh":
				linters.LintWhitespace(path, true, true, true)
				linters.ShellcheckLinter.Exec(path, debug).ForEachViolation(linters.PrintViolation)
			case ".ts", ".vue", ".js":
				linters.LintWhitespace(path, true, true, true)
				linters.EsLintLinter.Exec(path, debug).ForEachViolation(linters.PrintViolation)
			case ".yaml", ".yml":
				linters.LintWhitespace(path, true, true, true)
				linters.YamlLinter.Exec(path, debug).ForEachViolation(linters.PrintViolation)
			default:
				linters.LintWhitespace(path, true, true, true)
			}
		}

		return nil
	},
}

func init() {
	runLintCmd.Flags().Bool("debug", false, "Debug linter output")

	lintCmd.AddCommand(runLintCmd)
}
