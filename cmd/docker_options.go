package cmd

import (
	"fmt"
	"strings"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var dockerOptionsCmd = &cobra.Command{
	Use:     "options",
	Aliases: []string{"option", "opt", "opts", "op"},
	Short:   "Allows you to set options on the host docker backend",
	Long: `
This command allows you to set configuration options for on the host
docker backend running on your Home Assistant system.`,
	Example: `
  ha docker options --enable-ipv6 true`,
	ValidArgsFunction: cobra.NoFileCompletions,
	Args:              cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("docker options")

		section := "docker"
		command := "options"

		options := make(map[string]any)

		for _, value := range []string{
			"enable_ipv6",
		} {
			data, err := cmd.Flags().GetBool(value)
			if err == nil && cmd.Flags().Changed(value) {
				options[strings.ReplaceAll(value, "-", "_")] = data
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
	dockerOptionsCmd.Flags().BoolP("enable-ipv6", "", false, "Enable IPv6")
	dockerOptionsCmd.Flags().SetNormalizeFunc(func(set *pflag.FlagSet, name string) pflag.NormalizedName {
		return pflag.NormalizedName(strings.ReplaceAll(name, "_", "-"))
	})

	dockerCmd.AddCommand(dockerOptionsCmd)
}
