package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var hostShutdownCmd = &cobra.Command{
	Use:     "shutdown",
	Aliases: []string{"sh"},
	Short:   "Shutdown the host machine",
	Long: `
Shuts down the machine that your Home Assistant is running on.
WARNING: This is turning off the computer/device.`,
	Example: `
  ha host shutdown`,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("host shutdown")

		section := "host"
		command := "shutdown"

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
	hostCmd.AddCommand(hostShutdownCmd)
}
