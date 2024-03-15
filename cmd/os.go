package cmd

import (
	"os"

	helper "github.com/home-assistant/cli/client"

	"github.com/spf13/cobra"
)

var osCmd = &cobra.Command{
	Use:     "os",
	Aliases: []string{"hassos"},
	Short:   "Operating System specific for updating, info and configuration imports",
	Long: `
This command set is specifically designed for the Home Assistant Operating System
and only works on those systems. It provides an interface to get information
about your Home Assistant Operating System, but also provides command to
upgrade the operating system and the operating system CLI. Finally,
it provides a command to import configurations from a USB stick.`,
	Example: `
  ha os info
  ha os update`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		for _, arg := range os.Args {
			if arg == "hassos" {
				cmd.PrintErrf("The use of '%s' is deprecated, please use 'os' instead!\n", arg)
			}
		}
		rootCmd.PersistentPreRun(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(osCmd)
}

func osGetBootSlots() (map[string]interface{}, error) {
	resp, err := helper.GenericJSONGet("os", "info")
	if err != nil || !resp.IsSuccess() {
		return nil, err
	}

	data := resp.Result().(*helper.Response)
	if data.Result == "ok" && data.Data["boot_slots"] != nil {
		if bootSlots, ok := data.Data["boot_slots"].(map[string]interface{}); ok {
			return bootSlots, nil
		}
	}

	return nil, nil
}

func osBootSlotCompletions(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	var ret []string
	bootSlots, _ := osGetBootSlots()
	for bootSlot := range bootSlots {
		ret = append(ret, bootSlot)
	}
	return ret, cobra.ShellCompDirectiveNoFileComp
}
