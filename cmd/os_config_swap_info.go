package cmd

import (
	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var osConfigSwapInfoCmd = &cobra.Command{
	Use:     "info",
	Aliases: []string{"in", "info"},
	Short:   "Show HAOS swap settings",
	Long: `
This command allows you to see how swap is used by the Home Assistant OS.`,
	Example: `
  ha os config swap info`,
	ValidArgsFunction: cobra.NoFileCompletions,
	Args:              cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("os config swap info")

		section := "os"
		command := "config/swap"

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
	osConfigSwapCmd.AddCommand(osConfigSwapInfoCmd)
}
