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

		options := make(map[string]interface{})

		for _, value := range []string{
			"hostname",
			"channel",
			"timezone",
		} {
			val, err := cmd.Flags().GetString(value)
			if val != "" && err == nil && cmd.Flags().Changed(value) {
				options[value] = val
			}
		}

		waitboot, err := cmd.Flags().GetInt("wait-boot")
		if cmd.Flags().Changed("wait-boot") {
			options["wait_boot"] = waitboot
		}

		repos, err := cmd.Flags().GetStringArray("repositories")
		log.WithField("repositories", repos).Debug("repos")

		if len(repos) >= 0 && err == nil {
			options["addons_repositories"] = repos
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
	supervisorOptionsCmd.Flags().StringP("hostname", "", "", "Hostname to set")
	supervisorOptionsCmd.Flags().StringP("channel", "c", "", "Channel to track (stable|beta|dev)")
	supervisorOptionsCmd.Flags().StringP("timezone", "t", "", "Timezone")
	supervisorOptionsCmd.Flags().IntP("wait-boot", "w", 0, "Seconds to wait after boot")
	supervisorOptionsCmd.Flags().StringArrayP("repositories", "r", []string{}, "repositories to track, can be supplied multiple times")
	supervisorCmd.AddCommand(supervisorOptionsCmd)
}
