package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var coreOptionsCmd = &cobra.Command{
	Use:     "options",
	Aliases: []string{"option", "opt", "opts", "op"},
	Short:   "Allow to set options on Home Assistant Core instance",
	Long: `
This command allows you to set configuration options for the Home Assistant Core
instance running on your Home Assistant system.`,
	Example: `
  ha core options --wait_boot 600`,
	ValidArgsFunction: cobra.NoFileCompletions,
	Args:              cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("core options")

		section := "core"
		command := "options"

		options := make(map[string]interface{})

		for _, value := range []string{
			"image",
			"refresh_token",
			"audio_output",
			"audio_input",
		} {
			val, err := cmd.Flags().GetString(value)
			if err == nil && cmd.Flags().Changed(value) {
				if val == "" {
					options[value] = nil
				} else {
					options[value] = val
				}
			}
		}

		port, err := cmd.Flags().GetInt("port")
		if port != 0 && err == nil && cmd.Flags().Changed("port") {
			options["port"] = port
		}

		for _, value := range []string{
			"boot",
			"ssl",
			"watchdog",
		} {
			val, err := cmd.Flags().GetBool(value)
			if err == nil && cmd.Flags().Changed(value) {
				options[value] = val
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
	coreOptionsCmd.Flags().Bool("boot", true, "Start Core on boot")
	coreOptionsCmd.Flags().String("image", "", "Optional image")
	coreOptionsCmd.Flags().Int("port", 8123, "Port to access Home Assistant Core")
	coreOptionsCmd.Flags().Bool("ssl", false, "Use SSL")
	coreOptionsCmd.Flags().Bool("watchdog", true, "Use watchdog")
	coreOptionsCmd.Flags().String("refresh_token", "", "Refresh token")
	coreOptionsCmd.Flags().String("audio_input", "", "Profile name for audio input")
	coreOptionsCmd.Flags().String("audio_output", "", "Profile name for audio output")
	coreOptionsCmd.RegisterFlagCompletionFunc("boot", boolCompletions)
	coreOptionsCmd.RegisterFlagCompletionFunc("image", cobra.NoFileCompletions)
	coreOptionsCmd.RegisterFlagCompletionFunc("port", cobra.NoFileCompletions)
	coreOptionsCmd.RegisterFlagCompletionFunc("ssl", boolCompletions)
	coreOptionsCmd.RegisterFlagCompletionFunc("watchdog", boolCompletions)
	coreOptionsCmd.RegisterFlagCompletionFunc("refresh_token", cobra.NoFileCompletions)
	coreOptionsCmd.RegisterFlagCompletionFunc("audio_input", cobra.NoFileCompletions)
	coreOptionsCmd.RegisterFlagCompletionFunc("audio_output", cobra.NoFileCompletions)
	coreCmd.AddCommand(coreOptionsCmd)
}
