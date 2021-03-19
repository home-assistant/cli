package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var networkReloadCmd = &cobra.Command{
	Use:     "reload",
	Aliases: []string{"re"},
	Short:   "Reload Network information the host",
	Long: `
Reload information about the host network and interfaces.
`,
	Example: `
  ha network reload
`,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("network reload")

		section := "network"
		command := "reload"
		base := viper.GetString("endpoint")

		resp, err := helper.GenericJSONPost(base, section, command, nil)
		if err != nil {
			fmt.Println(err)
			ExitWithError = true
		} else {
			ExitWithError = !helper.ShowJSONResponse(resp)
		}
	},
}

func init() {
	networkCmd.AddCommand(networkReloadCmd)
}
