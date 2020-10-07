package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		log.WithField("args", args).Debug("docker registries")

		section := "docker"
		command := "registries"
		base := viper.GetString("endpoint")

		resp, err := helper.GenericJSONGet(base, section, command)
		if err != nil {
			fmt.Println(err)
			ExitWithError = true
		} else {
			ExitWithError = !helper.ShowJSONResponse(resp)
		}
	},
}

func init() {
	dockerCmd.AddCommand(dockerRegistriesCmd)
}
