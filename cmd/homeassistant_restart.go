package cmd

import (
	"fmt"

	helper "github.com/home-assistant/hassio-cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var homeassistantRestartCmd = &cobra.Command{
	Use:     "restart",
	Aliases: []string{"reboot"},
	Short:   "Restarts Home Assistant",
	Long: `
Restart the Home Assistant instance running on your Hass.io system`,
	Example: `
  hassio homeassistant restart`,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("homeassistant restart")

		section := "homeassistant"
		command := "restart"
		base := viper.GetString("endpoint")

		ProgressSpinner.Start()
		resp, err := helper.GenericJSONPost(base, section, command, nil)
		ProgressSpinner.Stop()
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
	homeassistantCmd.AddCommand(homeassistantRestartCmd)
}
