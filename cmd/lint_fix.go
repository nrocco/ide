package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nrocco/ide/pkg/ide/tools"
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
				tools.FixWhitespace(path)
				tools.FixClrf(path)
				tools.FixGofmt(path)
				tools.FixGoimports(path)
			case ".html":
				tools.FixWhitespace(path)
				tools.FixClrf(path)
			case ".json":
				tools.FixWhitespace(path)
				tools.FixClrf(path)
				tools.FixJq(path)
			case ".php":
				tools.FixWhitespace(path)
				tools.FixClrf(path)
				tools.FixPhpcsfixer(path)
			case ".py":
				tools.FixWhitespace(path)
				tools.FixClrf(path)
			case ".rb":
				tools.FixWhitespace(path)
				tools.FixClrf(path)
			case ".sh":
				tools.FixWhitespace(path)
				tools.FixClrf(path)
			case ".ts", ".vue", ".js":
				tools.FixWhitespace(path)
				tools.FixClrf(path)
				tools.FixEslint(path)
			case ".yaml", ".yml":
				tools.FixWhitespace(path)
				tools.FixClrf(path)
			default:
				tools.FixWhitespace(path)
				tools.FixClrf(path)
			}
		}

		return nil
	},
}

func init() {
	lintCmd.AddCommand(fixLintCmd)
}
