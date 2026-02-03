package cmd

import (
	"log/slog"

	helper "github.com/home-assistant/cli/client"
	"github.com/spf13/cobra"
)

var resolutionHealthCheckCmd = &cobra.Command{
	Use:     "healthcheck",
	Aliases: []string{"evaluate", "run"},
	Short:   "Execute system healthcheck and fixups",
	Long: `
This command execute a full system check and auto fixups. It check all issues again to see if they
are still around and try to fix it again.`,
	Example: `
  ha resolution healthcheck`,
	ValidArgsFunction: cobra.NoFileCompletions,
	Args:              cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		slog.Debug("resolution", "args", args)

		section := "resolution"
		command := "healthcheck"

		ProgressSpinner.Start()
		resp, err := helper.GenericJSONPost(section, command, nil)
		ProgressSpinner.Stop()

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
	resolutionCmd.AddCommand(resolutionHealthCheckCmd)
}
