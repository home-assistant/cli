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
Assistant is running on. A host reboot is required for changes to take effect.`,
	Example: `
  ha os boards green options --activity-led=false`,
	ValidArgsFunction: cobra.NoFileCompletions,
	Args:              cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("os boards green options")

		section := "os"
		command := "boards/green"

		options := make(map[string]interface{})

		for _, value := range []string{
			"activity-led",
			"power-led",
			"user-led",
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
	osBoardsGreenOptionsCmd.Flags().Bool("activity-led", true, "Enable/disable the activity LED")
	osBoardsGreenOptionsCmd.Flags().Bool("power-led", true, "Enable/disable the power LED")
	osBoardsGreenOptionsCmd.Flags().Bool("user-led", true, "Enable/disable the user LED")
	osBoardsGreenOptionsCmd.Flags().Lookup("activity-led").NoOptDefVal = "true"
	osBoardsGreenOptionsCmd.Flags().Lookup("power-led").NoOptDefVal = "true"
	osBoardsGreenOptionsCmd.Flags().Lookup("user-led").NoOptDefVal = "true"

	osBoardsGreenCmd.AddCommand(osBoardsGreenOptionsCmd)
}
