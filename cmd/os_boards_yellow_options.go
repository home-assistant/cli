package cmd

import (
	"log/slog"
	"strings"

	helper "github.com/home-assistant/cli/client"
	"github.com/spf13/cobra"
)

var osBoardsYellowOptionsCmd = &cobra.Command{
	Use:     "options",
	Aliases: []string{"option", "opt", "opts", "op"},
	Short:   "Change settings of the current Yellow board",
	Long: `
This command allows you to change settings of the Yellow board that Home
Assistant is running on. A host reboot is required for changes to take effect.`,
	Example: `
  ha os boards yellow options --heartbeat-led=false`,
	ValidArgsFunction: cobra.NoFileCompletions,
	Args:              cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		slog.Debug("os boards yellow options", "args", args)

		section := "os"
		command := "boards/yellow"

		options := make(map[string]any)

		for _, value := range []string{
			"disk-led",
			"heartbeat-led",
			"power-led",
		} {
			data, err := cmd.Flags().GetBool(value)
			if err == nil && cmd.Flags().Changed(value) {
				options[strings.ReplaceAll(value, "-", "_")] = data
			}
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
	osBoardsYellowOptionsCmd.Flags().Bool("disk-led", true, "Enable/disable the disk LED")
	osBoardsYellowOptionsCmd.Flags().Bool("heartbeat-led", true, "Enable/disable the heartbeat LED")
	osBoardsYellowOptionsCmd.Flags().Bool("power-led", true, "Enable/disable the power LED")
	osBoardsYellowOptionsCmd.Flags().Lookup("disk-led").NoOptDefVal = "true"
	osBoardsYellowOptionsCmd.Flags().Lookup("heartbeat-led").NoOptDefVal = "true"
	osBoardsYellowOptionsCmd.Flags().Lookup("power-led").NoOptDefVal = "true"

	osBoardsYellowCmd.AddCommand(osBoardsYellowOptionsCmd)
}
