package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var coreInfoCmd = &cobra.Command{
	Use:     "info",
	Aliases: []string{"in", "inf"},
	Short:   "Provides information about Home Assistant Core",
	Long: `
This command provides information about the running Home Assistant Core instance
running on your Home Assistant system.`,
	Example: `
  ha core info`,
	ValidArgsFunction: cobra.NoFileCompletions,
	Args:              cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("core info")

		section := "core"
		command := "info"

		resp, err := helper.GenericJSONGet(section, command)
		if err != nil {
			fmt.Println(err)
			ExitWithError = true
		} else {
			ExitWithError = !helper.ShowJSONResponse(resp)
		}
	},
}

func init() {
	coreCmd.AddCommand(coreInfoCmd)
}
