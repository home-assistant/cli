package cmd

import (
	"fmt"
	"strings"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var storeAddonsCmd = &cobra.Command{
	Use:     "addons",
	Aliases: []string{"add-on", "addon", "add-ons"},
	Short:   "Install and update Home Asistant add-ons",
	Long: `
The store command allows you to manage Home Assistant add-ons by exposing
commands for installing or update them.`,
	Example: `
  ha store addons install core_ssh
  ha store addons update core_ssh`,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("store")

		section := "store"
		command := "addons"

		resp, err := helper.GenericJSONGet(section, command)
		if err != nil {
			fmt.Println(err)
			ExitWithError = true
		} else {
			ExitWithError = !helper.ShowJSONResponse(resp)
		}
	},
}

func init() {
	storeCmd.AddCommand(storeAddonsCmd)
}

func storeAddonCompletions(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
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
		if addons, ok := data.Data["addons"].([]interface{}); ok {
			for _, addon := range addons {
				var m map[string]interface{}
				if m, ok = addon.(map[string]interface{}); !ok {
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
