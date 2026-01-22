package cmd

import (
	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var appsRebuildForce bool

var appsRebuildCmd = &cobra.Command{
	Use:     "rebuild [slug]",
	Aliases: []string{"rb", "reinstall"},
	Short:   "Rebuild a locally built Home Assistant app",
	Long: `
Most apps provide pre-built images Home Assistant can download and use.
However, some don't. This is usually the case for local or development version
of apps. This command allows you to trigger a rebuild of a locally built app.
`,
	Example: `
  ha apps rebuild local_my_app
  ha apps rebuild local_my_app --force
`,
	ValidArgsFunction: appsCompletions,
	Args:              cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("apps rebuild")

		section := "addons"
		command := "{slug}/rebuild"

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

		if appsRebuildForce {
			request.SetBody(map[string]interface{}{
				"force": true,
			})
		}

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
	appsRebuildCmd.Flags().BoolVar(&appsRebuildForce, "force", false, "Force rebuild of the app even if pre-built images are provided")
	appsCmd.AddCommand(appsRebuildCmd)
}
