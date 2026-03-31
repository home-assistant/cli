package cmd

import (
	"log/slog"

	helper "github.com/home-assistant/cli/client"
	"github.com/spf13/cobra"
)

var dockerRegistriesDeleteCmd = &cobra.Command{
	Use:     "delete [host]",
	Aliases: []string{"del", "remove"},
	Short:   "Delete docker registry login for specific host",
	Long: `
Remove login for the Docker OCI registry server.
`,
	Example: `
  ha docker registries delete my-docker.example.com"
`,
	ValidArgsFunction: dockerRegistriesDeleteCompletions,
	Args:              cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		slog.Debug("registries delete", "args", args)

		section := "docker"
		command := "registries/{host}"

		url, err := helper.URLHelper(section, command)
		if err != nil {
			helper.PrintError(err)
			ExitWithError = true
			return
		}

		request := helper.GetJSONRequest()

		host := args[0]

		request.SetPathParams(map[string]string{
			"host": host,
		})

		resp, err := request.Delete(url)
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
	dockerRegistriesCmd.AddCommand(dockerRegistriesDeleteCmd)
}

func dockerRegistriesDeleteCompletions(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	resp, err := helper.GenericJSONGet("docker", "registries")
	if err != nil || !resp.IsSuccess() {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	var ret []string
	data := resp.Result().(*helper.Response)
	if data.Result == "ok" && data.Data["registries"] != nil {
		if registries, ok := data.Data["registries"].(map[string]any); ok {
			for k := range registries {
				ret = append(ret, k)
			}
		}
	}
	return ret, cobra.ShellCompDirectiveNoFileComp
}
