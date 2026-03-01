package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Remove all ide configuration for a repository",
	Long:  "Remove all ide configuration for a repository",
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return nil, cobra.ShellCompDirectiveNoFileComp
	},
	PreRunE: loadProject,
	RunE: func(cmd *cobra.Command, args []string) error {
		destroyForce, _ := cmd.Flags().GetBool("no-private")

		if !destroyForce {
			fmt.Print("This will remove all ide configuration from the repository. Continue? [y/N] ")
			reader := bufio.NewReader(os.Stdin)
			response, err := reader.ReadString('\n')
			if err != nil {
				return err
			}
			if strings.ToLower(strings.TrimSpace(response)) != "y" {
				fmt.Println("Aborted")
				return nil
			}
		}

		if err := project.Destroy(); err != nil {
			return err
		}

		fmt.Println("Repository is no longer an ide project")

		return nil
	},
}

func init() {
	destroyCmd.Flags().Bool("force", false, "Actually destroy ide configuration")

	rootCmd.AddCommand(destroyCmd)
}
