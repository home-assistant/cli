package cmd

import (
	"fmt"

	helper "github.com/home-assistant/hassio-cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// optionsCmd represents the options command
var supervisorOptionsCmd = &cobra.Command{
	Use:     "options",
	Aliases: []string{"op"},
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("supervisor options")

		section := "supervisor"
		command := "options"
		base := viper.GetString("endpoint")

		url, err := helper.URLHelper(base, section, command)
		if err != nil {
			// TODO: error handler
			fmt.Printf("Error: %v", err)
			return
		}

		options := make(map[string]interface{})

		for _, value := range []string{
			"hostname",
			"channel",
			"timezone",
		} {
			val, err := cmd.Flags().GetString(value)
			if val != "" && err == nil {
				options[value] = val
			}
		}

		waitboot, err := cmd.Flags().GetInt("wait-boot")
		if waitboot != 0 {
			options["wait_boot"] = waitboot
		}

		repos, err := cmd.Flags().GetStringArray("repositories")
		log.WithField("repositories", repos).Debug("repos")

		if len(repos) >= 0 && err == nil {
			options["addons_repositories"] = repos
		}

		request := helper.GetJSONRequest()
		if len(options) > 0 {
			log.WithField("options", options).Debug("Sending options")
			request.SetBody(options)
		}
		resp, err := request.Post(url)
		log.WithField("Request", resp.Request.RawRequest).Debug("Request")

		// returns 200 OK or 400
		if resp.StatusCode() != 200 && resp.StatusCode() != 400 {
			fmt.Println("Unexpected server response")
			fmt.Println(resp.String())
		} else {
			helper.ShowJSONResponse(resp)
		}

		return
	},
}

func init() {
	supervisorOptionsCmd.Flags().StringP("hostname", "", "", "Hostname to set")
	supervisorOptionsCmd.Flags().StringP("channel", "c", "", "Channel to track (stable|beta|dev)")
	supervisorOptionsCmd.Flags().StringP("timezone", "t", "", "Timezone")
	supervisorOptionsCmd.Flags().IntP("wait-boot", "w", 0, "Seconds to wait after boot")
	supervisorOptionsCmd.Flags().StringArrayP("repositories", "r", []string{}, "repositories to track, can be supplied multiple times")
	supervisorCmd.AddCommand(supervisorOptionsCmd)
}
