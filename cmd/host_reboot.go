package cmd

import (
	"fmt"

	helper "github.com/home-assistant/hassio-cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var hostRebootCmd = &cobra.Command{
	Use:     "reboot",
	Aliases: []string{"restart", "rb"},
	Short:   "Reboots the host machine",
	Long: `
Reboot the machine that your Hass.io is running on.`,
	Example: `
  hassio host reboot`,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("host reboot")

		section := "host"
		command := "reboot"
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
	hostCmd.AddCommand(hostRebootCmd)
}
