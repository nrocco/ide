package cmd

import (
	"os"
	"path/filepath"

	"github.com/nrocco/ide/pkg/ide/tools"
	"github.com/spf13/cobra"
)

var fixCmd = &cobra.Command{
	Use:   "fix",
	Short: "Fix source code",
	Long:  "Fix source code",
	RunE: func(cmd *cobra.Command, args []string) error {
		for _, path := range args {
			if _, err := os.Stat(path); os.IsNotExist(err) {
				return err
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
			case ".ts":
				tools.FixWhitespace(path)
				tools.FixClrf(path)
			case ".yaml":
				tools.FixWhitespace(path)
				tools.FixClrf(path)
			case ".yml":
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
	rootCmd.AddCommand(fixCmd)
}
