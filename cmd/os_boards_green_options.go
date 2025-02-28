package cmd

import (
	"fmt"
	"strings"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var osBoardsGreenOptionsCmd = &cobra.Command{
	Use:     "options",
	Aliases: []string{"option", "opt", "opts", "op"},
	Short:   "Change settings of the current Green board",
	Long: `
This command allows you to change settings of the Green board that Home
Assistant is running on.`,
	Example: `
  ha os boards green options --activity-led=false`,
	ValidArgsFunction: cobra.NoFileCompletions,
	Args:              cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("os boards green options")

		section := "os"
		command := "boards/green"

		options := make(map[string]any)

		for _, value := range []string{
			"activity-led",
			"power-led",
			"system-health-led",
		} {
			data, err := cmd.Flags().GetBool(value)
			if err == nil && cmd.Flags().Changed(value) {
				options[strings.ReplaceAll(value, "-", "_")] = data
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
	osBoardsGreenOptionsCmd.Flags().Bool("activity-led", true, "Enable/disable the green activity LED")
	osBoardsGreenOptionsCmd.Flags().Bool("power-led", true, "Enable/disable the white power LED")
	osBoardsGreenOptionsCmd.Flags().Bool("system-health-led", true, "Enable/disable the yellow system health LED")
	osBoardsGreenOptionsCmd.Flags().Lookup("activity-led").NoOptDefVal = "true"
	osBoardsGreenOptionsCmd.Flags().Lookup("power-led").NoOptDefVal = "true"
	osBoardsGreenOptionsCmd.Flags().Lookup("system-health-led").NoOptDefVal = "true"

	osBoardsGreenCmd.AddCommand(osBoardsGreenOptionsCmd)
}
