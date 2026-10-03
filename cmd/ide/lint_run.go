package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/nrocco/ide/pkg/ide/linters"
	"github.com/spf13/cobra"
)

var runLintCmdOpts struct {
	Debug bool
}

var (
	goFileRe     = regexp.MustCompile(`\.go$`)
	htmlFileRe   = regexp.MustCompile(`\.html$`)
	jsonFileRe   = regexp.MustCompile(`\.json$`)
	makefileRe   = regexp.MustCompile(`^(GNUmakefile|[Mm]akefile|.*\.mk)$`)
	phpFileRe    = regexp.MustCompile(`\.php$`)
	pythonFileRe = regexp.MustCompile(`\.py$`)
	rubyFileRe   = regexp.MustCompile(`\.rb$`)
	shellFileRe  = regexp.MustCompile(`\.sh$`)
	scriptFileRe = regexp.MustCompile(`\.(ts|vue|js)$`)
	yamlFileRe   = regexp.MustCompile(`\.ya?ml$`)
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

			name := filepath.Base(path)

			switch {
			case goFileRe.MatchString(name):
				linters.LintWhitespace(path, true, true, false)
				linters.GovetLinter.Exec(path, runLintCmdOpts.Debug).ForEachViolation(linters.PrintViolation)
				linters.GolintLinter.Exec(path, runLintCmdOpts.Debug).ForEachViolation(linters.PrintViolation)
				linters.GobuildLinter.Exec(path, runLintCmdOpts.Debug).ForEachViolation(linters.PrintViolation)
			case htmlFileRe.MatchString(name):
				linters.LintWhitespace(path, true, true, true)
			case jsonFileRe.MatchString(name):
				linters.LintWhitespace(path, true, true, true)
				linters.JqLinter.Exec(path, runLintCmdOpts.Debug).ForEachViolation(linters.PrintViolation)
			case makefileRe.MatchString(name):
				linters.LintWhitespace(path, true, true, false)
			case phpFileRe.MatchString(name):
				linters.LintWhitespace(path, true, true, true)
				linters.PhpLinter.Exec(path, runLintCmdOpts.Debug).ForEachViolation(linters.PrintViolation)
				// TODO linters.LintPhpstan(path)
			case pythonFileRe.MatchString(name):
				linters.LintWhitespace(path, true, true, true)
				linters.Flake8Linter.Exec(path, runLintCmdOpts.Debug).ForEachViolation(linters.PrintViolation)
			case rubyFileRe.MatchString(name):
				linters.LintWhitespace(path, true, true, true)
				linters.CookstyleLinter.Exec(path, runLintCmdOpts.Debug).ForEachViolation(linters.PrintViolation)
			case shellFileRe.MatchString(name):
				linters.LintWhitespace(path, true, true, true)
				linters.ShellcheckLinter.Exec(path, runLintCmdOpts.Debug).ForEachViolation(linters.PrintViolation)
			case scriptFileRe.MatchString(name):
				linters.LintWhitespace(path, true, true, true)
				linters.EsLintLinter.Exec(path, runLintCmdOpts.Debug).ForEachViolation(linters.PrintViolation)
			case yamlFileRe.MatchString(name):
				linters.LintWhitespace(path, true, true, true)
				linters.YamlLinter.Exec(path, runLintCmdOpts.Debug).ForEachViolation(linters.PrintViolation)
			default:
				linters.LintWhitespace(path, true, true, true)
			}
		}

		return nil
	},
}

func init() {
	runLintCmd.Flags().BoolVar(&runLintCmdOpts.Debug, "debug", false, "Debug linter output")

	lintCmd.AddCommand(runLintCmd)
}
