package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/nrocco/ide/pkg/ide"
	"github.com/spf13/cobra"
)

var parseCmdOpts struct {
	Debug     bool
	NoPublic  bool
	NoPrivate bool
}

var parseCmd = &cobra.Command{
	Use:   "parse <path>...",
	Short: "Parse ctags output for a path",
	Args:  cobra.MinimumNArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return nil, cobra.ShellCompDirectiveDefault
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		currentFile := ""

		return project.CtagsParseCode(func(entry ide.CtagsEntry) {
			if parseCmdOpts.NoPublic && entry.IsPublic() {
				return
			}
			if parseCmdOpts.NoPrivate && entry.IsPrivate() {
				return
			}
			if parseCmdOpts.Debug {
				fmt.Printf("%+v\n", entry)
				return
			}
			if currentFile != entry.Path {
				if currentFile != "" {
					fmt.Printf("\n\n")
				}
				currentFile = entry.Path
				color.Magenta(entry.Path)
				fmt.Printf("\n")
			}
			fmt.Print(color.HiBlackString("%5d  ", entry.Line))
			if entry.Language == "python" && entry.Kind == "member" {
				fmt.Printf("%s.", entry.Scope)
			}
			if entry.Language == "go" && entry.ScopeKind == "struct" {
				fmt.Printf("(%s) ", entry.Scope)
			}
			if entry.IsPublic() {
				fmt.Print(color.HiBlueString("%s", entry.Name))
			} else {
				fmt.Print(color.HiRedString("%s", entry.Name))
			}
			fmt.Print(color.HiGreenString("%s", entry.Signature))
			if entry.Typeref != "" {
				fmt.Print(color.HiBlackString(": "))
				fmt.Print(color.HiBlueString("%s", entry.Typeref))
			}
			fmt.Printf("\n")
		}, args...)
	},
}

func init() {
	parseCmd.Flags().BoolVar(&parseCmdOpts.Debug, "debug", false, "Output raw JSON instead of formatted output")
	parseCmd.Flags().BoolVar(&parseCmdOpts.NoPublic, "no-public", false, "Exclude public functions")
	parseCmd.Flags().BoolVar(&parseCmdOpts.NoPrivate, "no-private", false, "Exclude private functions")

	toolCmd.AddCommand(parseCmd)
}
