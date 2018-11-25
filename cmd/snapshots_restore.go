package cmd

import (
	"errors"
	"fmt"
	"net/http"

	helper "github.com/home-assistant/hassio-cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	resty "gopkg.in/resty.v1"
)

var snapshotsRestoreCmd = &cobra.Command{
	Use: "restore",
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("snapshots restore")

		section := "snapshots/{slug}"
		command := "restore/full"
		base := viper.GetString("endpoint")

		request := helper.GetJSONRequest()

		options := make(map[string]interface{})

		slug, _ := cmd.Flags().GetString("slug")

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
		fmt.Println(url)
		if err != nil {
			fmt.Println(err)
			return
		}

		if len(options) > 0 {
			log.WithField("options", options).Debug("Request body")
			request.SetBody(options)
		}

		resp, err := request.Post(url)

		// returns 200 OK or 400, everything else is wrong
		if err == nil {
			if resp.StatusCode() != 200 && resp.StatusCode() != 400 {
				err = errors.New("Unexpected server response")
				log.Error(err)
			} else if !resty.IsJSONType(resp.Header().Get(http.CanonicalHeaderKey("Content-Type"))) {
				err = errors.New("api did not return a json response")
				log.Error(err)
			}
		}

		if err != nil {
			fmt.Println(err)
		} else {
			helper.ShowJSONResponse(resp)
		}

		return
	},
}

func init() {
	snapshotsRestoreCmd.Flags().StringP("slug", "", "", "Slug of the snapshot")
	snapshotsRestoreCmd.MarkFlagRequired("slug")

	snapshotsRestoreCmd.Flags().StringP("password", "", "", "Password")
	snapshotsRestoreCmd.Flags().BoolP("homeassistant", "", true, "Restore homeassistant (default true), triggers a partial backup when se to false")
	snapshotsRestoreCmd.Flags().StringArrayP("addons", "a", []string{}, "addons to restore, triggers a partial backup")
	snapshotsRestoreCmd.Flags().StringArrayP("folders", "f", []string{}, "folders to restore, triggers a partial backup")

	snapshotsCmd.AddCommand(snapshotsRestoreCmd)
}
