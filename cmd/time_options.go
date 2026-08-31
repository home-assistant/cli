package cmd

import (
	"log/slog"

	helper "github.com/home-assistant/cli/client"
	"github.com/spf13/cobra"
)

var timeOptionsCmd = &cobra.Command{
	Use:     "options",
	Aliases: []string{"option", "opt", "opts", "op"},
	Short:   "Set Home Assistant OS time options",
	Long: `
This command allows you to set Home Assistant OS NTP servers and fallback NTP
servers.`,
	Example: `
  ha time options --servers pool.ntp.org --servers time.cloudflare.com
  ha time options --fallback-servers time.google.com
  ha time options --clear-fallback-servers`,
	ValidArgsFunction: cobra.NoFileCompletions,
	Args:              cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		slog.Debug("time options", "args", args)

		section := "time"
		command := "options"

		options := make(map[string]any)

		servers, err := cmd.Flags().GetStringArray("servers")
		if err == nil && cmd.Flags().Changed("servers") {
			options["servers"] = servers
		}

		clearServers, err := cmd.Flags().GetBool("clear-servers")
		if err == nil && clearServers {
			options["servers"] = []string{}
		}

		fallbackServers, err := cmd.Flags().GetStringArray("fallback-servers")
		if err == nil && cmd.Flags().Changed("fallback-servers") {
			options["fallback_servers"] = fallbackServers
		}

		clearFallbackServers, err := cmd.Flags().GetBool("clear-fallback-servers")
		if err == nil && clearFallbackServers {
			options["fallback_servers"] = []string{}
		}

		resp, err := helper.GenericJSONPost(section, command, options)
		if err != nil {
			helper.PrintError(err)
			ExitWithError = true
		} else {
			ExitWithError = !helper.ShowJSONResponse(resp)
		}
	},
}

func init() {
	timeOptionsCmd.Flags().StringArrayP("servers", "s", []string{}, "NTP servers to use. Use multiple times for multiple servers.")
	timeOptionsCmd.Flags().StringArrayP("fallback-servers", "f", []string{}, "Fallback NTP servers to use. Use multiple times for multiple servers.")
	timeOptionsCmd.Flags().Bool("clear-servers", false, "Clear configured NTP servers")
	timeOptionsCmd.Flags().Bool("clear-fallback-servers", false, "Clear configured fallback NTP servers")

	timeOptionsCmd.RegisterFlagCompletionFunc("servers", cobra.NoFileCompletions)
	timeOptionsCmd.RegisterFlagCompletionFunc("fallback-servers", cobra.NoFileCompletions)
	timeOptionsCmd.RegisterFlagCompletionFunc("clear-servers", boolCompletions)
	timeOptionsCmd.RegisterFlagCompletionFunc("clear-fallback-servers", boolCompletions)

	timeOptionsCmd.MarkFlagsMutuallyExclusive("servers", "clear-servers")
	timeOptionsCmd.MarkFlagsMutuallyExclusive("fallback-servers", "clear-fallback-servers")

	timeCmd.AddCommand(timeOptionsCmd)
}
