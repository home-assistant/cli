package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var resolutionHealthCheckCmd = &cobra.Command{
	Use:     "healthcheck",
	Aliases: []string{"check", "run"},
	Short:   "Execute system healthcheck and fixups",
	Long: `
This command execute a full system check and auto fixups. It check all issues again to see if they
are still around and try to fix it again.`,
	Example: `
  ha resolution healthcheck`,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("resolution")

		section := "resolution"
		command := "healthcheck"
		base := viper.GetString("endpoint")

		ProgressSpinner.Start()
		resp, err := helper.GenericJSONPost(base, section, command, nil)
		ProgressSpinner.Stop()

		if err != nil {
			fmt.Println(err)
			ExitWithError = true
		} else {
			ExitWithError = !helper.ShowJSONResponse(resp)
		}
	},
}

func init() {
	// add cmd to root command
	resolutionCmd.AddCommand(resolutionHealthCheckCmd)
}
