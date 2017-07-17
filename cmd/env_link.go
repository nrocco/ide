package cmd

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var linkEnvCmd = &cobra.Command{
	Use:   "link",
	Short: "Link and executable for this ide project",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return nil
		}

		// TODO: check format of 2nd argument
		//		 if php or /usr/bin/php5 create symlink to binary on host system
		//       if container_name:php or container_name:/usr/bin/php5 create symlink to ide
		// TODO: if latter, add git config php = container_name

		executable := filepath.Join(".git", "bin", args[0])
		source, _ := filepath.Abs(args[1])

		log.Printf("Linking %s to %s", executable, source)
		os.Symlink(source, executable)

		return nil
	},
}

func init() {
	envCmd.AddCommand(linkEnvCmd)
}
