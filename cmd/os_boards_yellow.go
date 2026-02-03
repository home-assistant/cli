package cmd

import (
	"log/slog"

	helper "github.com/home-assistant/cli/client"
	"github.com/spf13/cobra"
)

var osBoardsYellowCmd = &cobra.Command{
	Use:     "yellow",
	Aliases: []string{"yell"},
	Short:   "See or change settings of the current Yellow board",
	Long: `
This command allows you to see or change settings of the Yellow board that Home
Assistant is running on.`,
	Example: `
  ha os boards yellow`,
	ValidArgsFunction: cobra.NoFileCompletions,
	Args:              cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		slog.Debug("os boards yellow", "args", args)

		section := "os"
		command := "boards/yellow"

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
	osBoardsCmd.AddCommand(osBoardsYellowCmd)
}
