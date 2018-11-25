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

// reloadCmd represents the info command
var snapshotsRemoveCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"re"},
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("snapshots remove")

		section := "snapshots"
		command := "{slug}/remove"
		base := viper.GetString("endpoint")

		url, err := helper.URLHelper(base, section, command)
		fmt.Println(url)
		if err != nil {
			fmt.Println(err)
			return
		}

		request := helper.GetJSONRequest()

		slug, _ := cmd.Flags().GetString("slug")

		request.SetPathParams(map[string]string{
			"slug": slug,
		})

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
	snapshotsRemoveCmd.Flags().StringP("slug", "", "", "Slug of the snapshot")
	snapshotsRemoveCmd.MarkFlagRequired("slug")

	snapshotsCmd.AddCommand(snapshotsRemoveCmd)
}
