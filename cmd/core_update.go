package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var coreUpdateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"upgrade", "downgrade", "up", "down"},
	Short:   "Updates the Home Assistant Core",
	Long: `
Using this command you can upgrade or downgrade the Home Assistant Core instance
running on your system to the latest version or the version specified.`,
	Example: `
  ha core update
  ha core update --version 0.105.4`,
	ValidArgsFunction: cobra.NoFileCompletions,
	Args:              cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("core update")

		section := "core"
		command := "update"

		options := make(map[string]interface{})

		version, _ := cmd.Flags().GetString("version")
		if version != "" {
			options["version"] = version
		}

		backup, _ := cmd.Flags().GetBool("backup")
		if cmd.Flags().Changed("backup") {
			options["backup"] = backup
		}

		ProgressSpinner.Start()
		resp, err := helper.GenericJSONPostTimeout(section, command, options, helper.ContainerDownloadTimeout)
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
	coreUpdateCmd.Flags().StringP("version", "", "", "Version to update to")
	coreUpdateCmd.Flags().Bool("backup", false, "Create partial backup before update")
	coreUpdateCmd.RegisterFlagCompletionFunc("version", cobra.NoFileCompletions)
	coreUpdateCmd.RegisterFlagCompletionFunc("backup", boolCompletions)
	coreCmd.AddCommand(coreUpdateCmd)
}
