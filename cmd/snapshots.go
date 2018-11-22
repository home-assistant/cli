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

// snapshotsCmd represents the snapshots command
var snapshotsCmd = &cobra.Command{
	Use:     "snapshots",
	Aliases: []string{"sn"},
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("snapshots")

		section := "snapshots"
		command := ""
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
	log.Debug("Init snapshots")
	// add subcommands
	// TODO: add subcommand

	// add cmd to root command
	rootCmd.AddCommand(snapshotsCmd)
}
