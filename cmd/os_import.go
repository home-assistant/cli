package cmd

import (
	"log/slog"

	helper "github.com/home-assistant/cli/client"
	"github.com/spf13/cobra"
)

var osImportCmd = &cobra.Command{
	Use:     "import",
	Aliases: []string{"im", "sync", "load"},
	Short:   "Import configurations from a USB stick",
	Long: `
This commands triggers an import action from a connected USB stick with
configuration to load for the Home Assistant Operating System.
`,
	Example: `
  ha os import
`,
	ValidArgsFunction: cobra.NoFileCompletions,
	Args:              cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		slog.Debug("os import", "args", args)

		section := "os"
		command := "config/sync"

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
	osCmd.AddCommand(osImportCmd)
}
