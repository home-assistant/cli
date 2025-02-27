package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var audioVolumeOuputCmd = &cobra.Command{
	Use:     "output",
	Aliases: []string{"out"},
	Short:   "Set volume of a Home Assistant Audio output channel",
	Long: `
This command allows you to set the volume of a Home Assistant Audio
output channel or application on your Home Assistant system.`,
	Example: `
	ha audio volume output --index 1 --mute
	ha audio volume output --index 1 --unmute
	ha audio volume output --index 1 --volume 75
	ha audio volume output --index 1 --mute --application
	ha audio volume output --index 1 --unmute --application
	ha audio volume output --index 2 --volume 50 --application
`,
	ValidArgsFunction: cobra.NoFileCompletions,
	Args:              cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("audio volume output")

		section := "audio"
		command := "volume/output"

		options := make(map[string]any)

		volume, err := cmd.Flags().GetInt("volume")
		if volume != 0 && err == nil && cmd.Flags().Changed("volume") {
			options["volume"] = float64(volume) / 100
		}

		index, err := cmd.Flags().GetInt("index")
		if err == nil && cmd.Flags().Changed("index") {
			options["index"] = index
		}

		mute, err := cmd.Flags().GetBool("mute")
		if err == nil && cmd.Flags().Changed("mute") {
			options["active"] = mute
		}

		unmute, err := cmd.Flags().GetBool("unmute")
		if err == nil && cmd.Flags().Changed("unmute") {
			options["active"] = !unmute
		}

		application, _ := cmd.Flags().GetBool("application")
		if (mute || unmute) && application {
			command = "mute/output/application"
		} else if mute || unmute {
			command = "mute/output"
		} else if application {
			command = "volume/output/application"
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
	audioVolumeOuputCmd.Flags().Bool("application", false, "The index provided is an application")
	audioVolumeOuputCmd.Flags().Int("index", 0, "Channel index number")
	audioVolumeOuputCmd.Flags().Int("volume", 0, "Volume level in a percentage")
	audioVolumeOuputCmd.Flags().Bool("mute", false, "Mute the channel")
	audioVolumeOuputCmd.Flags().Bool("unmute", false, "Unmute the channel")
	audioVolumeOuputCmd.MarkFlagRequired("index")
	audioVolumeOuputCmd.RegisterFlagCompletionFunc("application", boolCompletions)
	audioVolumeOuputCmd.RegisterFlagCompletionFunc("index", cobra.NoFileCompletions)
	audioVolumeOuputCmd.RegisterFlagCompletionFunc("volume", volumePercentCompletions)
	audioVolumeOuputCmd.RegisterFlagCompletionFunc("mute", boolCompletions)
	audioVolumeOuputCmd.RegisterFlagCompletionFunc("unmute", boolCompletions)
	audioVolumeCmd.AddCommand(audioVolumeOuputCmd)
}
