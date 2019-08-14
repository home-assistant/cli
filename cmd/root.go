package cmd

import (
	"os"
	"path"
	"strings"
	"time"
	"io/ioutil"

	"github.com/home-assistant/hassio-cli/client"
	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/briandowns/spinner"
)

var cfgFile string
var endPoint string
var logLevel string
var apiToken string
var rawJSON bool

// ExitWithError is a hint for the called that we want an non-zero exit code
var ExitWithError = false

// ProgressSpinner is a general spinner that can be used across the CLI
var ProgressSpinner = spinner.New(spinner.CharSets[11], 100*time.Millisecond)

var rootCmd = &cobra.Command{
	Use:   path.Base(os.Args[0]),
	Short: "A small CLI program to control Hass.io",
	Long: `
The Hass.io CLI is a small and simple command line utility that allows you to
control and configure different aspects of Hass.io`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {

		// set loglevel if posible
		logrusLevel, err := log.ParseLevel(viper.GetString("log-level"))

		if err == nil {
			log.SetLevel(logrusLevel)
		}

		client.RawJSON = viper.GetBool("raw-json")

		log.WithFields(log.Fields{
			"cfgFile":  viper.GetString("config"),
			"endpoint": viper.GetString("endpoint"),
			"logLevel": viper.GetString("log-level"),
			"apiToken": viper.GetString("api-token"),
			"rawJSON":  viper.GetBool("raw-json"),
		}).Debugln("Debug flags")
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if (ProgressSpinner.Active()) {
			ProgressSpinner.Stop()
		}
	},	
}

// Execute represents the entrypoint for when called without any subcommand
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error while executing rootCmd: %s", err)
	}
}

func init() {

	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Optional config file (default is $HOME/.homeassistant.yaml)")
	rootCmd.PersistentFlags().StringVar(&endPoint, "endpoint", "", "Endpoint for Hass.io Supervisor ( default is 'hassio' )")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "", "Log level (defaults to Warn)")
	rootCmd.PersistentFlags().StringVar(&apiToken, "api-token", "", "Hass.io API token")
	rootCmd.PersistentFlags().BoolVar(&rawJSON, "raw-json", false, "Output raw JSON from the API")

	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
	viper.BindPFlag("endpoint", rootCmd.PersistentFlags().Lookup("endpoint"))
	viper.BindPFlag("log-level", rootCmd.PersistentFlags().Lookup("log-level"))
	viper.BindPFlag("api-token", rootCmd.PersistentFlags().Lookup("api-token"))
	viper.BindPFlag("raw-json", rootCmd.PersistentFlags().Lookup("raw-json"))

	viper.SetDefault("endpoint", "hassio")
	viper.SetDefault("log-level", "warn")
	viper.SetDefault("api-token", "")

	// Configure global spinner
	ProgressSpinner.Suffix = " Processing..."
	ProgressSpinner.FinalMSG = "Processing... Done.\n\n"
	ProgressSpinner.Writer = ioutil.Discard

	// Only shows spinner output when not in a pipe
	info, _ := os.Stdin.Stat()
	if (info.Mode() & os.ModeCharDevice) == os.ModeCharDevice {
		ProgressSpinner.Writer = os.Stderr
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetEnvPrefix("HASSIO")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	// bind to current API token ENV variable
	viper.BindEnv("api-token", "HASSIO_TOKEN")

	// set loglevel if possible
	logLevel, err := log.ParseLevel(viper.GetString("log-level"))

	if err == nil {
		log.SetLevel(logLevel)
	}

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			log.Fatalf("Error while finding home directory: %s", err)
		}

		// Search config in home directory with name ".homeassistant" (without extension).
		viper.AddConfigPath(home)
		log.WithField("homedir", home).Debug("Adding homedir to searchpath")
		viper.SetConfigName(".homeassistant")
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.WithField("configfile", viper.ConfigFileUsed()).Info("Using configfile")
	} else {
		log.Info("No configfile found")
	}
}
