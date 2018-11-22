package cmd

import (
	"fmt"

	helper "github.com/home-assistant/hassio-cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var hostHostname = ""

// optionsCmd represents the options command
var hostOptionsCmd = &cobra.Command{
	Use:     "options",
	Aliases: []string{"op"},
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("host options")

		section := "host"
		command := "options"
		base := viper.GetString("endpoint")

		url, err := helper.URLHelper(base, section, command)
		if err != nil {
			// TODO: error handler
			fmt.Printf("Error: %v", err)
			return
		}

		request := helper.GetJSONRequest()

		// TODO: submit hostname
		if hostHostname != "" {
			request.SetBody(map[string]interface{}{"hostname": hostHostname})
		}

		resp, err := request.Post(url)

		// returns 200 OK or 400
		if resp.StatusCode() != 200 && resp.StatusCode() != 400 {
			fmt.Println("Unexpected server response")
			fmt.Println(resp.String())
		} else {
			helper.ShowJSONResponse(resp)
		}

		return
	},
}

func init() {
	hostOptionsCmd.Flags().StringVarP(&hostHostname, "hostname", "", "", "Hostname to set")
	hostCmd.AddCommand(hostOptionsCmd)
}
