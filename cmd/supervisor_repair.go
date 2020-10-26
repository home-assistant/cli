package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var supervisorRepairCmd = &cobra.Command{
	Use:     "repair",
	Aliases: []string{"rep", "fix"},
	Short:   "Repair Docker issue automatically using the Supervisor",
	Long: `
There are cases where the Docker file system running on your Home Assistant
system, encounters issue or corruptions. Running this command,
the Home Assistant Supervisor will try to resolve these.
`,
	Example: `
  ha supervisor repair`,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("supervisor repair")

		section := "supervisor"
		command := "repair"
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
	supervisorCmd.AddCommand(supervisorRepairCmd)
}
