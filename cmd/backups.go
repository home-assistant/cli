package cmd

import (
	"log/slog"
	"os"
	"strings"

	helper "github.com/home-assistant/cli/client"
	"github.com/spf13/cobra"
)

var backupsFolders = []string{"addons", "media", "share", "ssl"}

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
		slog.Debug("backups", "args", args)

		section := "backups"
		command := "info"

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
	slog.Debug("Init backups")
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
		if backups, ok := data.Data["backups"].([]any); ok {
			for _, backup := range backups {
				var m map[string]any
				if m, ok = backup.(map[string]any); !ok {
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

func backupsLocationsCompletions(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	resp, err := helper.GenericJSONGet("mounts", "")
	if err != nil || !resp.IsSuccess() {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	var ret []string
	ret = append(ret, ".local\tLocal storage, /backups")
	data := resp.Result().(*helper.Response)
	if data.Result == "ok" && data.Data["mounts"] != nil {
		if mounts, ok := data.Data["mounts"].([]any); ok {
			for _, mount := range mounts {
				var m map[string]any
				if m, ok = mount.(map[string]any); !ok {
					continue
				}
				var s string
				if s, ok = m["usage"].(string); !ok || s != "backup" {
					continue
				}
				if s, ok = m["state"].(string); !ok || s != "active" {
					continue
				}
				if s, ok = m["name"].(string); !ok {
					continue
				}
				ret = append(ret, s)
				var ds []string
				if s, ok = m["server"].(string); ok && s != "" {
					ds = append(ds, s)
				}
				if s, ok = m["share"].(string); ok && s != "" {
					ds = append(ds, s)
				}
				if s, ok = m["path"].(string); ok && s != "" {
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

func backupsAppsCompletions(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	ret, directive := appsCompletions(cmd, args, toComplete)
	ret = append(ret, "ALL\tAll currently installed apps")
	return ret, directive
}

func backupsFoldersCompletions(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	return backupsFolders, cobra.ShellCompDirectiveNoFileComp
}
