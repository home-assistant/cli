package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
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

		options := make(map[string]interface{})

		name, err := cmd.Flags().GetString("name")
		if name != "" && err == nil && cmd.Flags().Changed("name") {
			options["name"] = name
		}

		card, err := cmd.Flags().GetString("card")
		if name != "" && err == nil && cmd.Flags().Changed("card") {
			options["card"] = card
		}

		resp, err := helper.GenericJSONPost(section, command, options)
		if err != nil {
			fmt.Println(err)
			ExitWithError = true
		} else {
			ExitWithError = !helper.ShowJSONResponse(resp)
		}
	},
}

func init() {
	audioProfileCmd.Flags().String("name", "", "Name of the profile")
	audioProfileCmd.Flags().String("card", "", "The card to set the profile for")
	audioProfileCmd.MarkFlagRequired("name")
	audioProfileCmd.MarkFlagRequired("card")
	audioCmd.AddCommand(audioProfileCmd)
}
