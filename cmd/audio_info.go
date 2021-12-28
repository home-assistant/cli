package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var audioInfoCmd = &cobra.Command{
	Use:     "info",
	Aliases: []string{"in", "inf"},
	Short:   "Provides information about Home Assistant Audio devices",
	Long: `
This command provides information about the running Home Assistant Audio instance
running on your Home Assistant system, including its devices.`,
	Example: `
  ha audio info`,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("audio info")

		section := "audio"
		command := "info"

		resp, err := helper.GenericJSONGet(section, command)
		if err != nil {
			fmt.Println(err)
			ExitWithError = true
		} else {
			ExitWithError = !helper.ShowJSONResponse(resp)
		}
	},
}

func init() {
	audioCmd.AddCommand(audioInfoCmd)
}
