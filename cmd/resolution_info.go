package cmd

import (
	"log/slog"

	helper "github.com/home-assistant/cli/client"
	"github.com/spf13/cobra"
)

var resolutionInfoCmd = &cobra.Command{
	Use:     "info",
	Aliases: []string{"in", "inf"},
	Short:   "Show issues and suggestions",
	Long: `
This command provides general information about the issues, suggestion and the supported state of the system.`,
	Example: `
  ha resolution info`,
	ValidArgsFunction: cobra.NoFileCompletions,
	Args:              cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		slog.Debug("resolution", "args", args)

		section := "resolution"
		command := "info"

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
	// add cmd to root command
	resolutionCmd.AddCommand(resolutionInfoCmd)
}
