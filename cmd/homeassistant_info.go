package cmd

import (
	"fmt"
	"net/http"

	helper "github.com/home-assistant/hassio-cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	resty "gopkg.in/resty.v1"
)

// infoCmd represents the info command
var homeassistantInfoCmd = &cobra.Command{
	Use:     "info",
	Aliases: []string{"in"},
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("homeassistant info")

		section := "homeassistant"
		command := "info"
		base := viper.GetString("endpoint")

		url, err := helper.URLHelper(base, section, command)
		if err != nil {
			// TODO: error handler
			fmt.Printf("Error: %v", err)
			return
		}

		request := helper.GetClient()
		resp, err := request.Get(url)

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

		if !resty.IsJSONType(resp.Header().Get(http.CanonicalHeaderKey("Content-Type"))) {
			// TODO: return error
			fmt.Println("Error: api did not return a json response")
			return
		}

		helper.ShowJSONResponse(resp.Body())

		return
	},
}
