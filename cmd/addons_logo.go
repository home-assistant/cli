package cmd

import (
	"errors"
	"fmt"

	helper "github.com/home-assistant/hassio-cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var addonsLogoCmd = &cobra.Command{
	Use:     "logo [slug]",
	Aliases: []string{"in"},
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("addons logo")

		section := "addons"
		command := "{slug}/logo"
		base := viper.GetString("endpoint")

		url, err := helper.URLHelper(base, section, command)

		if err != nil {
			fmt.Println(err)
			ExitWithError = true
			return
		}

		request := helper.GetRequest()

		slug := args[0]
		output := args[1]

		request.SetOutput(output)

		request.SetPathParams(map[string]string{
			"slug": slug,
		})

		resp, err := request.Get(url)

		// returns 200 OK or 400, everything else is wrong
		if err == nil && resp.StatusCode() != 200 && resp.StatusCode() != 400 {
			err = errors.New("Unexpected server response")
			log.Error(err)
		}

		if err != nil {
			fmt.Println(err)
			ExitWithError = true
		}

		return
	},
}

func init() {
	addonsCmd.AddCommand(addonsLogoCmd)
}
