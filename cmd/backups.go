package cmd

import (
	"fmt"
	"os"
	"strings"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var backupsCmd = &cobra.Command{
	Use:     "backups",
	Aliases: []string{"backup", "back", "backups", "bk", "snapshots", "snapshot", "snap", "shot", "sn"},
	Short:   "Create, restore and remove backups",
	Long: `
Backups of your Home Assistant system, which you can create,
restore, and delete using this command.`,
	Example: `
  ha backups
  ha backups new`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		for idx, arg := range os.Args {
			if idx != 0 && (arg == "snapshots" || arg == "snapshot" || arg == "snap" || arg == "shot" || arg == "sn") {
				cmd.PrintErrf("The use of '%s' is deprecated, please use 'backups' instead!\n", arg)
			}
		}
		rootCmd.PersistentPreRun(cmd, args)
	},
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("backups")

		section := "backups"
		command := ""

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
	log.Debug("Init backups")
	// add cmd to root command
	rootCmd.AddCommand(backupsCmd)
}

func backupsCompletions(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	resp, err := helper.GenericJSONGet("backups", "")
	if err != nil || !resp.IsSuccess() {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	var ret []string
	data := resp.Result().(*helper.Response)
	if data.Result == "ok" && data.Data["backups"] != nil {
		if backups, ok := data.Data["backups"].([]interface{}); ok {
			for _, backup := range backups {
				var m map[string]interface{}
				if m, ok = backup.(map[string]interface{}); !ok {
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
				if s, ok = m["date"].(string); ok && s != "" {
					ds = append(ds, s)
				}
				if s, ok = m["type"].(string); ok && s != "" {
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
