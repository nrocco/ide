package cmd

import (
	"fmt"

	"github.com/nrocco/ide/pkg/ide"
	"github.com/spf13/cobra"
)

// ctrlpCmd represents the status command
var ctrlpCmd = &cobra.Command{
	Use:   "ctrlp",
	Short: "(re)generate cache file for the ctrlp vim plugin",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctrlp, err := ide.LoadCtrlp(Project)
		if err != nil {
			return err
		}

		fmt.Printf("Writing ctrlp cache to %s\n", ctrlp.CacheFile)
		return ctrlp.Refresh()
	},
}

func init() {
	execCmd.AddCommand(ctrlpCmd)
}
