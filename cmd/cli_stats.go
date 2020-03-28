package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cliStatsCmd = &cobra.Command{
	Use:     "stats",
	Aliases: []string{"status", "stat"},
	Short:   "Provides system usage stats of the Home Assistant cli backend",
	Long: `
Provides insight into the system usage stats of the Home Assistant CLI backend.
It shows you how much CPU, memory, disk & network resources it uses.
`,
	Example: `
  ha cli stats
`,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("cli stats")

		section := "cli"
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
	cliCmd.AddCommand(cliStatsCmd)
}
