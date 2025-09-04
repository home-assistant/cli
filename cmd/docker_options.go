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
  ha docker options --enable-ipv6=true
  ha docker options --mtu=1450`,
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

		if cmd.Flags().Changed("mtu") {
			mtu, err := cmd.Flags().GetInt("mtu")
			if err == nil {
				if mtu == 0 {
					options["mtu"] = nil
				} else if mtu < 68 || mtu > 65535 {
					helper.PrintError(fmt.Errorf("MTU value must be between 68 and 65535, or 0 to reset"))
					ExitWithError = true
					return
				} else {
					options["mtu"] = mtu
				}
			}
		}

		resp, err := helper.GenericJSONPost(section, command, options)
		if err != nil {
			helper.PrintError(err)
			ExitWithError = true
		} else {
			if cmd.Flags().Changed("enable-ipv6") {
				fmt.Println("Note: System restart required to apply new IPv6 configuration.")
			}
			if cmd.Flags().Changed("mtu") {
				fmt.Println("Note: System restart required to apply new MTU configuration.")
			}
			ExitWithError = !helper.ShowJSONResponse(resp)
		}
	},
}

func init() {
	dockerOptionsCmd.Flags().BoolP("enable-ipv6", "", false, "Enable IPv6")
	dockerOptionsCmd.Flags().IntP("mtu", "", 0, "Set Docker MTU (68-65535, 0 to reset)")
	dockerOptionsCmd.Flags().SetNormalizeFunc(func(set *pflag.FlagSet, name string) pflag.NormalizedName {
		return pflag.NormalizedName(strings.ReplaceAll(name, "_", "-"))
	})

	dockerCmd.AddCommand(dockerOptionsCmd)
}
