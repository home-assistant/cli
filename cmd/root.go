package cmd

import (
	"errors"
	"io/fs"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"time"

	helper "github.com/home-assistant/cli/client"
	"github.com/home-assistant/cli/config"
	"github.com/home-assistant/cli/spinner"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var cfgFile string
var endPoint string
var logLevel string
var apiToken string
var hassSource string
var rawJSON bool
var noProgress bool

// ExitWithError is a hint for the called that we want a non-zero exit code
var ExitWithError = false

// ProgressSpinner is a general spinner that can be used across the CLI
var ProgressSpinner = spinner.New(spinner.CharSets[0], 125*time.Millisecond)

var rootCmd = &cobra.Command{
	Use:   path.Base(os.Args[0]),
	Short: "A small CLI program to control Home Assistant",
	Long: `
The Home Assistant CLI is a small and simple command line utility that allows
you to control and configure different aspects of Home Assistant`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// set loglevel if possible
		var level slog.Level
		if level.UnmarshalText([]byte(config.GetString("log-level"))) == nil {
			slog.SetLogLoggerLevel(level)
		}

		helper.RawJSON = config.GetBool("raw-json") && cmd.Name() != cobra.ShellCompRequestCmd

		// Only shows spinner output when we have a TTY
		if !noProgress && term.IsTerminal(int(os.Stdout.Fd())) {
			// Write to Stderr, helps when redirecting output, e.g., to a file
			ProgressSpinner.Writer = os.Stderr
		}

		slog.Debug("Debug flags",
			"apiToken", config.GetString("api-token"),
			"cfgFile", config.GetString("config"),
			"endpoint", config.GetString("endpoint"),
			"hassSource", config.GetString("hass-source"),
			"logLevel", config.GetString("log-level"),
			"noProgress", config.GetBool("no-progress"),
			"rawJSON", config.GetBool("raw-json"),
		)
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if ProgressSpinner.Active() {
			ProgressSpinner.Stop()
		}
	},
}

// Execute represents the entrypoint for when called without any subcommand
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		helper.PrintError(err)
		ExitWithError = true
	}
}

func init() {

	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Optional config file (default is $HOME/.homeassistant.yaml)")
	rootCmd.PersistentFlags().StringVar(&endPoint, "endpoint", "", "Endpoint for Home Assistant Supervisor, use unix:///path for a Unix socket (default is $SUPERVISOR_SOCKET if set, else 'supervisor')")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "", "Log level (defaults to "+slog.LevelWarn.String()+")")
	rootCmd.PersistentFlags().StringVar(&apiToken, "api-token", "", "Home Assistant Supervisor API token")
	rootCmd.PersistentFlags().StringVar(&hassSource, "hass-source", "", "Identify the caller to Supervisor in the X-Hass-Source header (not sent by default)")
	rootCmd.PersistentFlags().BoolVar(&rawJSON, "raw-json", false, "Output raw JSON from the API")
	rootCmd.PersistentFlags().BoolVar(&noProgress, "no-progress", false, "Disable the progress spinner")

	rootCmd.RegisterFlagCompletionFunc("endpoint", cobra.NoFileCompletions)
	rootCmd.RegisterFlagCompletionFunc("log-level", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{
			slog.LevelDebug.String(),
			slog.LevelInfo.String(),
			slog.LevelWarn.String(),
			slog.LevelError.String(),
		}, cobra.ShellCompDirectiveNoFileComp
	})
	rootCmd.RegisterFlagCompletionFunc("api-token", cobra.NoFileCompletions)
	rootCmd.RegisterFlagCompletionFunc("hass-source", cobra.NoFileCompletions)
	rootCmd.RegisterFlagCompletionFunc("raw-json", boolCompletions)
	rootCmd.RegisterFlagCompletionFunc("no-progress", boolCompletions)
	rootCmd.MarkPersistentFlagFilename("config", "yaml", "yml")

	config.BindFlag(rootCmd.PersistentFlags().Lookup("config"))
	config.BindFlag(rootCmd.PersistentFlags().Lookup("endpoint"))
	config.BindFlag(rootCmd.PersistentFlags().Lookup("log-level"))
	config.BindFlag(rootCmd.PersistentFlags().Lookup("api-token"))
	config.BindFlag(rootCmd.PersistentFlags().Lookup("hass-source"))
	config.BindFlag(rootCmd.PersistentFlags().Lookup("raw-json"))
	config.BindFlag(rootCmd.PersistentFlags().Lookup("no-progress"))

	config.SetDefault("endpoint", defaultEndpoint())
	config.SetDefault("log-level", "warn")
	config.SetDefault("api-token", "")

	// Configure global spinner
	ProgressSpinner.Suffix = " Processing..."
	ProgressSpinner.FinalMSG = "Processing... Done.\n\n"
}

// defaultEndpoint returns the endpoint used when none is configured. Supervisor
// gives apps a socket of their own in SUPERVISOR_SOCKET, prefer it over TCP.
func defaultEndpoint() string {
	if socket := os.Getenv("SUPERVISOR_SOCKET"); socket != "" {
		return helper.UnixSocketScheme + socket
	}
	return "supervisor"
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// bind to current API token ENV variable
	config.BindEnv("api-token", "SUPERVISOR_TOKEN")

	// set loglevel if possible
	var level slog.Level
	if level.UnmarshalText([]byte(config.GetString("log-level"))) == nil {
		slog.SetLogLoggerLevel(level)
	}

	// Use config file from the flag, or search the home directory for
	// ".homeassistant" with one of the supported extensions.
	candidates := []string{cfgFile}
	if cfgFile == "" {
		candidates = nil
		// Without a home directory, e.g. if $HOME is not set, there is no
		// config file to look for
		if home, err := os.UserHomeDir(); err == nil {
			slog.Debug("Adding homedir to searchpath", "homedir", home)
			for _, ext := range []string{"yaml", "yml", "json"} {
				candidates = append(candidates, filepath.Join(home, ".homeassistant."+ext))
			}
		}
	}

	// If a config file is found, read it in.
	for _, candidate := range candidates {
		err := config.ReadFile(candidate)
		if err == nil {
			slog.Info("Using configfile", "configfile", candidate)
			return
		}
		if !errors.Is(err, fs.ErrNotExist) {
			slog.Info("Cannot read configfile", "configfile", candidate, "error", err)
			return
		}
	}
	slog.Info("No configfile found")
}

func boolCompletions(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{"true", "false"}, cobra.ShellCompDirectiveNoFileComp
}
