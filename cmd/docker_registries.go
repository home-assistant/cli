package cmd

import (
	"log/slog"

	helper "github.com/home-assistant/cli/client"
	"github.com/spf13/cobra"
)

var dockerRegistriesCmd = &cobra.Command{
	Use:     "registries",
	Aliases: []string{"reg", "re"},
	Short:   "Manage private OCI docker registry",
	Long: `
Manage private OCI registry server on the local Docker host.
`,
	Example: `
	ha docker registries
	ha docker registries add my-docker.example.com
	ha docker registries delete my-docker.example.com
`,
	Run: func(cmd *cobra.Command, args []string) {
		slog.Debug("docker registries", "args", args)

		section := "docker"
		command := "registries"

		resp, err := helper.GenericJSONGet(section, command)
		if err != nil {
			helper.PrintError(err)
			ExitWithError = true
		} else {
			ExitWithError = !helper.ShowJSONResponse(resp)
		}
	},
}

func init() {
	dockerCmd.AddCommand(dockerRegistriesCmd)
}
