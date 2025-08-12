package cmd

import (
	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var addonsRebuildForce bool

var addonsRebuildCmd = &cobra.Command{
	Use:     "rebuild [slug]",
	Aliases: []string{"rb", "reinstall"},
	Short:   "Rebuild a locally built Home Assistant add-on",
	Long: `
Most add-ons provide pre-built images Home Assistant can download and use.
However, some don't. This is usually the case for local or development version
of add-ons. This command allows you to trigger a rebuild of a locally built
add-on.
`,
	Example: `
  ha addons rebuild local_my_addon
  ha addons rebuild local_my_addon --force
`,
	ValidArgsFunction: addonsCompletions,
	Args:              cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("addons rebuild")

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

		if addonsRebuildForce {
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
	addonsRebuildCmd.Flags().BoolVar(&addonsRebuildForce, "force", false, "Force rebuild of the add-on even if pre-built images are provided")
	addonsCmd.AddCommand(addonsRebuildCmd)
}
