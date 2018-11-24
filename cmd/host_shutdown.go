package cmd

import (
	"fmt"

	helper "github.com/home-assistant/hassio-cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// shutdownCmd represents the shutdown command
var hostShutdownCmd = &cobra.Command{
	Use:     "shutdown",
	Aliases: []string{"sh"},
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("host shutdown")

		section := "host"
		command := "shutdown"
		base := viper.GetString("endpoint")

		resp, err := helper.GenericJSONPost(base, section, command, nil)
		if err != nil {
			fmt.Println(err)
		} else {
			helper.ShowJSONResponse(resp)
		}
		return
	},
}

func init() {
	hostCmd.AddCommand(hostShutdownCmd)
}
