package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var audioDefaultInputCmd = &cobra.Command{
	Use:     "input",
	Aliases: []string{"in"},
	Short:   "Set the default Home Assistant Audio input channel",
	Long: `
This command allows you to set the default input channel of the
Home Assistant Audio on your Home Assistant system.`,
	Example: `
	ha audio default input --name "..."
`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("audio default input")

		section := "audio"
		command := "default/input"

		options := make(map[string]interface{})

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
	audioDefaultInputCmd.Flags().String("name", "", "The name of the audio device")
	audioDefaultInputCmd.MarkFlagRequired("name")
	audioDefaultCmd.AddCommand(audioDefaultInputCmd)
}
