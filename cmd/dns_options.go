package cmd

import (
	"fmt"

	helper "github.com/home-assistant/hassio-cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var dnsOptionsCmd = &cobra.Command{
	Use:     "options",
	Aliases: []string{"op"},
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("dns options")

		section := "dns"
		command := "options"
		base := viper.GetString("endpoint")

		options := make(map[string]interface{})

		servers, err := cmd.Flags().GetStringArray("servers")
		log.WithField("servers", servers).Debug("servers")

		if len(servers) >= 1 && err == nil {
			options["servers"] = servers
		}

		resp, err := helper.GenericJSONPost(base, section, command, options)
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
	dnsOptionsCmd.Flags().StringArrayP("servers", "r", []string{}, "Upstream DNS servers to use. Use multiple times for multiple servers.")
	dnsCmd.AddCommand(dnsOptionsCmd)
}
