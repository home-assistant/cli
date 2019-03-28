package cmd

import (
	"fmt"

	helper "github.com/home-assistant/hassio-cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var homeassistantRebuildCmd = &cobra.Command{
	Use:     "rebuild",
	Aliases: []string{"re"},
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("homeassistant rebuild")

		section := "homeassistant"
		command := "rebuild"
		base := viper.GetString("endpoint")

		resp, err := helper.GenericJSONPost(base, section, command, nil)
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
	homeassistantCmd.AddCommand(homeassistantRebuildCmd)
}
