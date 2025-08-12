package cmd

import (
	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var resolutionCheckOptionsCmd = &cobra.Command{
	Use:     "options",
	Aliases: []string{"option", "opt", "opts", "op"},
	Short:   "Options apply to checks managed by the Resolution center",
	Long: `
This command allows to apply options to an specific check managed by the system.`,
	Example: `
  ha resolution check options [slug]`,
	ValidArgsFunction: resolutionCheckCompletions,
	Args:              cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("check options")

		section := "resolution"
		command := "check/{check}/options"
		options := make(map[string]any)

		url, err := helper.URLHelper(section, command)
		if err != nil {
			helper.PrintError(err)
			ExitWithError = true
			return
		}

		request := helper.GetJSONRequest()

		check := args[0]

		request.SetPathParams(map[string]string{
			"check": check,
		})

		enabled, err := cmd.Flags().GetBool("enabled")
		if err == nil && cmd.Flags().Changed("enabled") {
			options["enabled"] = enabled
		}

		if len(options) > 0 {
			log.WithField("options", options).Debug("Request body")
			request.SetBody(options)
		}
		resp, err := request.Post(url)
		resp, err = helper.GenericJSONErrorHandling(resp, err)

		if err != nil {
			helper.PrintError(err)
			ExitWithError = true
		} else {
			ExitWithError = !helper.ShowJSONResponse(resp)
		}
	},
}

func init() {
	resolutionCheckOptionsCmd.Flags().BoolP("enabled", "", true, "Enable/Disable check on the system")
	resolutionCheckOptionsCmd.Flags().Lookup("enabled").NoOptDefVal = "true"
	resolutionCheckOptionsCmd.RegisterFlagCompletionFunc("enabled", boolCompletions)
	resolutionCheckCmd.AddCommand(resolutionCheckOptionsCmd)
}
