package cmd

import (
	"fmt"
	"strings"

	humanize "github.com/dustin/go-humanize"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get the current status of your ide project",
	Long:  "Get the current status of your ide project",
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return nil, cobra.ShellCompDirectiveNoFileComp
	},
	PreRunE: loadProject,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Ide\n")
		fmt.Printf("  Name: %s\n", project.Name())
		fmt.Printf("  Branch: %s\n", project.Branch())
		fmt.Printf("  Email: %s\n", project.Email())
		fmt.Printf("  Location: %s\n", project.Location())
		fmt.Printf("  Ctags:\n")
		fmt.Printf("    File: %s\n", project.CtagsFile())
		fmt.Printf("    Age: %s\n", humanize.Time(project.CtagsFileAge()))
		fmt.Printf("    Size: %s\n", humanize.Bytes(project.CtagsFileSize()))

		if project.HasDirEnv() {
			fmt.Printf("  HasDirEnv: yes\n")
		} else {
			fmt.Printf("  HasDirEnv: no\n")
		}

		if project.HasGitBinInPath() {
			fmt.Printf("  GitBinInPath: yes\n")
		} else {
			fmt.Printf("  GitBinInPath: no\n")
		}

		if hooks := project.ListHooks(); len(hooks) > 0 {
			fmt.Printf("  Hooks: %s\n", strings.Join(hooks, " "))
		} else {
			fmt.Printf("  Hooks: ~\n")
		}

		if binaries := project.ListBinaries(); len(binaries) > 0 {
			fmt.Printf("  Binaries:\n")
			for binary, command := range binaries {
				fmt.Printf("    %s: %s\n", binary, command)
			}
		} else {
			fmt.Printf("  Binaries: ~\n")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
