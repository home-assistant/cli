package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var audioDefaultOutputCmd = &cobra.Command{
	Use:     "output",
	Aliases: []string{"out"},
	Short:   "Set the default Home Assistant Audio output channel",
	Long: `
This command allows you to set the default output channel of the
Home Assistant Audio on your Home Assistant system.`,
	Example: `
	ha audio default output --name "..."
`,
	ValidArgsFunction: cobra.NoFileCompletions,
	Args:              cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("audio default output")

		section := "audio"
		command := "default/output"

		options := make(map[string]any)

		name, err := cmd.Flags().GetString("name")
		if name != "" && err == nil && cmd.Flags().Changed("name") {
			options["name"] = name
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
	audioDefaultOutputCmd.Flags().String("name", "", "The name of the audio device")
	audioDefaultOutputCmd.MarkFlagRequired("name")
	audioDefaultOutputCmd.RegisterFlagCompletionFunc("name", cobra.NoFileCompletions)
	audioDefaultCmd.AddCommand(audioDefaultOutputCmd)
}
