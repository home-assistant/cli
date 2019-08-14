package cmd

import (
	"fmt"

	helper "github.com/home-assistant/hassio-cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var homeassistantUpdateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"upgrade", "downgrade", "up", "down"},
	Short:   "Updates Home Assistant",
	Long: `
Using this command you can upgrade or downgrade the Home Assistant instance
running on your Hass.io system to the latest version or the version specified.`,
	Example: `
  hassio homeassistant update
  hassio homeassistant update --version 0.97.2`,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("homeassistant update")

		section := "homeassistant"
		command := "update"
		base := viper.GetString("endpoint")

		var options map[string]interface{}

		version, err := cmd.Flags().GetString("version")
		if version != "" {
			options = map[string]interface{}{"version": version}
		}

		ProgressSpinner.Start()
		resp, err := helper.GenericJSONPost(base, section, command, options)
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
	homeassistantUpdateCmd.Flags().StringP("version", "", "", "Version to update to")
	homeassistantCmd.AddCommand(homeassistantUpdateCmd)
}
