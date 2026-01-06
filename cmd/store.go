package cmd

import (
	"strings"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var storeCmd = &cobra.Command{
	Use:     "store",
	Aliases: []string{"shop", "stor"},
	Short:   "Install and update Home Assistant add-ons and manage stores",
	Long: `
The store command allows you to manage Home Assistant add-ons by exposing
commands for installing or update them. It also provides functionality
for managing stores that provide later add-ons.`,
	Example: `
  ha store addons install core_ssh
  ha store add https://github.com/home-assistant/addons-example
  ha store delete 94cfad5a 
  ha store reload`,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("store")

		section := "store"
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
	log.Debug("Init store")

	rootCmd.AddCommand(storeCmd)
}

func storeRepositoriesCompletions(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	resp, err := helper.GenericJSONGet("store", "")
	if err != nil || !resp.IsSuccess() {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	var ret []string
	data := resp.Result().(*helper.Response)
	if data.Result == "ok" && data.Data["repositories"] != nil {
		if repos, ok := data.Data["repositories"].([]any); ok {
			for _, repo := range repos {
				var m map[string]any
				if m, ok = repo.(map[string]any); !ok {
					continue
				}
				var s string
				if s, ok = m["slug"].(string); !ok {
					continue
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
