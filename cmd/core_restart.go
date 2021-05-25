package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var coreRestartCmd = &cobra.Command{
	Use:     "restart",
	Aliases: []string{"reboot"},
	Short:   "Restarts the Home Assistant Core",
	Long: `
Restart the Home Assistant Core instance running on your system`,
	Example: `
  ha core restart`,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("core restart")

		section := "core"
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
	coreCmd.AddCommand(coreRestartCmd)
}
