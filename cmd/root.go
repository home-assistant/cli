package cmd

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var debug bool = false
var logFormat string = "text"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "hassio",
	Short: "A brief description of your application",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("hassio")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.homeassistant.yaml)")
	rootCmd.PersistentFlags().StringVarP(&logFormat, "log-format", "", logFormat, "log format to use, valid options are text and json. Default is text")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", debug, "Prints Debug information")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".homeassistant" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".homeassistant.yaml")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
