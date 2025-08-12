package cmd

import (
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
	ValidArgsFunction: cobra.NoFileCompletions,
	Args:              cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("audio profile")

		section := "audio"
		command := "profile"

		options := make(map[string]any)

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
			helper.PrintError(err)
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
	audioProfileCmd.RegisterFlagCompletionFunc("name", cobra.NoFileCompletions)
	audioProfileCmd.RegisterFlagCompletionFunc("card", cobra.NoFileCompletions)
	audioCmd.AddCommand(audioProfileCmd)
}
