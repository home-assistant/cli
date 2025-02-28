package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"strings"
)

var hostLogsBootsCmd = &cobra.Command{
	Use:     "boots",
	Aliases: []string{"list-boots", "lb"},
	Short:   "Show all boot IDs by offset",
	Long: `
Show all values that can be used with the boot arg to find logs.
`,
	Example: `
  ha host logs boots
`,
	ValidArgsFunction: cobra.NoFileCompletions,
	Args:              cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("host logs boots")

		section := "host"
		command := "logs/boots"

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
	hostLogsCmd.AddCommand(hostLogsBootsCmd)
}

func hostBootCompletions(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	resp, err := helper.GenericJSONGet("host/logs/boots", "")
	if err != nil || !resp.IsSuccess() {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	var ret []string
	data := resp.Result().(*helper.Response)
	if data.Result == "ok" && data.Data["boots"] != nil {
		if boots, ok := data.Data["boots"].(map[string]any); ok {
			for bootID, bootName := range boots {
				s := bootName.(string)
				if toComplete == "" || strings.HasPrefix(s, toComplete) {
					ret = append(ret, s+"\tboot offset "+bootID)
				}
			}
		}
	}

	return ret, cobra.ShellCompDirectiveNoFileComp
}
