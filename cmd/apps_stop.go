package cmd

import (
	"log/slog"

	helper "github.com/home-assistant/cli/client"
	"github.com/spf13/cobra"
)

var appsStopCmd = &cobra.Command{
	Use:     "stop [slug]",
	Aliases: []string{"halt", "shutdown", "quit"},
	Short:   "Manually stop a running Home Assistant app",
	Long: `
This command allows you to manually stop a Home Assistant app
`,
	Example: `
  ha apps stop core_ssh
`,
	ValidArgsFunction: appsCompletions,
	Args:              cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		slog.Debug("apps stop", "args", args)

		section := "addons"
		command := "{slug}/stop"

		url, err := helper.URLHelper(section, command)
		if err != nil {
			helper.PrintError(err)
			ExitWithError = true
			return
		}

		request := helper.GetJSONRequestTimeout(helper.ContainerOperationTimeout)

		slug := args[0]

		request.SetPathParams(map[string]string{
			"slug": slug,
		})

		resp, err := request.Post(url)
		resp, err = helper.GenericJSONErrorHandling(resp, err)

		if err != nil {
			helper.PrintError(err)
			ExitWithError = true
		} else {
			ExitWithError = !helper.ShowJSONResponse(resp)
		}
	},
}

func init() {

	appsCmd.AddCommand(appsStopCmd)
}
