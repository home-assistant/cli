package cmd

import (
	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var storeRepositoriesRepairCmd = &cobra.Command{
	Use:     "repair [slug]",
	Aliases: []string{"reset"},
	Short:   "Repair/reset repository from Home Assistant store",
	Long: `
Repair/reset a repository of add-ons that is missing from store, showing
incorrect information, or otherwise working incorrectly.
`,
	Example: `
ha store repair 94cfad5a
`,
	ValidArgsFunction: storeRepositoriesCompletions,
	Args:              cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("store repair")

		section := "store"
		command := "repositories/{slug}/repair"

		url, err := helper.URLHelper(section, command)
		if err != nil {
			helper.PrintError(err)
			ExitWithError = true
			return
		}

		request := helper.GetJSONRequest()

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
	storeCmd.AddCommand(storeRepositoriesRepairCmd)
}
