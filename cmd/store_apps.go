package cmd

import (
	"log/slog"
	"os"
	"strings"

	helper "github.com/home-assistant/cli/client"
	"github.com/spf13/cobra"
)

var storeAppsCmd = &cobra.Command{
	Use:     "apps",
	Aliases: []string{"app", "addons", "add-on", "addon", "add-ons"},
	Short:   "Install and update Home Assistant apps",
	Long: `
The store command allows you to manage Home Assistant apps by exposing
commands for installing or update them.`,
	Example: `
  ha store apps install core_ssh
  ha store apps update core_ssh`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		for idx, arg := range os.Args {
			if idx != 0 && (arg == "addons" || arg == "addon" || arg == "add-on" || arg == "add-ons") {
				cmd.PrintErrf("The use of '%s' is deprecated, please use 'apps' instead!\n", arg)
			}
		}
		rootCmd.PersistentPreRun(cmd, args)
	},
	Run: func(cmd *cobra.Command, args []string) {
		slog.Debug("store apps", "args", args)

		section := "store"
		command := "addons"

		resp, err := helper.GenericJSONGet(section, command)
		if err != nil {
			helper.PrintError(err)
			ExitWithError = true
		} else {
			ExitWithError = !helper.ShowJSONResponse(resp)
		}
	},
}

func init() {
	storeCmd.AddCommand(storeAppsCmd)
}

func storeAppCompletions(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	resp, err := helper.GenericJSONGet("store", "")
	if err != nil || !resp.IsSuccess() {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	var ret []string
	data := resp.Result().(*helper.Response)
	if data.Result == "ok" && data.Data["addons"] != nil {
		if addons, ok := data.Data["addons"].([]any); ok {
			for _, addon := range addons {
				var m map[string]any
				if m, ok = addon.(map[string]any); !ok {
					continue
				}
				var s string
				if s, ok = m["slug"].(string); !ok {
					continue
				}
				var b bool
				switch cmd.Name() {
				case "install":
					if b, ok = m["available"].(bool); ok && !b {
						continue
					}
					if b, ok = m["installed"].(bool); ok && b {
						continue
					}
				case "update":
					if b, ok = m["available"].(bool); ok && !b {
						continue
					}
					if b, ok = m["installed"].(bool); ok && !b {
						continue
					}
					if b, ok = m["update_available"].(bool); ok && !b {
						continue
					}
				}
				ret = append(ret, s)
				var ds []string
				if s, ok = m["name"].(string); ok && s != "" {
					ds = append(ds, s)
				}
				if s, ok = m["url"].(string); ok && s != "" {
					ds = append(ds, s)
				}
				if len(ds) != 0 {
					ret[len(ret)-1] += "\t" + strings.Join(ds, ", ")
				}
			}
		}
	}
	return ret, cobra.ShellCompDirectiveNoFileComp
}
