package cmd

import (
	"fmt"
	"strings"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
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
		log.WithField("args", args).Debug("os boards yellow options")

		section := "os"
		command := "boards/yellow"

		options := make(map[string]interface{})

		for _, value := range []string{
			"disk-led",
			"heartbeat-led",
			"power-led",
		} {
			data, err := cmd.Flags().GetBool(value)
			if err == nil && cmd.Flags().Changed(value) {
				options[strings.Replace(value, "-", "_", -1)] = data
			}
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
	osBoardsYellowOptionsCmd.Flags().Bool("disk-led", true, "Enable/disable the disk LED")
	osBoardsYellowOptionsCmd.Flags().Bool("heartbeat-led", true, "Enable/disable the heartbeat LED")
	osBoardsYellowOptionsCmd.Flags().Bool("power-led", true, "Enable/disable the power LED")
	osBoardsYellowOptionsCmd.Flags().Lookup("disk-led").NoOptDefVal = "true"
	osBoardsYellowOptionsCmd.Flags().Lookup("heartbeat-led").NoOptDefVal = "true"
	osBoardsYellowOptionsCmd.Flags().Lookup("power-led").NoOptDefVal = "true"

	osBoardsYellowCmd.AddCommand(osBoardsYellowOptionsCmd)
}
