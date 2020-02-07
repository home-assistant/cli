package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var coreStatsCmd = &cobra.Command{
	Use:     "stats",
	Aliases: []string{"status", "stat", "st"},
	Short:   "Provides system usage stats of Home Assistant Core",
	Long: `
Provides insight into the system usage stats of Home Assistant Core.
It shows you how much CPU, memory, disk & network resources it uses.`,
	Example: `
  ha core stats`,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("core stats")

		section := "core"
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
	coreCmd.AddCommand(coreStatsCmd)
}
