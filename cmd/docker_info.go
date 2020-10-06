package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var dockerInfoCmd = &cobra.Command{
	Use:     "info",
	Aliases: []string{"in", "inf"},
	Short:   "Shows information about the host docker backend",
	Long: `
Shows information about the local docker backend on the host system
`,
	Example: `
  ha docker info
`,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("docker info")

		section := "docker"
		command := "info"
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
	dockerCmd.AddCommand(dockerInfoCmd)
}
