package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nrocco/ide/pkg/ide/tools"
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

			switch filepath.Ext(path) {
			case ".go":
				tools.LintWhitespace(path, true, true, false)
				tools.GovetLinter.Exec(path).ForEachViolation(tools.PrintViolation)
				tools.GovetLinter.Exec(path).ForEachViolation(tools.PrintViolation)
				tools.GolintLinter.Exec(path).ForEachViolation(tools.PrintViolation)
				tools.GobuildLinter.Exec(path).ForEachViolation(tools.PrintViolation)
			case ".html":
				tools.LintWhitespace(path, true, true, true)
			case ".json":
				tools.LintWhitespace(path, true, true, true)
				tools.JqLinter.Exec(path).ForEachViolation(tools.PrintViolation)
			case ".php":
				tools.LintWhitespace(path, true, true, true)
				tools.PhpLinter.Exec(path).ForEachViolation(tools.PrintViolation)
				// TODO tools.LintPhpstan(path)
			case ".py":
				tools.LintWhitespace(path, true, true, true)
				tools.Flake8Linter.Exec(path).ForEachViolation(tools.PrintViolation)
			case ".rb":
				tools.LintWhitespace(path, true, true, true)
				tools.CookstyleLinter.Exec(path).ForEachViolation(tools.PrintViolation)
			case ".sh":
				tools.LintWhitespace(path, true, true, true)
				tools.ShellcheckLinter.Exec(path).ForEachViolation(tools.PrintViolation)
			case ".ts", ".vue", ".js":
				tools.LintWhitespace(path, true, true, true)
				tools.EsLintLinter.Exec(path).ForEachViolation(tools.PrintViolation)
			case ".yaml", ".yml":
				tools.LintWhitespace(path, true, true, true)
				tools.YamlLinter.Exec(path).ForEachViolation(tools.PrintViolation)
			default:
				tools.LintWhitespace(path, true, true, true)
			}
		}

		return nil
	},
}

func init() {
	lintCmd.AddCommand(runLintCmd)
}
