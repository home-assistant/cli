package cmd

import (
	"github.com/spf13/cobra"

	helper "github.com/home-assistant/cli/client"
)

var resolutionSuggestionCmd = &cobra.Command{
	Use:     "suggestion",
	Aliases: []string{"su", "solution"},
	Short:   "Suggestion management reported by Resolution center",
	Long: `
This command allows to dismiss or apply suggestion reported by the system.`,
	Example: `
  ha resolution suggestion dismiss [id]
  ha resolution suggestion apply [id]`,
}

func init() {
	resolutionCmd.AddCommand(resolutionSuggestionCmd)
}

func resolutionSuggestionCompletions(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	resp, err := helper.GenericJSONGet("resolution", "info")
	if err != nil || !resp.IsSuccess() {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	var ret []string
	data := resp.Result().(*helper.Response)
	if data.Result == "ok" && data.Data["suggestions"] != nil {
		if suggestions, ok := data.Data["suggestions"].([]interface{}); ok {
			for _, suggestion := range suggestions {
				var m map[string]interface{}
				if m, ok = suggestion.(map[string]interface{}); !ok {
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
