package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var hostReloadCmd = &cobra.Command{
	Use:     "reload",
	Aliases: []string{"update", "refresh", "re"},
	Short:   "Reload information from the host machine",
	Long: `
This commands reload the information Home Assistant has on the hostmachine.
If some setting are changed outside of Home Assistant, this commands updates
the internals of Home Assistant.`,
	Example: `
  ha host reload`,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("host reload")

		section := "host"
		command := "reload"
		base := viper.GetString("endpoint")

		resp, err := helper.GenericJSONPost(base, section, command, nil)
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
	hostCmd.AddCommand(hostReloadCmd)
}
