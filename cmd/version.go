package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	Version   string
	Commit    string
	BuildDate string
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display version and build information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("homectl:\n")
		fmt.Printf("  Version:    %s\n", Version)
		fmt.Printf("  Commit:     %s\n", Commit)
		fmt.Printf("  Build Date: %s\n", BuildDate)
		fmt.Printf("  Platform:   %s (%s)\n", runtime.GOOS, runtime.GOARCH)
		fmt.Printf("  Build Info: %s (%s)\n", runtime.Version(), runtime.Compiler)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
