package cmd

import (
	"fmt"

	helper "github.com/home-assistant/hassio-cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// infoCmd represents the info command
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
			// TODO: error handler
			fmt.Printf("Error: %v", err)
			return
		}

		request := helper.GetClient()
		resp, err := request.SetHeader("Accept", "text/plain").Get(url)

		// explore response object
		log.WithFields(log.Fields{
			"statuscode":  resp.StatusCode(),
			"status":      resp.Status(),
			"time":        resp.Time(),
			"recieved-at": resp.ReceivedAt(),
			"headers":     resp.Header(),
			"request":     resp.Request.RawRequest,
			"body":        resp,
		}).Debug("Response")

		fmt.Print(resp.String())
		return
	},
}
