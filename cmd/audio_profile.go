package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var audioProfileCmd = &cobra.Command{
	Use:     "profile",
	Aliases: []string{"pro"},
	Short:   "Set the Home Assistant Audio profile for a card",
	Long: `
This command allows you to set the audio profile on a audio card.`,
	Example: `
	ha audio profile --card "..." --name "..."
`,

	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("audio profile")

		section := "audio"
		command := "profile"
		base := viper.GetString("endpoint")

		options := make(map[string]interface{})

		name, err := cmd.Flags().GetString("name")
		if name != "" && err == nil && cmd.Flags().Changed("name") {
			options["name"] = name
		}

		profile, err := cmd.Flags().GetString("profile")
		if name != "" && err == nil && cmd.Flags().Changed("profile") {
			options["profile"] = profile
		}

		resp, err := helper.GenericJSONPost(base, section, command, options)
		if err != nil {
			fmt.Println(err)
			ExitWithError = true
		} else {
			ExitWithError = !helper.ShowJSONResponse(resp)
		}

		return
	},
}

func init() {
	audioProfileCmd.Flags().String("name", "", "Name of the sound card")
	audioProfileCmd.Flags().String("profile", "", "Name of the profile")
	audioProfileCmd.MarkFlagRequired("name")
	audioProfileCmd.MarkFlagRequired("profile")
	audioCmd.AddCommand(audioProfileCmd)
}
