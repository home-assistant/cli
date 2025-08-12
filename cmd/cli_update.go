package cmd

import (
	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var cliUpdateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"upgrade", "downgrade", "up", "down"},
	Short:   "Updates the internal Home Assistant CLI backend",
	Long: `
Using this command you can upgrade or downgrade the internal Home Assistant 
CLI backend, to the latest version or the version specified.
`,
	Example: `
  ha cli update --version 5
`,
	ValidArgsFunction: cobra.NoFileCompletions,
	Args:              cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("cli update")

		section := "cli"
		command := "update"

		var options map[string]any

		version, _ := cmd.Flags().GetString("version")
		if version != "" {
			options = map[string]any{"version": version}
		}

		ProgressSpinner.Start()
		resp, err := helper.GenericJSONPostTimeout(section, command, options, helper.ContainerDownloadTimeout)
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
	cliUpdateCmd.Flags().StringP("version", "", "", "Version to update to")
	cliUpdateCmd.RegisterFlagCompletionFunc("version", cobra.NoFileCompletions)
	cliCmd.AddCommand(cliUpdateCmd)
}
