package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nrocco/ide/pkg/ide/fixers"
	"github.com/spf13/cobra"
)

var fixLintCmd = &cobra.Command{
	Use:   "fix",
	Short: "Fix source code",
	Long:  "Fix source code",
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
				fixers.FixWhitespace(path)
				fixers.FixClrf(path)
				fixers.FixGofmt(path)
				fixers.FixGoimports(path)
			case ".html":
				fixers.FixWhitespace(path)
				fixers.FixClrf(path)
			case ".json":
				fixers.FixWhitespace(path)
				fixers.FixClrf(path)
				fixers.FixJq(path)
			case ".php":
				fixers.FixWhitespace(path)
				fixers.FixClrf(path)
				fixers.FixPhpcsfixer(path)
			case ".py":
				fixers.FixWhitespace(path)
				fixers.FixClrf(path)
			case ".rb":
				fixers.FixWhitespace(path)
				fixers.FixClrf(path)
			case ".sh":
				fixers.FixWhitespace(path)
				fixers.FixClrf(path)
			case ".ts", ".vue", ".js":
				fixers.FixWhitespace(path)
				fixers.FixClrf(path)
				fixers.FixEslint(path)
			case ".yaml", ".yml":
				fixers.FixWhitespace(path)
				fixers.FixClrf(path)
			default:
				fixers.FixWhitespace(path)
				fixers.FixClrf(path)
			}
		}

		return nil
	},
}

func init() {
	lintCmd.AddCommand(fixLintCmd)
}
