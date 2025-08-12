package cmd

import (
	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
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
		log.WithField("args", args).Debug("os boards yellow")

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
