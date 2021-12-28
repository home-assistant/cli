package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var dockerInfoCmd = &cobra.Command{
	Use:     "info",
	Aliases: []string{"in", "inf"},
	Short:   "Shows information about the host docker backend",
	Long: `
Shows information about the local Docker backend on the host system
`,
	Example: `
  ha docker info
`,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("docker info")

		section := "docker"
		command := "info"

		resp, err := helper.GenericJSONGet(section, command)
		if err != nil {
			fmt.Println(err)
			ExitWithError = true
		} else {
			ExitWithError = !helper.ShowJSONResponse(resp)
		}
	},
}

func init() {
	dockerCmd.AddCommand(dockerInfoCmd)
}
