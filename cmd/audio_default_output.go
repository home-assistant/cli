package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("audio default output")

		section := "audio"
		command := "default/output"
		base := viper.GetString("endpoint")

		options := make(map[string]interface{})

		name, err := cmd.Flags().GetString("name")
		if name != "" && err == nil && cmd.Flags().Changed("name") {
			options["name"] = name
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
	audioDefaultOutputCmd.Flags().String("name", "", "The name of the audio device")
	audioDefaultOutputCmd.MarkFlagRequired("name")
	audioDefaultCmd.AddCommand(audioDefaultOutputCmd)
}
