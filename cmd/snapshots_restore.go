package cmd

import (
	"errors"
	"fmt"
	"time"

	resty "github.com/go-resty/resty/v2"
	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var snapshotsRestoreCmd = &cobra.Command{
	Use:   "restore [slug]",
	Short: "Restores a Home Assistant snapshot backup",
	Long: `
When something goes wrong, this command allows you to restore a previously
take Home Assistant snapshot backup on your system.`,
	Example: `
  ha snapshots restore c1a07617
  ha snapshots restore c1a07617 --addons core_ssh --addon core_mosquitto
  ha snapshots restore c1a07617 --folders config`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("snapshots restore")

		section := "snapshots/{slug}"
		command := "restore/full"
		base := viper.GetString("endpoint")

		request := helper.GetJSONRequestTimeout(3 * time.Hour)

		options := make(map[string]interface{})

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

		url, err := helper.URLHelper(base, section, command)
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

		// returns 200 OK or 400, everything else is wrong
		if err == nil {
			if resp.StatusCode() != 200 && resp.StatusCode() != 400 {
				err = errors.New("Unexpected server response")
				log.Error(err)
			} else if !resty.IsJSONType(resp.Header().Get("Content-Type")) {
				err = errors.New("API did not return a JSON response")
				log.Error(err)
			}
		}

		if err != nil {
			fmt.Println(err)
			ExitWithError = true
		} else {
			ExitWithError = !helper.ShowJSONResponse(resp)
		}
	},
}

func init() {
	snapshotsRestoreCmd.Flags().StringP("password", "", "", "Password")
	snapshotsRestoreCmd.Flags().BoolP("homeassistant", "", true, "Restore homeassistant (default true), triggers a partial backup when se to false")
	snapshotsRestoreCmd.Flags().StringArrayP("addons", "a", []string{}, "addons to restore, triggers a partial backup")
	snapshotsRestoreCmd.Flags().StringArrayP("folders", "f", []string{}, "folders to restore, triggers a partial backup")

	snapshotsCmd.AddCommand(snapshotsRestoreCmd)
}
