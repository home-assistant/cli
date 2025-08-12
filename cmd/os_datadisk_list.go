package cmd

import (
	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var osDataDiskListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"in", "inf", "info", "show"},
	Short:   "Provides information about the running Home Assistant Operating System",
	Long: `
This command provides general information about available Harddisk for using with Home Assistant Operating System.
`,
	Example: `
  ha os datadisk list
`,
	ValidArgsFunction: cobra.NoFileCompletions,
	Args:              cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("os datadisk list")

		section := "os"
		command := "datadisk/list"

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
	osDataDiskCmd.AddCommand(osDataDiskListCmd)
}
