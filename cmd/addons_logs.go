package cmd

import (
	"errors"
	"fmt"

	helper "github.com/home-assistant/hassio-cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var addonsLogsCmd = &cobra.Command{
	Use:     "logs [slug]",
	Aliases: []string{"in"},
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("addons logs")

		section := "addons"
		command := "{slug}/logs"
		base := viper.GetString("endpoint")

		url, err := helper.URLHelper(base, section, command)

		if err != nil {
			fmt.Println(err)
			ExitWithError = true
			return
		}

		request := helper.GetRequest()

		slug := args[0]

		request.SetPathParams(map[string]string{
			"slug": slug,
		})

		resp, err := request.SetHeader("Accept", "text/plain").Get(url)

		// returns 200 OK or 400, everything else is wrong
		if err == nil && resp.StatusCode() != 200 && resp.StatusCode() != 400 {
			err = errors.New("Unexpected server response")
			log.Error(err)
		}

		if err != nil {
			fmt.Println(err)
			ExitWithError = true
		} else {
			fmt.Println(string(resp.Body()))
		}

		return
	},
}

func init() {
	addonsCmd.AddCommand(addonsLogsCmd)
}
