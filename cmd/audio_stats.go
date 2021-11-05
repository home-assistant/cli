package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var audioStatsCmd = &cobra.Command{
	Use:     "stats",
	Aliases: []string{"status", "stat", "st"},
	Short:   "Provides system usage stats of Home Assistant Audio",
	Long: `
Provides insight into the system usage stats of Home Assistant Audio.
It shows you how much CPU, memory, disk & network resources it uses.`,
	Example: `
  ha audio stats`,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("audio stats")

		section := "audio"
		command := "stats"

		resp, err := helper.GenericJSONGet(section, command)
		if err != nil {
			fmt.Println(err)
			ExitWithError = true
		} else {
			ExitWithError = !helper.ShowJSONResponse(resp)
		}

	},
}

func init() {
	audioCmd.AddCommand(audioStatsCmd)
}
