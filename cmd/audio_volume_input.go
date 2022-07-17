package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var audioVolumeInputCmd = &cobra.Command{
	Use:     "input",
	Aliases: []string{"in"},
	Short:   "Set volume of a Home Assistant Audio input channel",
	Long: `
This command allows you to set the volume of a Home Assistant Audio
input channel or application on your Home Assistant system.`,
	Example: `
	ha audio volume input --index 1 --mute
	ha audio volume input --index 1 --unmute
	ha audio volume input --index 1 --volume 75
	ha audio volume input --index 1 --mute --application
	ha audio volume input --index 1 --unmute --application
	ha audio volume input --index 2 --volume 50 --application
`,
	ValidArgsFunction: cobra.NoFileCompletions,
	Args:              cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("audio volume input")

		section := "audio"
		command := "volume/input"

		options := make(map[string]interface{})

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
			command = "mute/input/application"
		} else if mute || unmute {
			command = "mute/input"
		} else if application {
			command = "volume/input/application"
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
	audioVolumeInputCmd.Flags().Bool("application", false, "The index provided is an application")
	audioVolumeInputCmd.Flags().Int("index", 0, "Channel index number")
	audioVolumeInputCmd.Flags().Int("volume", 0, "Volume level in a percentage")
	audioVolumeInputCmd.Flags().Bool("mute", false, "Mute the channel")
	audioVolumeInputCmd.Flags().Bool("unmute", false, "Unmute the channel")
	audioVolumeInputCmd.MarkFlagRequired("index")
	audioVolumeInputCmd.RegisterFlagCompletionFunc("application", boolCompletions)
	audioVolumeInputCmd.RegisterFlagCompletionFunc("index", cobra.NoFileCompletions)
	audioVolumeInputCmd.RegisterFlagCompletionFunc("volume", cobra.NoFileCompletions)
	audioVolumeInputCmd.RegisterFlagCompletionFunc("mute", boolCompletions)
	audioVolumeInputCmd.RegisterFlagCompletionFunc("unmute", boolCompletions)
	audioVolumeCmd.AddCommand(audioVolumeInputCmd)
}
