package cmd

import (
	"github.com/spf13/cobra"

	helper "github.com/home-assistant/cli/client"
)

var resolutionIssueCmd = &cobra.Command{
	Use:     "issue",
	Aliases: []string{"is", "trouble"},
	Short:   "Issues management reported by Resolution center",
	Long: `
This command allows dismissing issues reported by the system.`,
	Example: `
  ha resolution issue dismiss [id]`,
}

func init() {
	resolutionCmd.AddCommand(resolutionIssueCmd)
}

func resolutionIssueCompletions(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	resp, err := helper.GenericJSONGet("resolution", "info")
	if err != nil || !resp.IsSuccess() {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	var ret []string
	data := resp.Result().(*helper.Response)
	if data.Result == "ok" && data.Data["issues"] != nil {
		if issues, ok := data.Data["issues"].([]interface{}); ok {
			for _, issue := range issues {
				var m map[string]interface{}
				if m, ok = issue.(map[string]interface{}); !ok {
					continue
				}
				var s string
				if s, ok = m["uuid"].(string); !ok {
					continue
				}
				if t, ok := m["type"].(string); ok {
					s += "\t" + t
				}
				ret = append(ret, s)
			}
		}
	}
	return ret, cobra.ShellCompDirectiveNoFileComp
}
