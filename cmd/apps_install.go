package cmd

import (
	"log/slog"

	helper "github.com/home-assistant/cli/client"
	"github.com/spf13/cobra"
)

var appsInstallCmd = &cobra.Command{
	Use:     "install [slug]",
	Aliases: []string{"i", "inst"},
	Short:   "Installs a Home Assistant app",
	Long: `
This command allows you to install a Home Assistant app from the commandline.
`,
	Example: `
  ha apps install core_ssh
`,
	ValidArgsFunction: storeAppCompletions,
	Args:              cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		slog.Debug("apps install", "args", args)

		section := "addons"
		command := "{slug}/install"

		url, err := helper.URLHelper(section, command)
		if err != nil {
			helper.PrintError(err)
			ExitWithError = true
			return
		}

		request := helper.GetJSONRequestTimeout(helper.ContainerDownloadTimeout)

		slug := args[0]

		request.SetPathParams(map[string]string{
			"slug": slug,
		})

		ProgressSpinner.Start()
		resp, err := request.Post(url)
		ProgressSpinner.Stop()

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

	appsCmd.AddCommand(appsInstallCmd)
}
