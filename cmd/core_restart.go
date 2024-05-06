package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var coreRestartCmd = &cobra.Command{
	Use:     "restart",
	Aliases: []string{"reboot"},
	Short:   "Restarts the Home Assistant Core",
	Long: `
Restart the Home Assistant Core instance running on your system`,
	Example: `
  ha core restart`,
	ValidArgsFunction: cobra.NoFileCompletions,
	Args:              cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("core restart")

		section := "core"
		command := "restart"

		options := make(map[string]interface{})

		safeMode, err := cmd.Flags().GetBool("safe-mode")
		if err == nil && safeMode {
			options["safe_mode"] = safeMode
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
	coreRestartCmd.Flags().BoolP("safe-mode", "s", false, "Restart Home Assistant in safe mode")
	coreRestartCmd.Flags().Lookup("safe-mode").NoOptDefVal = "true"
	coreRestartCmd.RegisterFlagCompletionFunc("safe-mode", boolCompletions)

	coreCmd.AddCommand(coreRestartCmd)
}
