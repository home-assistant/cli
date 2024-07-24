package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var hostRebootCmd = &cobra.Command{
	Use:     "reboot",
	Aliases: []string{"restart", "rb"},
	Short:   "Reboots the host machine",
	Long: `
Reboot the machine that your Home Assistant is running on.`,
	Example: `
  ha host reboot`,
	ValidArgsFunction: cobra.NoFileCompletions,
	Args:              cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("host reboot")

		section := "host"
		command := "reboot"

		options := make(map[string]interface{})

		force, err := cmd.Flags().GetBool("force")
		if err == nil && force {
			options["force"] = force
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
	hostRebootCmd.Flags().BoolP("force", "f", false, "Force reboot during an offline db migration")
	hostRebootCmd.Flags().Lookup("force").NoOptDefVal = "true"
	hostRebootCmd.RegisterFlagCompletionFunc("force", boolCompletions)
	hostCmd.AddCommand(hostRebootCmd)
}
