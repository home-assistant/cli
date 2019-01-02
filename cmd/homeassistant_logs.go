package cmd

import (
	"fmt"

	helper "github.com/home-assistant/hassio-cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var homeassistantLogsCmd = &cobra.Command{
	Use:     "logs",
	Aliases: []string{"lo"},
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("homeassistant logs")

		section := "homeassistant"
		command := "logs"
		base := viper.GetString("endpoint")

		url, err := helper.URLHelper(base, section, command)
		if err != nil {
			fmt.Printf("Error: %v", err)
			return
		}

		request := helper.GetRequest()
		resp, err := request.SetHeader("Accept", "text/plain").Get(url)

		fmt.Println(resp.String())
		return
	},
}

func init() {
	homeassistantCmd.AddCommand(homeassistantLogsCmd)
}
