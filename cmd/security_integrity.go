package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var securityIntegrityCmd = &cobra.Command{
	Use:     "integrity",
	Aliases: []string{"int", "trust"},
	Short:   "Execute security integrity check",
	Long: `
This command execute a full system integrity check.
This need content trust to be enabled.`,
	Example: `
  ha security integrity`,
	ValidArgsFunction: cobra.NoFileCompletions,
	Args:              cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("security")

		section := "security"
		command := "integrity"

		ProgressSpinner.Start()
		resp, err := helper.GenericJSONPost(section, command, nil)
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
	// add cmd to root command
	securityCmd.AddCommand(securityIntegrityCmd)
}
