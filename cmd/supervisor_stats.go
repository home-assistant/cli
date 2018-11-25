package cmd

import (
	"fmt"

	helper "github.com/home-assistant/hassio-cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var supervisorStatsCmd = &cobra.Command{
	Use:     "stats",
	Aliases: []string{"st"},
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("supervisor stats")

		section := "supervisor"
		command := "stats"
		base := viper.GetString("endpoint")

		resp, err := helper.GenericJSONGet(base, section, command)
		if err != nil {
			fmt.Println(err)
		} else {
			helper.ShowJSONResponse(resp)
		}

	},
}

func init() {
	supervisorCmd.AddCommand(supervisorStatsCmd)
}
