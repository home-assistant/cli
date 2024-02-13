package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
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
  ha backups new --addons core_ssh --addons core_mosquitto
  ha backups new --folders homeassistant
  ha backups new --uncompressed
`,
	ValidArgsFunction: cobra.NoFileCompletions,
	Args:              cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("backups new")

		section := "backups"
		command := "new/full"

		options := make(map[string]interface{})

		name, err := cmd.Flags().GetString("name")
		if name != "" && err == nil && cmd.Flags().Changed("name") {
			options["name"] = name
		}

		password, err := cmd.Flags().GetString("password")
		if password != "" && err == nil && cmd.Flags().Changed("password") {
			options["password"] = password
		}

		addons, err := cmd.Flags().GetStringArray("addons")
		log.WithField("addons", addons).Debug("addons")

		if len(addons) != 0 && err == nil && cmd.Flags().Changed("addons") {
			options["addons"] = addons
			command = "new/partial"
		}

		folders, err := cmd.Flags().GetStringArray("folders")
		log.WithField("folders", folders).Debug("folders")

		if len(folders) != 0 && err == nil && cmd.Flags().Changed("folders") {
			options["folders"] = folders
			command = "new/partial"
		}

		if cmd.Flags().Changed("uncompressed") {
			options["compressed"] = false
		}

		location, err := cmd.Flags().GetString("location")
		if err == nil && cmd.Flags().Changed("location") {
			if location == "" {
				options["location"] = nil
			} else {
				options["location"] = location
			}
		}

		ExcludeDB, err := cmd.Flags().GetBool("homeassistant-exclude-database")
		if err == nil && cmd.Flags().Changed("homeassistant-exclude-database") {
			options["homeassistant-exclude-database"] = ExcludeDB
		}

		ProgressSpinner.Start()
		resp, err := helper.GenericJSONPostTimeout(section, command, options, helper.BackupTimeout)
		ProgressSpinner.Stop()
		if err != nil {
			fmt.Println(err)
			ExitWithError = true
		} else {
			ExitWithError = !helper.ShowJSONResponse(resp)
		}
	},
}

func init() {
	backupsNewCmd.Flags().StringP("name", "", "", "Name of the backup")
	backupsNewCmd.Flags().StringP("password", "", "", "Password")
	backupsNewCmd.Flags().Bool("uncompressed", false, "Use Uncompressed archives")
	backupsNewCmd.Flags().StringArrayP("addons", "a", []string{}, "addons to backup, triggers a partial backup")
	backupsNewCmd.Flags().StringArrayP("folders", "f", []string{}, "folders to backup, triggers a partial backup")
	backupsNewCmd.Flags().StringP("location", "l", "", "where to put backup file (backup mount or local)")
	backupsNewCmd.Flags().Bool("homeassistant-exclude-database", false, "Exclude the Home Assistant database file from backup")

	backupsNewCmd.Flags().Lookup("uncompressed").NoOptDefVal = "false"
	backupsNewCmd.Flags().Lookup("location").NoOptDefVal = ""
	backupsNewCmd.Flags().Lookup("homeassistant-exclude-database").NoOptDefVal = "false"

	backupsNewCmd.RegisterFlagCompletionFunc("name", cobra.NoFileCompletions)
	backupsNewCmd.RegisterFlagCompletionFunc("password", cobra.NoFileCompletions)
	backupsNewCmd.RegisterFlagCompletionFunc("uncompressed", boolCompletions)
	backupsNewCmd.RegisterFlagCompletionFunc("addons", cobra.NoFileCompletions)
	backupsNewCmd.RegisterFlagCompletionFunc("folders", cobra.NoFileCompletions)
	backupsNewCmd.RegisterFlagCompletionFunc("location", backupsLocationsCompletions)
	backupsNewCmd.RegisterFlagCompletionFunc("homeassistant-exclude-database", boolCompletions)

	backupsCmd.AddCommand(backupsNewCmd)
}
