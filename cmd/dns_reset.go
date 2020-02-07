package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var dnsResetCmd = &cobra.Command{
	Use:     "reset",
	Short:   "Resets the internal Home Assistant DNS server configuration",
	Long:    `Resets the internal Home Assistant DNS server configuration.`,
	Example: `ha dns reset`,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("dns reset")

		section := "dns"
		command := "reset"
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
	dnsCmd.AddCommand(dnsResetCmd)
}
