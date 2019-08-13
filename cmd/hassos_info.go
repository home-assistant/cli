package cmd

import (
	"fmt"

	helper "github.com/home-assistant/hassio-cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var hassosInfoCmd = &cobra.Command{
	Use:     "info",
	Aliases: []string{"in", "inf"},
	Short:   "Provides information about the running HassOS",
	Long: `
This command provides general information about the running HassOS system.
`,
	Example: `
  hassio hassos info
`,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("hassos info")

		section := "hassos"
		command := "info"
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
	hassosCmd.AddCommand(hassosInfoCmd)
}
