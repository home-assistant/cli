package cmd

import (
	"fmt"

	helper "github.com/home-assistant/hassio-cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var version = ""

// updateCmd represents the update command
var homeassistantUpdateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"up"},
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("homeassistant update")

		section := "homeassistant"
		command := "update"
		base := viper.GetString("endpoint")

		url, err := helper.URLHelper(base, section, command)
		if err != nil {
			// TODO: error handler
			fmt.Printf("Error: %v", err)
			return
		}

		request := helper.GetClient()

		// TODO: submit version
		if version != "" {

			request.SetBody(map[string]interface{}{"version": version})
		}

		resp, err := request.Post(url)

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

func init() {
	fmt.Println("ha update")
	homeassistantUpdateCmd.Flags().StringVarP(&version, "version", "", "", "Version to update to")
	homeassistantCmd.AddCommand(homeassistantUpdateCmd)
}