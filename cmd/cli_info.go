package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cliInfoCmd = &cobra.Command{
	Use:     "info",
	Aliases: []string{"in", "inf"},
	Short:   "Shows information about the internal Home Assistant CLI backend",
	Long: `
Shows information about the internally running Home Assistant CLI backend
`,
	Example: `
  ha cli info
`,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("cli info")

		section := "cli"
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
	cliCmd.AddCommand(cliInfoCmd)
}
