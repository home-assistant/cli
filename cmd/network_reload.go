package cmd

import (
	"log/slog"

	helper "github.com/home-assistant/cli/client"
	"github.com/spf13/cobra"
)

var networkReloadCmd = &cobra.Command{
	Use:     "reload",
	Aliases: []string{"re"},
	Short:   "Reload Network information the host",
	Long: `
Reload information about the host network and interfaces.
`,
	Example: `
  ha network reload
`,
	ValidArgsFunction: cobra.NoFileCompletions,
	Args:              cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		slog.Debug("network reload", "args", args)

		section := "network"
		command := "reload"

		resp, err := helper.GenericJSONPost(section, command, nil)
		if err != nil {
			helper.PrintError(err)
			ExitWithError = true
		} else {
			ExitWithError = !helper.ShowJSONResponse(resp)
		}
	},
}

func init() {
	networkCmd.AddCommand(networkReloadCmd)
}
