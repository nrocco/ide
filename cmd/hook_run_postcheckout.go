package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var runPostCheckoutHookCmd = &cobra.Command{
	Use:   "post-checkout",
	Short: "Run the post-checkout hook for the ide project",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Printf("Writing ctags file to %s\n", Project.GetCtagsFile())
		errCtags := Project.RefreshCtags()
		if errCtags != nil {
			return errCtags
		}

		log.Printf("Writing ctrlp cache to %s\n", Project.GetCtrlpCachFile())
		errCtrlp := Project.RefreshCtrlp()
		if errCtrlp != nil {
			return errCtrlp
		}

		return nil
	},
}

func init() {
	runHookCmd.AddCommand(runPostCheckoutHookCmd)
}
