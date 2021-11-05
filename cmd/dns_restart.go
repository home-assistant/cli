package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var dnsRestartCmd = &cobra.Command{
	Use:     "restart",
	Aliases: []string{"reboot"},
	Short:   "Restarts the internal Home Assistant DNS server",
	Long:    `Restart the internal Home Assistant DNS server running`,
	Example: `
  ha dns restart`,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("dns restart")

		section := "dns"
		command := "restart"

		ProgressSpinner.Start()
		resp, err := helper.GenericJSONPostTimeout(section, command, nil, helper.ContainerOperationTimeout)
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
	dnsCmd.AddCommand(dnsRestartCmd)
}
