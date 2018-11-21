package cmd

import (
	"fmt"

	helper "github.com/home-assistant/hassio-cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// reloadCmd represents the reload command
var supervisorReloadCmd = &cobra.Command{
	Use:     "reload",
	Aliases: []string{"ch"},
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("supervisor reload")

		section := "supervisor"
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
	supervisorCmd.AddCommand(supervisorReloadCmd)
}
