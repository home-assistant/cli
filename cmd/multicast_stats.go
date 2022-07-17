package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var multicastStatsCmd = &cobra.Command{
	Use:     "stats",
	Aliases: []string{"status", "stat"},
	Short:   "Provides system usage stats of the Home Assistant Multicast server",
	Long: `
Provides insight into the system usage stats of the Home Assistant Multicast server.
It shows you how much CPU, memory, disk & network resources it uses.
`,
	Example: `
  ha multicast stats
`,
	ValidArgsFunction: cobra.NoFileCompletions,
	Args:              cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("multicast stats")

		section := "multicast"
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
	multicastCmd.AddCommand(multicastStatsCmd)
}
