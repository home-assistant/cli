package cmd

import (
	"log/slog"

	helper "github.com/home-assistant/cli/client"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:     "info",
	Aliases: []string{"in", "inf"},
	Short:   "Provides a general Home Assistant information overview",
	Long: `
This command provide a general information about your Home Assistant system.
The information provide can be useful for sharing when you are encountering
issues or when reporting one on GitHub.
	`,
	Example: `
  ha info
	`,
	ValidArgsFunction: cobra.NoFileCompletions,
	Args:              cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		slog.Debug("info", "args", args)

		section := "info"
		command := ""

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
	rootCmd.AddCommand(infoCmd)
}
