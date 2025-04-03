package cmd

import (
	"fmt"
	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var osConfigSwapOptionsCmd = &cobra.Command{
	Use:     "options",
	Aliases: []string{"option", "opt", "opts", "op"},
	Short:   "Change HAOS swap settings",
	Long: `
This command allows you to override how the Home Assistant OS uses swap.`,
	Example: `
  ha os config swap options --swap-size=2G --swappiness=10`,
	ValidArgsFunction: cobra.NoFileCompletions,
	Args:              cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("os config swap options")

		section := "os"
		command := "config/swap"

		options := make(map[string]any)

		swapSize, err := cmd.Flags().GetString("swap-size")
		if err == nil && cmd.Flags().Changed("swap-size") {
			options["swap_size"] = swapSize
		}

		swappiness, err := cmd.Flags().GetInt("swappiness")
		if err == nil && cmd.Flags().Changed("swappiness") {
			options["swappiness"] = swappiness
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
	osConfigSwapOptionsCmd.Flags().String("swap-size", "", "Swap size in bytes with optional units (K/M/G)")
	osConfigSwapOptionsCmd.Flags().Int("swappiness", 1, "Kernel swappiness value (0-200)")

	osConfigSwapCmd.AddCommand(osConfigSwapOptionsCmd)
}
