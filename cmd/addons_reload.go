package cmd

import (
	"fmt"

	helper "github.com/home-assistant/hassio-cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var addonsReloadCmd = &cobra.Command{
	Use:     "reload",
	Aliases: []string{"refresh", "re"},
	Short:   "Reloads/Refreshes the Hass.io add-on store",
	Long: `
This commands allows you to force a reload/refresh of the Hass.io add-on store.
Using this, you can force the download of the most recent version information
of an add-on. This might be helpful when you know a new version of an add-on
is released, but not yet available as an upgrade on your Hass.io dashboard.
`,
	Example: `
  hassio addons reload
`,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("addons reload")

		section := "addons"
		command := "reload"
		base := viper.GetString("endpoint")

		resp, err := helper.GenericJSONPost(base, section, command, nil)
		if err != nil {
			fmt.Println(err)
			ExitWithError = true
		} else {
			ExitWithError = !helper.ShowJSONResponse(resp)
		}

		return
	},
}

func init() {
	addonsCmd.AddCommand(addonsReloadCmd)
}
