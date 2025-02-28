package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var coreStopCmd = &cobra.Command{
	Use:     "stop",
	Aliases: []string{},
	Short:   "Manually stop Home Assistant Core",
	Long: `
This command allows you to manually stop the Home Assistant Core instance on
your system.`,
	Example: `
  ha core stop`,
	ValidArgsFunction: cobra.NoFileCompletions,
	Args:              cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("core stop")

		section := "core"
		command := "stop"

		options := make(map[string]any)

		force, err := cmd.Flags().GetBool("force")
		if err == nil && force {
			options["force"] = force
		}

		ProgressSpinner.Start()
		resp, err := helper.GenericJSONPostTimeout(section, command, options, helper.ContainerOperationTimeout)
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
	coreStopCmd.Flags().BoolP("force", "f", false, "Force stop during an offline db migration")
	coreStopCmd.Flags().Lookup("force").NoOptDefVal = "true"
	coreStopCmd.RegisterFlagCompletionFunc("force", boolCompletions)
	coreCmd.AddCommand(coreStopCmd)
}
