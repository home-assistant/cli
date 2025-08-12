package cmd

import (
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var audioVolumeCmd = &cobra.Command{
	Use:     "volume",
	Aliases: []string{"vol", "sound", "snd"},
	Short:   "Audio device volume control.",
	Long: `
Control the volume of your audio devices.
`,
	Example: `
	ha audio volume input --index 1 --mute
	ha audio volume input --index 2 --volume 75
	ha audio volume output --index 3 --unmute
	ha audio volume output --index 4 --volume 75 --application
`,
}

func init() {
	log.Debug("Init audio volume")

	audioCmd.AddCommand(audioVolumeCmd)
}

func volumePercentCompletions(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	vals := make([]string, 0, 101)
	for i := 0; i <= 100; i++ {
		vals = append(vals, strconv.Itoa(i))
	}
	return vals, cobra.ShellCompDirectiveNoFileComp
}
