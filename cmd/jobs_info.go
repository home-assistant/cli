package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var jobsInfoCmd = &cobra.Command{
	Use:     "info",
	Aliases: []string{"in", "inf"},
	Short:   "Provides information about the Home Assistant Job Manager",
	Long: `
This command provides general information about the Home Assistant Job backend.
`,
	Example: `
  ha jobs info
`,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("jobs info")

		section := "jobs"
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
	jobsCmd.AddCommand(jobsInfoCmd)
}
