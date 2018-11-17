package cmd

import (
	"fmt"

	helper "github.com/home-assistant/hassio-cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// stopCmd represents the stop command
var homeassistantStopCmd = &cobra.Command{
	Use:     "stop",
	Aliases: []string{"ch"},
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("homeassistant stop")

		section := "homeassistant"
		command := "stop"
		base := viper.GetString("endpoint")

		url, err := helper.URLHelper(base, section, command)
		if err != nil {
			// TODO: error handler
			fmt.Printf("Error: %v", err)
			return
		}

		request := helper.GetClient()
		resp, err := request.Post(url)

		// returns 200 OK or 400
		if resp.StatusCode() != 200 && resp.StatusCode() != 400 {
			fmt.Println("Unexpected server response")
			fmt.Println(resp.String())
		} else {
			helper.ShowJSONResponse(resp.Body())
		}

		return
	},
}
