package cmd

import (
	"log"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/nrocco/ide/pkg/ide"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Version holds the version number of the tpm cli tool
	Version string

	// cfgFile holds the location to the cli configuration file
	cfgFile string

	// Project represents an instance of ide.Project
	Project ide.Project
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "ide",
	Short: "ide provides a powerful ide that gets out of your way",
	Long:  ``,
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

// init initializes the application
func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ide.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			viper.AddConfigPath(home)
			viper.SetConfigName(".ide")
		}
	}

	viper.SetEnvPrefix("ide")
	viper.AutomaticEnv()

	viper.ReadInConfig()
}
