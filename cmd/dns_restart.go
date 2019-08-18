package cmd

import (
	"fmt"

	helper "github.com/home-assistant/hassio-cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var dnsRestartCmd = &cobra.Command{
	Use:     "restart",
	Aliases: []string{"reboot"},
	Short:   "Restarts the internal Hass.io DNS server",
	Long: `
Restart the internal Hass.io DNS server running on your Hass.io system`,
	Example: `
  hassio dns restart`,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("dns restart")

		section := "dns"
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
	dnsCmd.AddCommand(dnsRestartCmd)
}
