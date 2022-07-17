package cmd

import (
	"fmt"
	"strings"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var securityOptionsCmd = &cobra.Command{
	Use:     "options",
	Aliases: []string{"option", "opt", "opts", "op"},
	Short:   "Allow to set options for the Security backend",
	Long: `
This command allows you to set configuration options for the internally
Home Assistant Security backend.
`,
	Example: `
  ha security options --force-security=True
`,
	ValidArgsFunction: cobra.NoFileCompletions,
	Args:              cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("security options")

		section := "security"
		command := "options"

		options := make(map[string]interface{})

		for _, value := range []string{
			"pwned",
			"content-trust",
			"force-security",
		} {
			data, err := cmd.Flags().GetBool(value)
			if err == nil && cmd.Flags().Changed(value) {
				options[strings.Replace(value, "-", "_", -1)] = data
			}
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
	securityOptionsCmd.Flags().BoolP("pwned", "", true, "Enable/Disable pwned check the backend")
	securityOptionsCmd.Flags().BoolP("content-trust", "", true, "Enable/Disable content-trust on the backend")
	securityOptionsCmd.Flags().BoolP("force-security", "", false, "Enable/Disable force-security on the backend")

	securityOptionsCmd.Flags().Lookup("pwned").NoOptDefVal = "true"
	securityOptionsCmd.Flags().Lookup("content-trust").NoOptDefVal = "true"
	securityOptionsCmd.Flags().Lookup("force-security").NoOptDefVal = "false"

	securityOptionsCmd.RegisterFlagCompletionFunc("pwned", cobra.NoFileCompletions)
	securityOptionsCmd.RegisterFlagCompletionFunc("content-trust", cobra.NoFileCompletions)
	securityOptionsCmd.RegisterFlagCompletionFunc("force-security", cobra.NoFileCompletions)

	securityCmd.AddCommand(securityOptionsCmd)
}
