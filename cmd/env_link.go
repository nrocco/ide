package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var linkEnvCmd = &cobra.Command{
	Use:   "link",
	Short: "Link and executable for this ide project",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return nil
		}

		var source string
		var err error

		parts := strings.SplitN(args[1], ":", 2)

		if len(parts) == 1 {
			source, err = filepath.Abs(args[1])
			if err != nil {
				return err
			}

			fileInfo, err := os.Stat(source)
			if err != nil {
				return err
			}

			if fileInfo.IsDir() {
				return nil // TODO return an error here
			}
		} else {
			source, err = os.Executable()
			if err != nil {
				return err
			}

			Project.AddExecutable(parts[1], parts[0])
		}

		executable := filepath.Join(".git", "bin", args[0])

		log.Printf("Linking %s to %s", executable, source)
		os.Symlink(source, executable)

		return nil
	},
}

func init() {
	envCmd.AddCommand(linkEnvCmd)
}
