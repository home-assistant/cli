package cmd

import (
	"log/slog"

	helper "github.com/home-assistant/cli/client"
	"github.com/spf13/cobra"
)

var dnsStatsCmd = &cobra.Command{
	Use:     "stats",
	Aliases: []string{"status", "stat"},
	Short:   "Provides system usage stats of the Home Assistant DNS server",
	Long: `
Provides insight into the system usage stats of the Home Assistant DNS server.
It shows you how much CPU, memory, disk & network resources it uses.
`,
	Example: `
  ha dns stats
`,
	ValidArgsFunction: cobra.NoFileCompletions,
	Args:              cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		slog.Debug("dns stats", "args", args)

		section := "dns"
		command := "stats"

		resp, err := helper.GenericJSONGet(section, command)
		if err != nil {
			helper.PrintError(err)
			ExitWithError = true
		} else {
			ExitWithError = !helper.ShowJSONResponse(resp)
		}

	},
}

func init() {
	dnsCmd.AddCommand(dnsStatsCmd)
}
