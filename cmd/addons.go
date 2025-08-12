package cmd

import (
	"strings"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var addonsCmd = &cobra.Command{
	Use:     "addons",
	Aliases: []string{"addon", "add-on", "add-ons", "ad"},
	Short:   "Install, update, remove and configure Home Assistant add-ons",
	Long: `
The addons command allows you to manage Home Assistant add-ons by exposing
commands for installing, removing, configure and control them. It also provides
information commands for add-ons.`,
	Example: `
  ha addons logs core_ssh
  ha addons install core_ssh
  ha addons start core_ssh`,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("addons")

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
	log.Debug("Init addons")

	rootCmd.AddCommand(addonsCmd)
}

func addonsCompletions(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
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
