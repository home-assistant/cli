package cmd

import (
	"log/slog"

	helper "github.com/home-assistant/cli/client"
	"github.com/spf13/cobra"
)

var backupFreezeCmd = &cobra.Command{
	Use:     "freeze",
	Aliases: []string{"frz"},
	Short:   "Freeze supervisor for external backup",
	Long: `
This command tells Supervisor to prepare Home Assistant and apps for a backup
or snapshot taken by external software. Caller should call thaw when done.`,
	Example: `
  ha backups freeze --timeout 300`,
	ValidArgsFunction: cobra.NoFileCompletions,
	Args:              cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		slog.Debug("backups freeze", "args", args)

		section := "backups"
		command := "freeze"

		options := make(map[string]any)

		timeout, err := cmd.Flags().GetInt("timeout")
		if timeout != 0 && err == nil && cmd.Flags().Changed("timeout") {
			options["timeout"] = timeout
		}

		resp, err := helper.GenericJSONPost(section, command, options)
		if err != nil {
			helper.PrintError(err)
			ExitWithError = true
		} else {
			ExitWithError = !helper.ShowJSONResponse(resp)
		}
	},
}

func init() {
	backupFreezeCmd.Flags().Int("timeout", 0, "Seconds before freeze times out and thaw begins")
	backupFreezeCmd.RegisterFlagCompletionFunc("timeout", cobra.NoFileCompletions)
	backupsCmd.AddCommand(backupFreezeCmd)
}
