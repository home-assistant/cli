package cmd

import (
	"log/slog"

	helper "github.com/home-assistant/cli/client"
	"github.com/spf13/cobra"
)

var supervisorUpdateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"upgrade", "downgrade", "up", "down"},
	Short:   "Updates the Home Assistant Supervisor",
	Long: `
Using this command you can upgrade or downgrade the Home Assistant Supervisor
running on your Home Assistant  system to the latest version
or the version specified.`,
	Example: `
  ha supervisor update
  ha supervisor update --version 173`,
	ValidArgsFunction: cobra.NoFileCompletions,
	Args:              cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		slog.Debug("supervisor update", "args", args)

		section := "supervisor"
		command := "update"

		var options map[string]any

		version, _ := cmd.Flags().GetString("version")
		if version != "" {
			options = map[string]any{"version": version}
		}

		ProgressSpinner.Start()
		resp, err := helper.GenericJSONPostTimeout(section, command, options, helper.ContainerDownloadTimeout)
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
	supervisorUpdateCmd.Flags().StringP("version", "", "", "Version to update to")
	supervisorUpdateCmd.RegisterFlagCompletionFunc("version", cobra.NoFileCompletions)
	supervisorCmd.AddCommand(supervisorUpdateCmd)
}
