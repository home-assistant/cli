package cmd

import (
	"log/slog"

	helper "github.com/home-assistant/cli/client"
	"github.com/spf13/cobra"
)

var osBoardsGreenCmd = &cobra.Command{
	Use:     "green",
	Aliases: []string{"grn"},
	Short:   "See or change settings of the current Green board",
	Long: `
This command allows you to see or change settings of the Green board that Home
Assistant is running on.`,
	Example: `
  ha os boards green`,
	ValidArgsFunction: cobra.NoFileCompletions,
	Args:              cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		slog.Debug("os boards green", "args", args)

		section := "os"
		command := "boards/green"

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
	osBoardsCmd.AddCommand(osBoardsGreenCmd)
}
