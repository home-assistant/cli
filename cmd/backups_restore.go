package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
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
  ha backups restore c1a07617 --addons core_ssh --addons core_mosquitto
  ha backups restore c1a07617 --folders homeassistant`,
	ValidArgsFunction: backupsCompletions,
	Args:              cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("backups restore")

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

		addons, err := cmd.Flags().GetStringArray("addons")
		log.WithField("addons", addons).Debug("addons")

		if len(addons) > 0 && err == nil {
			options["addons"] = addons
			command = "restore/partial"
		}

		folders, err := cmd.Flags().GetStringArray("folders")
		log.WithField("folders", folders).Debug("folders")

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
			fmt.Println(err)
			ExitWithError = true
			return
		}

		if len(options) > 0 {
			log.WithField("options", options).Debug("Request body")
			request.SetBody(options)
		}

		ProgressSpinner.Start()
		resp, err := request.Post(url)
		ProgressSpinner.Stop()

		resp, err = helper.GenericJSONErrorHandling(resp, err)

		if err != nil {
			fmt.Println(err)
			ExitWithError = true
		} else {
			ExitWithError = !helper.ShowJSONResponse(resp)
		}
	},
}

func init() {
	backupsRestoreCmd.Flags().StringP("password", "", "", "Password")
	backupsRestoreCmd.Flags().BoolP("homeassistant", "", true, "Restore homeassistant (default true), triggers a partial backup when set to false")
	backupsRestoreCmd.Flags().StringArrayP("addons", "a", []string{}, "addons to restore, triggers a partial backup")
	backupsRestoreCmd.Flags().StringArrayP("folders", "f", []string{}, "folders to restore, triggers a partial backup")
	backupsRestoreCmd.Flags().StringP("location", "l", "", "where to put backup file (backup mount or local)")

	backupsRestoreCmd.Flags().Lookup("location").NoOptDefVal = ".local"

	backupsRestoreCmd.RegisterFlagCompletionFunc("password", cobra.NoFileCompletions)
	backupsRestoreCmd.RegisterFlagCompletionFunc("homeassistant", boolCompletions)
	backupsRestoreCmd.RegisterFlagCompletionFunc("addons", cobra.NoFileCompletions)
	backupsRestoreCmd.RegisterFlagCompletionFunc("folders", backupsFoldersCompletions)
	backupsRestoreCmd.RegisterFlagCompletionFunc("location", backupsLocationsCompletions)

	backupsCmd.AddCommand(backupsRestoreCmd)
}
