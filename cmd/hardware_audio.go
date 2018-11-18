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
var hardwareAudioCmd = &cobra.Command{
	Use:     "audio",
	Aliases: []string{"au"},
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("hardware info")

		section := "hardware"
		command := "audio"
		base := viper.GetString("endpoint")

		url, err := helper.URLHelper(base, section, command)
		if err != nil {
			// TODO: error handler
			fmt.Printf("Error: %v", err)
			return
		}

		request := helper.GetJSONRequest()
		resp, err := request.Get(url)

		if !resty.IsJSONType(resp.Header().Get(http.CanonicalHeaderKey("Content-Type"))) {
			// TODO: return error
			fmt.Println("Error: api did not return a json response")
			return
		}
		helper.ShowJSONResponse(resp)
		return
	},
}

func init() {
	hardwareCmd.AddCommand(hardwareAudioCmd)
}
