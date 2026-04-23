package cmd

import (
	"fmt"
	"log/slog"
	"strconv"
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
  ha supervisor options --channel beta
  ha supervisor options --feature-flag supervisor_v2_api=true`,
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

		featureFlags, err := cmd.Flags().GetStringArray("feature-flag")
		if err == nil && len(featureFlags) > 0 {
			flagMap := make(map[string]bool)
			for _, entry := range featureFlags {
				parts := strings.SplitN(entry, "=", 2)
				name := strings.TrimSpace(parts[0])
				if name == "" {
					helper.PrintError(fmt.Errorf("invalid feature flag %q: name must not be empty", entry))
					ExitWithError = true
					return
				}

				var val bool
				if len(parts) == 1 {
					val = true
				} else {
					val, err = strconv.ParseBool(parts[1])
					if err != nil {
						helper.PrintError(fmt.Errorf("invalid value for feature flag %q: %q, expected true or false", name, parts[1]))
						ExitWithError = true
						return
					}
				}
				flagMap[name] = val
			}
			options["feature_flags"] = flagMap
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
	supervisorOptionsCmd.Flags().StringArrayP("feature-flag", "", []string{}, "Set a development feature flag (name=true|false). Use multiple times for multiple flags.")

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
	supervisorOptionsCmd.RegisterFlagCompletionFunc("feature-flag", featureFlagCompletions)

	supervisorCmd.AddCommand(supervisorOptionsCmd)
}

func featureFlagCompletions(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	resp, err := helper.GenericJSONGet("supervisor", "info")
	if err != nil || !resp.IsSuccess() {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	data := resp.Result().(*helper.Response)
	if data.Result != "ok" || data.Data["feature_flags"] == nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	flags, ok := data.Data["feature_flags"].(map[string]any)
	if !ok {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	var ret []string
	for name, val := range flags {
		current, ok := val.(bool)
		if !ok {
			continue
		}
		ret = append(ret, fmt.Sprintf("%s=%v\tcurrently %v", name, !current, current))
	}
	return ret, cobra.ShellCompDirectiveNoFileComp
}
