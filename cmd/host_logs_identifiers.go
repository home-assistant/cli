package cmd

import (
	"log/slog"

	helper "github.com/home-assistant/cli/client"
	"github.com/spf13/cobra"
)

var hostLogsIdentifiersCmd = &cobra.Command{
	Use:     "identifiers",
	Aliases: []string{"ids", "list-identifiers", "li"},
	Short:   "Show all syslog identifiers",
	Long: `
Show all values that can be used with the identifier arg to find logs.
`,
	Example: `
  ha host logs identifiers
`,
	ValidArgsFunction: cobra.NoFileCompletions,
	Args:              cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		slog.Debug("host logs identifiers", "args", args)

		section := "host"
		command := "logs/identifiers"

		resp, err := helper.GenericJSONGet(section, command)
		if err != nil {
			helper.PrintError(err)
			ExitWithError = true
		} else {
			ExitWithError = !helper.ShowJSONResponse(resp)
		}
	},
}

func init() {
	hostLogsCmd.AddCommand(hostLogsIdentifiersCmd)
}
