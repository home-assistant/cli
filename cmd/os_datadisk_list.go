package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("os datadisk list")

		section := "os"
		command := "datadisk/list"
		base := viper.GetString("endpoint")

		resp, err := helper.GenericJSONGet(base, section, command)
		if err != nil {
			fmt.Println(err)
			ExitWithError = true
		} else {
			ExitWithError = !helper.ShowJSONResponse(resp)
		}
	},
}

func init() {
	osDataDiskCmd.AddCommand(osDataDiskListCmd)
}
