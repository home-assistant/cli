package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var multicastRestartCmd = &cobra.Command{
	Use:     "restart",
	Aliases: []string{"reboot"},
	Short:   "Restarts the internal Home Assistant Multicast server",
	Long:    `Restart the internal Home Assistant Multicast server running`,
	Example: `
  ha multicast restart`,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("multicast restart")

		section := "multicast"
		command := "restart"
		base := viper.GetString("endpoint")

		ProgressSpinner.Start()
		resp, err := helper.GenericJSONPostTimeout(base, section, command, nil, helper.ContainerOperationTimeout)
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
	multicastCmd.AddCommand(multicastRestartCmd)
}
