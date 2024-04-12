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

		url, err := helper.URLHelper(section, command)
		if err != nil {
			fmt.Println(err)
			ExitWithError = true
			return
		}

		request := helper.GetRequestTimeout(helper.ContainerOperationTimeout)

		safeMode, err := cmd.Flags().GetBool("safe-mode")
		if err == nil && safeMode {
			options := make(map[string]interface{})
			options["safe_mode"] = safeMode
			request.SetBody(options)
		}

		ProgressSpinner.Start()
		resp, err := request.Post(url)
		ProgressSpinner.Stop()

		if err != nil {
			fmt.Println(err)
			ExitWithError = true
		} else if resp.StatusCode() != 200 && resp.StatusCode() != 400 {
			err = fmt.Errorf("Unexpected server response. Status code: %d", resp.StatusCode())
			log.Error(err)
			ExitWithError = true
		}
	},
}

func init() {
	coreRestartCmd.Flags().BoolP("safe-mode", "s", false, "Restart Home Assistant in safe mode")
	coreRestartCmd.Flags().Lookup("safe-mode").NoOptDefVal = "true"
	coreRestartCmd.RegisterFlagCompletionFunc("safe-mode", boolCompletions)

	coreCmd.AddCommand(coreRestartCmd)
}
