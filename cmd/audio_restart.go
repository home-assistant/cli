package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var audioRestartCmd = &cobra.Command{
	Use:     "restart",
	Aliases: []string{"reboot"},
	Short:   "Restarts the internal Home Assistant Audio",
	Long:    `Restart the internal Home Assistant Audio`,
	Example: `
  ha audio restart`,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("audio restart")

		section := "audio"
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
	audioCmd.AddCommand(audioRestartCmd)
}
