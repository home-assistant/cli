package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var jobsOptionsCmd = &cobra.Command{
	Use:     "options",
	Aliases: []string{"option", "opt", "opts", "op"},
	Short:   "Allow to set options for the Job Manager backend",
	Long: `
This command allows you to set configuration options for the internally
Home Assistant Job Manager.
`,
	Example: `
  ha jobs options --ignore-conditions healthy
`,
	ValidArgsFunction: cobra.NoFileCompletions,
	Args:              cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("jobs options")

		section := "jobs"
		command := "options"

		options := make(map[string]interface{})

		conditions, err := cmd.Flags().GetStringArray("ignore-conditions")
		if len(conditions) >= 1 && err == nil {
			options["ignore_conditions"] = conditions
		}

		resp, err := helper.GenericJSONPost(section, command, options)
		if err != nil {
			fmt.Println(err)
			ExitWithError = true
		} else {
			ExitWithError = !helper.ShowJSONResponse(resp)
		}
	},
}

func init() {
	jobsOptionsCmd.Flags().StringArrayP("ignore-conditions", "i", []string{}, "Conditions to ignore on Job Manager. Use multiple times for ignored conditions.")
	jobsOptionsCmd.RegisterFlagCompletionFunc("ignore-conditions", cobra.NoFileCompletions)
	jobsCmd.AddCommand(jobsOptionsCmd)
}
