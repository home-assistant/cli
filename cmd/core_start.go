package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var coreStartCmd = &cobra.Command{
	Use:     "start",
	Aliases: []string{"run", "st"},
	Short:   "Manually start Home Assistant Core",
	Long: `
This command allows you to manually start the Home Assistant Core instance on
your system. This, of course, only applies when it has been stopped.`,
	Example: `
  ha core start`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("core start")

		section := "core"
		command := "start"

		ProgressSpinner.Start()
		resp, err := helper.GenericJSONPostTimeout(section, command, nil, helper.ContainerOperationTimeout)
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
	coreCmd.AddCommand(coreStartCmd)
}
