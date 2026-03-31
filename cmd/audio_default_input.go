package cmd

import (
	"log/slog"

	helper "github.com/home-assistant/cli/client"
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
	ValidArgsFunction: cobra.NoFileCompletions,
	Args:              cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		slog.Debug("audio default input", "args", args)

		section := "audio"
		command := "default/input"

		options := make(map[string]any)

		name, err := cmd.Flags().GetString("name")
		if name != "" && err == nil && cmd.Flags().Changed("name") {
			options["name"] = name
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
	audioDefaultInputCmd.Flags().String("name", "", "The name of the audio device")
	audioDefaultInputCmd.MarkFlagRequired("name")
	audioDefaultInputCmd.RegisterFlagCompletionFunc("name", cobra.NoFileCompletions)
	audioDefaultCmd.AddCommand(audioDefaultInputCmd)
}
