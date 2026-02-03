package cmd

import (
	"log/slog"

	helper "github.com/home-assistant/cli/client"
	"github.com/spf13/cobra"
)

var hardwareAudioCmd = &cobra.Command{
	Use:     "audio",
	Aliases: []string{"sounds", "snd", "au"},
	Short:   "Provides information about audio devices on your system",
	Long: `
The command provides information about audio devices available on your system.`,
	Example:           `ha hardware audio`,
	ValidArgsFunction: cobra.NoFileCompletions,
	Args:              cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		slog.Debug("hardware audio", "args", args)

		section := "hardware"
		command := "audio"

		resp, err := helper.GenericJSONGet(section, command)
		if err != nil {
			helper.PrintError(err)
			ExitWithError = true
		} else {
			ExitWithError = !helper.ShowJSONResponse(resp)
		}
	},
}

func init() {
	hardwareCmd.AddCommand(hardwareAudioCmd)
}
