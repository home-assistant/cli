package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var authResetCmd = &cobra.Command{
	Use:     "reset",
	Aliases: []string{"rst", "change"},
	Short:   "Reset the password of a Home Assistant user.",
	Long: `
This command allows you to change a password of a Home Assistant user.
Please note, this command is limited due to security reasons, and will
only work on some locations. For example, the Operating System CLI.
`,
	Example: `
  ha authentication reset --username "JohnDoe" --password "123SuperSecret!"
`,
	ValidArgsFunction: cobra.NoFileCompletions,
	Args:              cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("auth reset")

		section := "auth"
		command := "reset"

		options := make(map[string]any)

		for _, value := range []string{
			"username",
			"password",
		} {
			val, err := cmd.Flags().GetString(value)
			if val != "" && err == nil && cmd.Flags().Changed(value) {
				options[value] = val
			}
		}

		resp, err := helper.GenericJSONPost(section, command, options)
		if err != nil {
			cmd.PrintErrln("this command is limited due to security reasons, and will only work on some locations. For example, the Operating System terminal.")
			fmt.Println(err)
			ExitWithError = true
		} else {
			ExitWithError = !helper.ShowJSONResponse(resp)
		}
	},
}

func init() {
	authResetCmd.Flags().String("username", "", "Username to reset the password for")
	authResetCmd.Flags().String("password", "", "The new password to assign")
	authResetCmd.MarkFlagRequired("username")
	authResetCmd.MarkFlagRequired("password")
	authResetCmd.RegisterFlagCompletionFunc("username", cobra.NoFileCompletions)
	authResetCmd.RegisterFlagCompletionFunc("password", cobra.NoFileCompletions)
	authCmd.AddCommand(authResetCmd)
}
