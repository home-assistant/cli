package cmd

import (
	"fmt"
	"strings"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var supervisorOptionsCmd = &cobra.Command{
	Use:     "options",
	Aliases: []string{"option", "opt", "opts", "op"},
	Short:   "Allows you to set options on the Home Assistant Supervisor",
	Long: `
This command allows you to set configuration options for on the Home Assistant
Supervisor running on your Home Assistant system.`,
	Example: `
  ha supervisor options --channel beta`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("supervisor options")

		section := "supervisor"
		command := "options"

		options := make(map[string]interface{})

		for _, value := range []string{
			"hostname",
			"channel",
			"timezone",
			"logging",
		} {
			val, err := cmd.Flags().GetString(value)
			if val != "" && err == nil && cmd.Flags().Changed(value) {
				options[value] = val
			}
		}

		for _, value := range []string{
			"debug",
			"debug-block",
			"diagnostics",
			"content-trust",
			"force-security",
		} {
			data, err := cmd.Flags().GetBool(value)
			if err == nil && cmd.Flags().Changed(value) {
				options[strings.Replace(value, "-", "_", -1)] = data
			}
		}

		waitboot, _ := cmd.Flags().GetInt("wait-boot")
		if cmd.Flags().Changed("wait-boot") {
			options["wait_boot"] = waitboot
		}

		repos, err := cmd.Flags().GetStringArray("repositories")
		log.WithField("repositories", repos).Debug("repos")

		if len(repos) >= 1 && err == nil {
			options["addons_repositories"] = repos
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
	supervisorOptionsCmd.Flags().StringP("hostname", "", "", "Hostname to set")
	supervisorOptionsCmd.Flags().StringP("channel", "c", "", "Channel to track (stable|beta|dev)")
	supervisorOptionsCmd.Flags().StringP("timezone", "t", "", "Timezone")
	supervisorOptionsCmd.Flags().StringP("logging", "l", "", "Logging: debug|info|warning|error|critical")
	supervisorOptionsCmd.Flags().IntP("wait-boot", "w", 0, "Seconds to wait after boot")
	supervisorOptionsCmd.Flags().BoolP("debug", "", false, "Enable debug mode")
	supervisorOptionsCmd.Flags().BoolP("debug-block", "", false, "Enable debug mode with blocking startup")
	supervisorOptionsCmd.Flags().BoolP("diagnostics", "", false, "Enable diagnostics mode")
	supervisorOptionsCmd.Flags().StringArrayP("repositories", "r", []string{}, "repositories to track, can be supplied multiple times")

	supervisorOptionsCmd.Flags().Lookup("debug").NoOptDefVal = "false"
	supervisorOptionsCmd.Flags().Lookup("debug-block").NoOptDefVal = "false"
	supervisorOptionsCmd.Flags().Lookup("diagnostics").NoOptDefVal = "false"

	supervisorCmd.AddCommand(supervisorOptionsCmd)
}
