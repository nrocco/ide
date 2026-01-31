package cmd

import (
	"fmt"
	"strconv"

	"github.com/nrocco/ide/pkg/ide"
	"github.com/spf13/cobra"
)

var genpassOpts ide.PasswordOptions

var genpassCmd = &cobra.Command{
	Use:   "genpass [length]",
	Short: "Generate a random password",
	Long:  "Generate a random password with configurable character sets",
	Args:  cobra.MaximumNArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return nil, cobra.ShellCompDirectiveNoFileComp
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		length := 12
		if len(args) > 0 {
			var err error
			length, err = strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid length: %s", args[0])
			}
		}

		password, err := ide.GeneratePassword(length, genpassOpts)
		if err != nil {
			return err
		}

		fmt.Println(password)
		return nil
	},
}

func init() {
	genpassCmd.Flags().BoolVar(&genpassOpts.Upper, "upper", true, "include uppercase letters")
	genpassCmd.Flags().BoolVar(&genpassOpts.Lower, "lower", true, "include lowercase letters")
	genpassCmd.Flags().BoolVar(&genpassOpts.Numbers, "numbers", true, "include digits")
	genpassCmd.Flags().BoolVar(&genpassOpts.Punctuation, "punctuation", false, "include punctuation")
	genpassCmd.Flags().BoolVar(&genpassOpts.Hexdigits, "hexdigits", false, "use only hex digits")

	toolsCmd.AddCommand(genpassCmd)
}
