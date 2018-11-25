package cmd

import (
	"fmt"

	helper "github.com/home-assistant/hassio-cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var snapshotsNewCmd = &cobra.Command{
	Use:     "new",
	Aliases: []string{"re"},
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("snapshots new")

		section := "snapshots"
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

		if len(addons) >= 0 && err == nil {
			options["addons"] = addons
			command = "new/partial"
		}

		folders, err := cmd.Flags().GetStringArray("folders")
		log.WithField("folders", folders).Debug("folders")

		if len(folders) >= 0 && err == nil {
			options["folders"] = folders
			command = "new/partial"
		}

		resp, err := helper.GenericJSONPost(base, section, command, options)
		if err != nil {
			fmt.Println(err)
		} else {
			helper.ShowJSONResponse(resp)
		}

		return
	},
}

func init() {
	snapshotsNewCmd.Flags().StringP("name", "", "", "Name of the snapshot")
	snapshotsNewCmd.Flags().StringP("password", "", "", "Password")
	snapshotsNewCmd.Flags().StringArrayP("addons", "a", []string{}, "addons to backup, triggers a partial backup")
	snapshotsNewCmd.Flags().StringArrayP("folders", "f", []string{}, "folders to backup, triggers a partial backup")

	snapshotsCmd.AddCommand(snapshotsNewCmd)
}
