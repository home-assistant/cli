package cmd

import (
	"log/slog"
	"strings"

	helper "github.com/home-assistant/cli/client"
	"github.com/spf13/cobra"
)

var supervisorOptionsCmd = &cobra.Command{
	Use:     "options",
	Aliases: []string{"option", "opt", "opts", "op"},
	Short:   "Allows you to set options on the Home Assistant Supervisor",
	Long: `
This command allows you to set configuration options for on the Home Assistant
Supervisor running on your Home Assistant system.`,
	Example: `
  ha supervisor options --channel beta`,
	ValidArgsFunction: cobra.NoFileCompletions,
	Args:              cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		slog.Debug("supervisor options", "args", args)

		section := "supervisor"
		command := "options"

		options := make(map[string]any)

		for _, value := range []string{
			"hostname",
			"channel",
			"detect-blocking-io",
			"timezone",
			"logging",
		} {
			val, err := cmd.Flags().GetString(value)
			if val != "" && err == nil && cmd.Flags().Changed(value) {
				options[strings.ReplaceAll(value, "-", "_")] = val
			}
		}

		for _, value := range []string{
			"debug",
			"debug-block",
			"diagnostics",
			"auto-update",
		} {
			data, err := cmd.Flags().GetBool(value)
			if err == nil && cmd.Flags().Changed(value) {
				options[strings.ReplaceAll(value, "-", "_")] = data
			}
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
	supervisorOptionsCmd.Flags().StringP("hostname", "", "", "Hostname to set")
	supervisorOptionsCmd.Flags().StringP("channel", "c", "", "Channel to track (stable|beta|dev)")
	supervisorOptionsCmd.Flags().StringP("detect-blocking-io", "", "", "Detect blocking IO (on|on-at-startup|off)")
	supervisorOptionsCmd.Flags().StringP("timezone", "t", "", "Timezone")
	supervisorOptionsCmd.Flags().StringP("logging", "l", "", "Logging: debug|info|warning|error|critical")
	supervisorOptionsCmd.Flags().BoolP("debug", "", false, "Enable debug mode")
	supervisorOptionsCmd.Flags().BoolP("debug-block", "", false, "Enable debug mode with blocking startup")
	supervisorOptionsCmd.Flags().BoolP("diagnostics", "", false, "Enable diagnostics mode")
	supervisorOptionsCmd.Flags().BoolP("auto-update", "", true, "Enable/disable supervisor auto update")

	supervisorOptionsCmd.Flags().Lookup("debug").NoOptDefVal = "false"
	supervisorOptionsCmd.Flags().Lookup("debug-block").NoOptDefVal = "false"
	supervisorOptionsCmd.Flags().Lookup("diagnostics").NoOptDefVal = "false"
	supervisorOptionsCmd.Flags().Lookup("auto-update").NoOptDefVal = "true"

	supervisorOptionsCmd.RegisterFlagCompletionFunc("hostname", cobra.NoFileCompletions)
	supervisorOptionsCmd.RegisterFlagCompletionFunc("channel", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"stable", "beta", "dev"}, cobra.ShellCompDirectiveNoFileComp
	})
	supervisorOptionsCmd.RegisterFlagCompletionFunc("detect-blocking-io", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"on", "on-at-startup", "off"}, cobra.ShellCompDirectiveNoFileComp
	})
	supervisorOptionsCmd.RegisterFlagCompletionFunc("timezone", cobra.NoFileCompletions)
	supervisorOptionsCmd.RegisterFlagCompletionFunc("logging", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"debug", "info", "warning", "error", "critical"}, cobra.ShellCompDirectiveNoFileComp
	})
	supervisorOptionsCmd.RegisterFlagCompletionFunc("debug", boolCompletions)
	supervisorOptionsCmd.RegisterFlagCompletionFunc("debug-block", boolCompletions)
	supervisorOptionsCmd.RegisterFlagCompletionFunc("diagnostics", boolCompletions)

	supervisorCmd.AddCommand(supervisorOptionsCmd)
}
