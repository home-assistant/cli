package cmd

import (
	"github.com/spf13/cobra"

	helper "github.com/home-assistant/cli/client"
)

var resolutionCheckCmd = &cobra.Command{
	Use:     "check",
	Aliases: []string{"checks", "test", "che", "ch"},
	Short:   "Check management by the Resolution center",
	Long: `
This command allows to manage checks that are run by the system.`,
	Example: `
  ha resolution check options [slug]`,
}

func init() {
	resolutionCmd.AddCommand(resolutionCheckCmd)
}

func resolutionCheckCompletions(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	resp, err := helper.GenericJSONGet("resolution", "info")
	if err != nil || !resp.IsSuccess() {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	var ret []string
	data := resp.Result().(*helper.Response)
	if data.Result == "ok" && data.Data["checks"] != nil {
		if checks, ok := data.Data["checks"].([]any); ok {
			for _, check := range checks {
				var m map[string]any
				if m, ok = check.(map[string]any); !ok {
					continue
				}
				var s string
				if s, ok = m["slug"].(string); !ok {
					continue
				}
				ret = append(ret, s)
			}
		}
	}
	return ret, cobra.ShellCompDirectiveNoFileComp
}
