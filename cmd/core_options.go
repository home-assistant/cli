package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("core options")

		section := "core"
		command := "options"
		base := viper.GetString("endpoint")

		options := make(map[string]interface{})

		for _, value := range []string{
			"image",
			"last_version",
			"password",
			"refresh_token",
		} {
			val, err := cmd.Flags().GetString(value)
			if val != "" && err == nil && cmd.Flags().Changed(value) {
				options[value] = val
			}
		}

		port, err := cmd.Flags().GetInt("port")
		if port != 0 && err == nil && cmd.Flags().Changed("port") {
			options["port"] = port
		}

		waitBoot, err := cmd.Flags().GetInt("wait_boot")
		if err == nil && cmd.Flags().Changed("wait_boot") {
			options["wait_boot"] = waitBoot
		}

		ssl, err := cmd.Flags().GetBool("ssl")
		if err == nil && cmd.Flags().Changed("ssl") {
			options["ssl"] = ssl
		}

		watchdog, err := cmd.Flags().GetBool("watchdog")
		if err == nil && cmd.Flags().Changed("watchdog") {
			options["watchdog"] = watchdog
		}

		resp, err := helper.GenericJSONPost(base, section, command, options)
		if err != nil {
			fmt.Println(err)
			ExitWithError = true
		} else {
			ExitWithError = !helper.ShowJSONResponse(resp)
		}
	},
}

func init() {
	coreOptionsCmd.Flags().String("image", "", "Optional image")
	coreOptionsCmd.Flags().String("last_version", "", "Optional for custom image")
	coreOptionsCmd.Flags().Int("port", 8123, "Port for access Home Assistant")
	coreOptionsCmd.Flags().Bool("ssl", false, "Use SSL")
	coreOptionsCmd.Flags().String("password", "", "API password")
	coreOptionsCmd.Flags().String("refresh_token", "", "Refresh token")
	coreOptionsCmd.Flags().Bool("watchdog", true, "Use watchdog")
	coreOptionsCmd.Flags().Int("wait_boot", 600, "wait_boot")
	coreCmd.AddCommand(coreOptionsCmd)
}
