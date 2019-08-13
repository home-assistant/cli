package cmd

import (
	"fmt"

	helper "github.com/home-assistant/hassio-cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var dnsInfoCmd = &cobra.Command{
	Use:     "info",
	Aliases: []string{"in", "inf"},
	Short:   "Shows information about the internal Hass.io DNS server",
	Long: `
Shows information about the internally running Hass.io DNS server
`,
	Example: `
  hassio dns info
`,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("dns info")

		section := "dns"
		command := "info"
		base := viper.GetString("endpoint")

		resp, err := helper.GenericJSONGet(base, section, command)
		if err != nil {
			fmt.Println(err)
			ExitWithError = true
		} else {
			ExitWithError = !helper.ShowJSONResponse(resp)
		}
	},
}

func init() {
	dnsCmd.AddCommand(dnsInfoCmd)
}
