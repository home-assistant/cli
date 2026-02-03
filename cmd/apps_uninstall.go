package cmd

import (
	"log/slog"

	helper "github.com/home-assistant/cli/client"
	"github.com/spf13/cobra"
)

var appsUninstallCmd = &cobra.Command{
	Use:     "uninstall [slug]",
	Aliases: []string{"remove", "delete", "del", "rem", "un", "uninst"},
	Short:   "Uninstalls a Home Assistant app",
	Long: `
This command allows you to uninstall a Home Assistant app.
`,
	Example: `
  ha apps uninstall core_ssh
`,
	ValidArgsFunction: appsCompletions,
	Args:              cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		slog.Debug("apps uninstall", "args", args)

		section := "addons"
		command := "{slug}/uninstall"

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

		var options map[string]any

		removeConfig, _ := cmd.Flags().GetBool("remove-config")
		if removeConfig {
			options = map[string]any{"remove_config": removeConfig}
			request.SetBody(options)
		}

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
	appsUninstallCmd.Flags().Bool("remove-config", false, "Delete app's config folder (if used)")
	appsUninstallCmd.Flags().Lookup("remove-config").NoOptDefVal = "true"
	appsUninstallCmd.RegisterFlagCompletionFunc("remove-config", boolCompletions)

	appsCmd.AddCommand(appsUninstallCmd)
}
