package cmd

import (
	"log/slog"

	helper "github.com/home-assistant/cli/client"
	"github.com/spf13/cobra"
)

var backupsNewCmd = &cobra.Command{
	Use:     "new",
	Aliases: []string{"create", "backup"},
	Short:   "Create a new Home Assistant backup",
	Long: `
This command can be used to trigger the creation of a new Home Assistant
backup.`,
	Example: `
  ha backups new
  ha backups new --app core_ssh --app core_mosquitto
  ha backups new --folders homeassistant
  ha backups new --uncompressed
`,
	ValidArgsFunction: cobra.NoFileCompletions,
	Args:              cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		slog.Debug("backups new", "args", args)

		section := "backups"
		command := "new/full"

		options := make(map[string]any)

		name, err := cmd.Flags().GetString("name")
		if name != "" && err == nil && cmd.Flags().Changed("name") {
			options["name"] = name
		}

		password, err := cmd.Flags().GetString("password")
		if password != "" && err == nil && cmd.Flags().Changed("password") {
			options["password"] = password
		}

		apps, err := cmd.Flags().GetStringArray("app")
		addonsDeprecated, _ := cmd.Flags().GetStringArray("addons")
		apps = append(apps, addonsDeprecated...)
		slog.Debug("apps", "apps", apps)

		if len(apps) != 0 && err == nil && (cmd.Flags().Changed("app") || cmd.Flags().Changed("addons")) {
			options["addons"] = apps
			command = "new/partial"
		}

		folders, err := cmd.Flags().GetStringArray("folders")
		slog.Debug("folders", "folders", folders)

		if len(folders) != 0 && err == nil && cmd.Flags().Changed("folders") {
			options["folders"] = folders
			command = "new/partial"
		}

		if cmd.Flags().Changed("uncompressed") {
			options["compressed"] = false
		}

		location, err := cmd.Flags().GetStringArray("location")
		slog.Debug("location", "location", location)
		if len(location) > 0 && err == nil && cmd.Flags().Changed("location") {
			options["location"] = location
		}

		filename, err := cmd.Flags().GetString("filename")
		slog.Debug("filename", "filename", filename)
		if filename != "" && err == nil && cmd.Flags().Changed("filename") {
			options["filename"] = filename
		}

		ExcludeDB, err := cmd.Flags().GetBool("homeassistant-exclude-database")
		if err == nil && cmd.Flags().Changed("homeassistant-exclude-database") {
			options["homeassistant_exclude_database"] = ExcludeDB
		}

		ProgressSpinner.Start()
		resp, err := helper.GenericJSONPostTimeout(section, command, options, helper.BackupTimeout)
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
	backupsNewCmd.Flags().StringP("name", "", "", "Name of the backup")
	backupsNewCmd.Flags().StringP("password", "", "", "Password")
	backupsNewCmd.Flags().Bool("uncompressed", false, "Use uncompressed archives")
	backupsNewCmd.Flags().StringArrayP("app", "a", []string{}, "App to backup, triggers a partial backup. Use multiple times for multiple apps.")
	backupsNewCmd.Flags().StringArray("addons", []string{}, "")
	backupsNewCmd.Flags().MarkHidden("addons")
	backupsNewCmd.Flags().MarkDeprecated("addons", "use --app instead")
	backupsNewCmd.Flags().StringArrayP("folders", "f", []string{}, "Folders to backup, triggers a partial backup")
	backupsNewCmd.Flags().StringArrayP("location", "l", []string{}, "Where to put backup file (backup mount or local). Use multiple times for multiple locations.")
	backupsNewCmd.Flags().Bool("homeassistant-exclude-database", false, "Exclude the Home Assistant database file from backup")
	backupsNewCmd.Flags().String("filename", "", "Name to use for the backup file")

	backupsNewCmd.Flags().Lookup("uncompressed").NoOptDefVal = "false"
	backupsNewCmd.Flags().Lookup("location").NoOptDefVal = ".local"
	backupsNewCmd.Flags().Lookup("homeassistant-exclude-database").NoOptDefVal = "false"

	backupsNewCmd.RegisterFlagCompletionFunc("name", cobra.NoFileCompletions)
	backupsNewCmd.RegisterFlagCompletionFunc("password", cobra.NoFileCompletions)
	backupsNewCmd.RegisterFlagCompletionFunc("uncompressed", boolCompletions)
	backupsNewCmd.RegisterFlagCompletionFunc("app", backupsAppsCompletions)
	backupsNewCmd.RegisterFlagCompletionFunc("folders", backupsFoldersCompletions)
	backupsNewCmd.RegisterFlagCompletionFunc("location", backupsLocationsCompletions)
	backupsNewCmd.RegisterFlagCompletionFunc("homeassistant-exclude-database", boolCompletions)
	backupsNewCmd.RegisterFlagCompletionFunc("filename", cobra.NoFileCompletions)

	backupsCmd.AddCommand(backupsNewCmd)
}
