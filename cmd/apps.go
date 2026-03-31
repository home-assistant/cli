package cmd

import (
	"log/slog"
	"os"
	"strings"

	helper "github.com/home-assistant/cli/client"
	"github.com/spf13/cobra"
)

var appsCmd = &cobra.Command{
	Use:     "apps",
	Aliases: []string{"app", "addons", "addon", "add-on", "add-ons", "ad"},
	Short:   "Install, update, remove and configure Home Assistant apps",
	Long: `
The apps command allows you to manage Home Assistant apps by exposing
commands for installing, removing, configure and control them. It also provides
information commands for apps.`,
	Example: `
  ha apps logs core_ssh
  ha apps install core_ssh
  ha apps start core_ssh`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		for idx, arg := range os.Args {
			if idx != 0 && (arg == "addons" || arg == "addon" || arg == "add-on" || arg == "add-ons" || arg == "ad") {
				cmd.PrintErrf("The use of '%s' is deprecated, please use 'apps' instead!\n", arg)
			}
		}
		rootCmd.PersistentPreRun(cmd, args)
	},
	Run: func(cmd *cobra.Command, args []string) {
		slog.Debug("apps", "args", args)

		section := "addons"
		command := ""

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
	slog.Debug("Init apps")

	rootCmd.AddCommand(appsCmd)
}

func appsCompletions(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	resp, err := helper.GenericJSONGet("addons", "")
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
				var s, t string
				if s, ok = m["slug"].(string); !ok {
					continue
				}
				var b bool
				switch cmd.Name() {
				case "rebuild":
					if b, ok = m["build"].(bool); ok && !b {
						continue
					}
				case "start":
					if t, ok = m["state"].(string); ok && t == "started" {
						continue
					}
				case "stop":
					if t, ok = m["state"].(string); ok && t == "stopped" {
						continue
					}
				case "update":
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
