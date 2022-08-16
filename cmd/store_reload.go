package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var storeReloadCmd = &cobra.Command{
	Use:     "reload",
	Aliases: []string{"refresh", "re"},
	Short:   "Reloads/Refreshes the Home Assistant add-on store",
	Long: `
This commands allows you to force a reload/refresh of the Home Assistant add-on
store. Using this, you can force the download of the most recent version
information of an add-on. This might be helpful when you know a new version of
an add-on is released, but not yet available as an upgrade in Home Assistant.
`,
	Example: `
  ha store reload
`,
	ValidArgsFunction: cobra.NoFileCompletions,
	Args:              cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("store reload")

		section := "store"
		command := "reload"

		ProgressSpinner.Start()
		resp, err := helper.GenericJSONPost(section, command, nil)
		ProgressSpinner.Stop()

		if err != nil {
			fmt.Println(err)
			ExitWithError = true
		} else {
			ExitWithError = !helper.ShowJSONResponse(resp)
		}
	},
}

func init() {
	storeCmd.AddCommand(storeReloadCmd)
}
