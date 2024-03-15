package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var osRollbackCmd = &cobra.Command{
	Use:     "rollback",
	Aliases: []string{"revert", "rb"},
	Short:   "Rollback OS after an update",
	Long: `
Use this command after an OS update to switch to the previous boot slot
with the previous version of OS installed.
`,
	Example: `
  ha os rollback
`,
	ValidArgsFunction: cobra.NoFileCompletions,
	Args:              cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("os rollback")

		section := "os"
		command := "boot-slot"

		bootSlots, err := osGetBootSlots()
		if err != nil {
			fmt.Println(err)
			ExitWithError = true
			return
		}
		if bootSlots == nil {
			ExitWithError = true
			return
		}

		target := ""
		for bootSlot, v := range bootSlots {
			info, ok := v.(map[string]interface{})
			if !ok {
				continue
			}
			if state, ok := info["state"].(string); ok && state == "inactive" {
				target = bootSlot
				break
			}
		}

		options := make(map[string]interface{})
		options = map[string]interface{}{"boot_slot": target}

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
	osCmd.AddCommand(osRollbackCmd)
}
