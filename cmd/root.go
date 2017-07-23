package cmd

import (
	homedir "github.com/mitchellh/go-homedir"
	"github.com/nrocco/ide/pkg/ide"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"strings"
)

var (
	// Version holds the version number of the tpm cli tool
	Version string

	// cfgFile holds the location to the cli configuration file
	cfgFile string

	// Project represents an instance of ide.Project
	Project *ide.Project
)

var rootCmd = &cobra.Command{
	Use:          "ide",
	Short:        "ide provides a powerful ide that gets out of your way",
	Long:         ``,
	SilenceUsage: true,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		Project, err = ide.LoadProject(dir)
		if err != nil {
			log.Fatal(err)
		}
	},
}

// Execute executes the rootCmd logic and is the main entry point for ide
func Execute() error {
	if strings.Contains(os.Args[0], ".git/hooks") {
		elems := strings.Split(os.Args[0], "/")

		args := []string{os.Args[0], "hook", "run", elems[2]}
		os.Args = append(args, os.Args[1:]...)
	} else if !strings.Contains(os.Args[0], "ide") {
		args := []string{os.Args[0], "link", "exec", "--", os.Args[0]}
		os.Args = append(args, os.Args[1:]...)
	}

	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ide.yaml)")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err == nil {
			viper.AddConfigPath(home)
			viper.SetConfigName(".ide")
		}
	}

	viper.SetEnvPrefix("ide")
	viper.AutomaticEnv()

	viper.ReadInConfig()
}
