package cmd

import (
	"log/slog"

	helper "github.com/home-assistant/cli/client"
	"github.com/spf13/cobra"
)

var storeReloadCmd = &cobra.Command{
	Use:     "reload",
	Aliases: []string{"refresh", "re"},
	Short:   "Reloads/Refreshes the Home Assistant app store",
	Long: `
This command allows you to force a reload/refresh of the Home Assistant app
store. Using this, you can force the download of the most recent version
information of an app. This might be helpful when you know a new version of
an app is released, but not yet available as an upgrade in Home Assistant.
`,
	Example: `
  ha store reload
`,
	ValidArgsFunction: cobra.NoFileCompletions,
	Args:              cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		slog.Debug("store reload", "args", args)

		section := "store"
		command := "reload"

		ProgressSpinner.Start()
		resp, err := helper.GenericJSONPost(section, command, nil)
		ProgressSpinner.Stop()

		if err != nil {
			helper.PrintError(err)
			ExitWithError = true
		} else {
			ExitWithError = !helper.ShowJSONResponse(resp)
		}
	},
}

func init() {
	storeCmd.AddCommand(storeReloadCmd)
}
