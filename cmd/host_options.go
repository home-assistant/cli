package cmd

import (
	"fmt"

	helper "github.com/home-assistant/hassio-cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var hostOptionsCmd = &cobra.Command{
	Use:     "options",
	Aliases: []string{"op"},
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("host options")

		section := "host"
		command := "options"
		base := viper.GetString("endpoint")

		var options map[string]interface{}

		hostname, err := cmd.Flags().GetString("hostname")
		if hostname != "" {
			options = map[string]interface{}{"hostname": hostname}
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
	hostOptionsCmd.Flags().StringP("hostname", "", "", "Hostname to set")
	hostCmd.AddCommand(hostOptionsCmd)
}
