package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var resolutionInfoCmd = &cobra.Command{
	Use:     "info",
	Aliases: []string{"in", "inf"},
	Short:   "Show issues and suggestions",
	Long: `
This command provides general information about the issues, suggestion and the supported state of the system.`,
	Example: `
  ha resolution
  ha resolution suggestion apply [id]`,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("resolution")

		section := "resolution"
		command := "info"
		base := viper.GetString("endpoint")

		resp, err := helper.GenericJSONGet(base, section, command)
		if err != nil {
			fmt.Println(err)
			ExitWithError = true
		} else {
			ExitWithError = !helper.ShowJSONResponse(resp)
		}
		return
	},
}

func init() {
	// add cmd to root command
	resolutionCmd.AddCommand(resolutionInfoCmd)
}
