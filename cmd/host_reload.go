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

// reloadCmd represents the reload command
var hostReloadCmd = &cobra.Command{
	Use:     "reload",
	Aliases: []string{"re"},
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("host reload")

		section := "host"
		command := "reload"
		base := viper.GetString("endpoint")

		url, err := helper.URLHelper(base, section, command)
		if err != nil {
			// TODO: error handler
			fmt.Printf("Error: %v", err)
			return
		}

		request := helper.GetJSONRequest()
		resp, err := request.Post(url)

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
	hostCmd.AddCommand(hostReloadCmd)
}
