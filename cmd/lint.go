package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nrocco/ide/pkg/ide/tools"
	"github.com/spf13/cobra"
)

var lintCmd = &cobra.Command{
	Use:   "lint",
	Short: "Lint source code and report errors",
	Long:  "Lint source code and report errors",
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
				tools.LintGovet(path)
				tools.LintGolint(path)
				tools.LintGobuild(path)
			case ".html":
				tools.LintWhitespace(path, true, true, true)
			case ".json":
				tools.LintWhitespace(path, true, true, true)
				tools.LintJq(path)
			case ".php":
				tools.LintWhitespace(path, true, true, true)
				tools.LintPhp(path)
				// tools.LintPhpstan(path)
			case ".py":
				tools.LintWhitespace(path, true, true, true)
				tools.LintFlake8(path)
			case ".rb":
				tools.LintWhitespace(path, true, true, true)
				tools.LintCookstyle(path)
			case ".sh":
				tools.LintWhitespace(path, true, true, true)
				tools.LintShellcheck(path)
			case ".ts":
				tools.LintWhitespace(path, true, true, true)
				tools.LintEslint(path)
			case ".yaml":
				tools.LintWhitespace(path, true, true, true)
				tools.LintYaml(path)
			case ".yml":
				tools.LintWhitespace(path, true, true, true)
				tools.LintYaml(path)
			default:
				tools.LintWhitespace(path, true, true, true)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(lintCmd)
}
