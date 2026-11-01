package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/nrocco/ide/pkg/ide/fixers"
	"github.com/spf13/cobra"
)

var fixLintCmd = &cobra.Command{
	Use:   "fix",
	Short: "Fix source code",
	Long:  "Fix source code",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var failed []error

		for _, path := range args {
			fileInfo, err := os.Stat(path)
			if err != nil {
				return err
			} else if fileInfo.IsDir() {
				return fmt.Errorf("%s is a directory", path)
			}

			// Prevent a file name starting with "-" from being parsed as a flag by the fixer tools
			if strings.HasPrefix(path, "-") {
				path = "./" + path
			}

			errs := []error{fixers.FixWhitespace(path), fixers.FixClrf(path)}

			switch filepath.Ext(path) {
			case ".go":
				errs = append(errs, fixers.FixGofmt(path), fixers.FixGoimports(path))
			case ".json":
				errs = append(errs, fixers.FixJq(path))
			case ".php":
				errs = append(errs, fixers.FixPhpcsfixer(path))
			case ".ts", ".vue", ".js":
				errs = append(errs, fixers.FixEslint(path))
			}

			if err := errors.Join(errs...); err != nil {
				failed = append(failed, fmt.Errorf("%s: %w", path, err))
			}
		}

		return errors.Join(failed...)
	},
}

func init() {
	lintCmd.AddCommand(fixLintCmd)
}
