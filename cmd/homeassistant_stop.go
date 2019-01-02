package cmd

import (
	"fmt"

	helper "github.com/home-assistant/hassio-cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var homeassistantStopCmd = &cobra.Command{
	Use:     "stop",
	Aliases: []string{},
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("homeassistant stop")

		section := "homeassistant"
		command := "stop"
		base := viper.GetString("endpoint")

		resp, err := helper.GenericJSONPost(base, section, command, nil)
		if err != nil {
			fmt.Println(err)
		} else {
			helper.ShowJSONResponse(resp)
		}

		return
	},
}

func init() {
	homeassistantCmd.AddCommand(homeassistantStopCmd)
}
