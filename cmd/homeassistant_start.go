package cmd

import (
	"fmt"

	helper "github.com/home-assistant/hassio-cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var homeassistantStartCmd = &cobra.Command{
	Use:     "start",
	Aliases: []string{"run", "st"},
	Short:   "Manually start Home Assistant",
	Long: `
This command allows you to manually start the Home Assistant instance on your
Hass.io system. This, of course, only applies when it has been stopped.`,
	Example: `
  hassio homeassistant start`,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("homeassistant start")

		section := "homeassistant"
		command := "start"
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
	homeassistantCmd.AddCommand(homeassistantStartCmd)
}
