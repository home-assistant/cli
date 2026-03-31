package cmd

import (
	"log/slog"

	helper "github.com/home-assistant/cli/client"
	"github.com/spf13/cobra"
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
	ValidArgsFunction: cobra.NoFileCompletions,
	Args:              cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		slog.Debug("supervisor repair", "args", args)

		section := "supervisor"
		command := "repair"

		ProgressSpinner.Start()
		resp, err := helper.GenericJSONPostTimeout(section, command, nil, helper.ContainerDownloadTimeout)
		ProgressSpinner.Stop()
		if err != nil {
			helper.PrintError(err)
			ExitWithError = true
		} else {
			ExitWithError = !helper.ShowJSONResponse(resp)
		}
	},
}

func init() {
	supervisorCmd.AddCommand(supervisorRepairCmd)
}
