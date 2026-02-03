package cmd

import (
	"log/slog"

	helper "github.com/home-assistant/cli/client"
	"github.com/spf13/cobra"
)

var authListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all Home Assistant users.",
	Long: `
This command allows you to list all Home Assistant users on the system.
Please note, this command is limited due to security reasons, and will
only work on some locations. For example, the Operating System CLI.
`,
	Example: `
  ha authentication list
`,
	ValidArgsFunction: cobra.NoFileCompletions,
	Args:              cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		slog.Debug("auth list", "args", args)

		section := "auth"
		command := "list"

		resp, err := helper.GenericJSONGet(section, command)
		if err != nil {
			cmd.PrintErrln("this command is limited due to security reasons, and will only work on some locations. For example, the Operating System terminal.")
			helper.PrintError(err)
			ExitWithError = true
		} else {
			ExitWithError = !helper.ShowJSONResponse(resp)
		}
	},
}

func init() {
	authCmd.AddCommand(authListCmd)
}
