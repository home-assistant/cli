package cmd

import (
	"log/slog"

	helper "github.com/home-assistant/cli/client"
	"github.com/spf13/cobra"
)

var backupsRemoveCmd = &cobra.Command{
	Use:     "remove [slug]",
	Aliases: []string{"delete", "del", "rem", "rm"},
	Short:   "Deletes a backup from disk",
	Long: `
Backups can take quite a bit of diskspace, this command allows you to
clean backups from disk.`,
	Example: `
  ha backups remove c1a07617`,
	ValidArgsFunction: backupsCompletions,
	Args:              cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		slog.Debug("backups remove", "args", args)

		section := "backups"
		command := "{slug}"

		url, err := helper.URLHelper(section, command)
		if err != nil {
			helper.PrintError(err)
			ExitWithError = true
			return
		}

		request := helper.GetJSONRequest()
		options := make(map[string]any)

		slug := args[0]
		request.SetPathParams(map[string]string{
			"slug": slug,
		})

		location, err := cmd.Flags().GetStringArray("location")
		slog.Debug("location", "location", location)
		if len(location) > 0 && err == nil && cmd.Flags().Changed(("location")) {
			options["location"] = location
		}

		if len(options) > 0 {
			slog.Debug("Request body", "options", options)
			request.SetBody(options)
		}

		resp, err := request.Delete(url)
		resp, err = helper.GenericJSONErrorHandling(resp, err)

		if err != nil {
			helper.PrintError(err)
			ExitWithError = true
		} else {
			ExitWithError = !helper.ShowJSONResponse(resp)
		}
	},
}

func init() {
	backupsRemoveCmd.Flags().StringArrayP("location", "l", []string{}, "location(s) to remove backup from (instead of all), use multiple times for multiple locations.")
	backupsRemoveCmd.Flags().Lookup("location").NoOptDefVal = ".local"
	backupsRemoveCmd.RegisterFlagCompletionFunc("location", backupsLocationsCompletions)

	backupsCmd.AddCommand(backupsRemoveCmd)
}
