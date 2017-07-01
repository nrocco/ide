package cmd

import (
	"fmt"

	"github.com/nrocco/ide/pkg/ide"
	"github.com/spf13/cobra"
)

// ctagsCmd represents the status command
var ctagsCmd = &cobra.Command{
	Use:   "ctags",
	Short: "(re)generate a ctags file for this project",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctags, err := ide.LoadCtags(Project)
		if err != nil {
			return err
		}

		fmt.Printf("ctags file %s\n", ctags.TagsFile)
		return ctags.Refresh()
	},
}

func init() {
	execCmd.AddCommand(ctagsCmd)
}
