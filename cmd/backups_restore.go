package cmd

import (
	"log/slog"

	helper "github.com/home-assistant/cli/client"
	"github.com/spf13/cobra"
)

var backupsRestoreCmd = &cobra.Command{
	Use:   "restore [slug]",
	Short: "Restores a Home Assistant backup",
	Long: `
When something goes wrong, this command allows you to restore a previously
take Home Assistant backup on your system.`,
	Example: `
  ha backups restore c1a07617
  ha backups restore c1a07617 --app core_ssh --app core_mosquitto
  ha backups restore c1a07617 --folders homeassistant`,
	ValidArgsFunction: backupsCompletions,
	Args:              cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		slog.Debug("backups restore", "args", args)

		section := "backups/{slug}"
		command := "restore/full"

		request := helper.GetJSONRequestTimeout(helper.BackupTimeout)

		options := make(map[string]any)

		slug := args[0]

		request.SetPathParams(map[string]string{
			"slug": slug,
		})

		password, err := cmd.Flags().GetString("password")
		if password != "" && err == nil && cmd.Flags().Changed("password") {
			options["password"] = password
		}

		homeassistant, err := cmd.Flags().GetBool("homeassistant")
		if err == nil && cmd.Flags().Changed("homeassistant") {
			options["homeassistant"] = homeassistant
			command = "restore/partial"
		}

		apps, err := cmd.Flags().GetStringArray("app")
		addonsDeprecated, _ := cmd.Flags().GetStringArray("addons")
		apps = append(apps, addonsDeprecated...)
		slog.Debug("apps", "apps", apps)

		if len(apps) > 0 && err == nil {
			options["addons"] = apps
			command = "restore/partial"
		}

		folders, err := cmd.Flags().GetStringArray("folders")
		slog.Debug("folders", "folders", folders)

		if len(folders) > 0 && err == nil {
			options["folders"] = folders
			command = "restore/partial"
		}

		location, err := cmd.Flags().GetString("location")
		if err == nil && cmd.Flags().Changed("location") {
			options["location"] = location
		}

		url, err := helper.URLHelper(section, command)
		if err != nil {
			helper.PrintError(err)
			ExitWithError = true
			return
		}

		if len(options) > 0 {
			slog.Debug("Request body", "options", options)
			request.SetBody(options)
		}

		ProgressSpinner.Start()
		resp, err := request.Post(url)
		ProgressSpinner.Stop()

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
	backupsRestoreCmd.Flags().StringP("password", "", "", "Password")
	backupsRestoreCmd.Flags().BoolP("homeassistant", "", true, "Restore homeassistant (default true), performs a partial restore when set to false")
	backupsRestoreCmd.Flags().StringArrayP("app", "a", []string{}, "App to restore, performs a partial restore. Use multiple times for multiple apps.")
	backupsRestoreCmd.Flags().StringArray("addons", []string{}, "")
	backupsRestoreCmd.Flags().MarkHidden("addons")
	backupsRestoreCmd.Flags().MarkDeprecated("addons", "use --app instead")
	backupsRestoreCmd.Flags().StringArrayP("folders", "f", []string{}, "Folders to restore, performs a partial restore")
	backupsRestoreCmd.Flags().StringP("location", "l", "", "Where to put the backup file (backup mount or local)")

	backupsRestoreCmd.Flags().Lookup("location").NoOptDefVal = ".local"

	backupsRestoreCmd.RegisterFlagCompletionFunc("password", cobra.NoFileCompletions)
	backupsRestoreCmd.RegisterFlagCompletionFunc("homeassistant", boolCompletions)
	backupsRestoreCmd.RegisterFlagCompletionFunc("app", backupsAppsCompletions)
	backupsRestoreCmd.RegisterFlagCompletionFunc("folders", backupsFoldersCompletions)
	backupsRestoreCmd.RegisterFlagCompletionFunc("location", backupsLocationsCompletions)

	backupsCmd.AddCommand(backupsRestoreCmd)
}
