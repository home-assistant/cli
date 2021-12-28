package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var hostRebootCmd = &cobra.Command{
	Use:     "reboot",
	Aliases: []string{"restart", "rb"},
	Short:   "Reboots the host machine",
	Long: `
Reboot the machine that your Home Assistant is running on.`,
	Example: `
  ha host reboot`,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("host reboot")

		section := "host"
		command := "reboot"

		resp, err := helper.GenericJSONPost(section, command, nil)
		if err != nil {
			fmt.Println(err)
			ExitWithError = true
		} else {
			ExitWithError = !helper.ShowJSONResponse(resp)
		}
	},
}

func init() {
	hostCmd.AddCommand(hostRebootCmd)
}
