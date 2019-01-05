package cmd

import (
	"fmt"

	helper "github.com/home-assistant/hassio-cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var hardwareAudioCmd = &cobra.Command{
	Use:     "audio",
	Aliases: []string{"au"},
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("hardware info")

		section := "hardware"
		command := "audio"
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
	hardwareCmd.AddCommand(hardwareAudioCmd)
}
