package cmd

import (
	"log/slog"

	helper "github.com/home-assistant/cli/client"
	"github.com/spf13/cobra"
)

var securityIntegrityCmd = &cobra.Command{
	Use:     "integrity",
	Aliases: []string{"int", "trust"},
	Short:   "Execute security integrity check",
	Long: `
This command execute a full system integrity check.
This need content trust to be enabled.`,
	Example: `
  ha security integrity`,
	ValidArgsFunction: cobra.NoFileCompletions,
	Args:              cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		slog.Debug("security integrity", "args", args)

		section := "security"
		command := "integrity"

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
	securityCmd.AddCommand(securityIntegrityCmd)
}
