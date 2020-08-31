package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var networkInfoCmd = &cobra.Command{
	Use:     "info [interface]",
	Aliases: []string{"in", "inf"},
	Short:   "Shows information about the host network",
	Long: `
Shows information about the host network and interfaces or only from a specific interface.
`,
	Example: `
  ha network info
  ha network info eth0
`,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("network info")

		section := "network"
		command := "info"
		base := viper.GetString("endpoint")

		if len(args) > 0 {
			inet := args[0]
			command = inet + "/info"
		}

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
	networkCmd.AddCommand(networkInfoCmd)
}
