package cmd

import (
	"log/slog"

	helper "github.com/home-assistant/cli/client"
	"github.com/spf13/cobra"
)

var dockerRegistriesAddCmd = &cobra.Command{
	Use:     "add [host]",
	Aliases: []string{"set", "new"},
	Short:   "Add new docker registry login for specific host",
	Long: `
Add new login for the Docker OCI registry server.
`,
	Example: `
  ha docker registries add my-docker.example.com --username "test" --password "example"
`,
	ValidArgsFunction: cobra.NoFileCompletions,
	Args:              cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		slog.Debug("registries add", "args", args)

		section := "docker"
		command := "registries"

		url, err := helper.URLHelper(section, command)
		if err != nil {
			helper.PrintError(err)
			ExitWithError = true
			return
		}

		options := map[string]map[string]string{}

		request := helper.GetJSONRequest()

		host := args[0]
		options[host] = make(map[string]string)

		for _, value := range []string{
			"username",
			"password",
		} {
			val, err := cmd.Flags().GetString(value)
			if val != "" && err == nil && cmd.Flags().Changed(value) {
				options[host][value] = val
			}
		}

		if len(options) > 0 {
			slog.Debug("Request body", "options", options)
			request.SetBody(options)
		}

		resp, err := request.Post(url)
		resp, err = helper.GenericJSONErrorHandling(resp, err)

		if err != nil {
			helper.PrintError(err)
			ExitWithError = true
		} else {
			ExitWithError = !helper.ShowJSONResponse(resp)
		}
	},
}

func init() {
	dockerRegistriesAddCmd.Flags().StringP("username", "u", "", "Username for OCI auth")
	dockerRegistriesAddCmd.Flags().StringP("password", "p", "", "Password for OCI auth")
	dockerRegistriesAddCmd.RegisterFlagCompletionFunc("username", cobra.NoFileCompletions)
	dockerRegistriesAddCmd.RegisterFlagCompletionFunc("password", cobra.NoFileCompletions)
	dockerRegistriesCmd.AddCommand(dockerRegistriesAddCmd)
}
