package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display version and build information",
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return nil, cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Ide:\n")
		fmt.Printf("  Version: %s\n", version)
		fmt.Printf("  Commit: %s\n", commit)
		fmt.Printf("  Build Date: %s\n", date)
		fmt.Printf("  Platform: %s (%s)\n", runtime.GOOS, runtime.GOARCH)
		fmt.Printf("  Build Info: %s (%s)\n", runtime.Version(), runtime.Compiler)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
