package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
`,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("backups new")

		section := "backups"
		command := "new/full"
		base := viper.GetString("endpoint")

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

		if len(addons) >= 0 && err == nil && cmd.Flags().Changed("addons") {
			options["addons"] = addons
			command = "new/partial"
		}

		folders, err := cmd.Flags().GetStringArray("folders")
		log.WithField("folders", folders).Debug("folders")

		if len(folders) >= 0 && err == nil && cmd.Flags().Changed("folders") {
			options["folders"] = folders
			command = "new/partial"
		}

		ProgressSpinner.Start()
		resp, err := helper.GenericJSONPostTimeout(base, section, command, options, helper.BackupTimeout)
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
	backupsNewCmd.Flags().StringArrayP("addons", "a", []string{}, "addons to backup, triggers a partial backup")
	backupsNewCmd.Flags().StringArrayP("folders", "f", []string{}, "folders to backup, triggers a partial backup")

	backupsCmd.AddCommand(backupsNewCmd)
}
