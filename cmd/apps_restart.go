package cmd

import (
	"log/slog"

	helper "github.com/home-assistant/cli/client"
	"github.com/spf13/cobra"
)

var appsRestartCmd = &cobra.Command{
	Use:     "restart [slug]",
	Aliases: []string{"reboot"},
	Short:   "Restarts a Home Assistant app",
	Long: `
Restart a Home Assistant app
`,
	Example: `
  ha apps restart core_ssh
`,
	ValidArgsFunction: appsCompletions,
	Args:              cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		slog.Debug("apps restart", "args", args)

		section := "addons"
		command := "{slug}/restart"

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

	appsCmd.AddCommand(appsRestartCmd)
}
