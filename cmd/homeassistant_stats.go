package cmd

import (
	"fmt"

	helper "github.com/home-assistant/hassio-cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var homeassistantStatsCmd = &cobra.Command{
	Use:     "stats",
	Aliases: []string{"status", "stat", "st"},
	Short:   "Provides system usage stats of Home Assistant",
	Long: `
Provides insight into the system usage stats of Home Assistant.
It shows you how much CPU, memory, disk & network resources it uses.`,
	Example: `
  hassio homeassistant stats`,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("homeassistant stats")

		section := "homeassistant"
		command := "stats"
		base := viper.GetString("endpoint")

		resp, err := helper.GenericJSONGet(base, section, command)
		if err != nil {
			fmt.Println(err)
			ExitWithError = true
		} else {
			ExitWithError = !helper.ShowJSONResponse(resp)
		}

	},
}

func init() {
	homeassistantCmd.AddCommand(homeassistantStatsCmd)
}
