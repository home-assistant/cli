package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var osBootSlotCmd = &cobra.Command{
	Use:     "boot-slot [slot]",
	Aliases: []string{"boot"},
	Short:   "Changes the active boot slot",
	Long: `
Using this command you can change the active boot slot to rollback
an OS update without making more changes to the system.
`,
	Example: `
  ha os boot-slot other
  ha os boot-slot A
`,
	ValidArgsFunction: osBootSlotCompletions,
	Args:              cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("os boot-slot")

		section := "os"
		command := "boot-slot"

		bootSlot := args[0]
		if bootSlot == "other" {
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

			for name, v := range bootSlots {
				info, ok := v.(map[string]any)
				if !ok {
					continue
				}
				if state, ok := info["state"].(string); ok && state == "inactive" {
					bootSlot = name
					break
				}
			}
			if bootSlot == "other" {
				ExitWithError = true
				return
			}
		}

		options := map[string]any{"boot_slot": bootSlot}
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
	osCmd.AddCommand(osBootSlotCmd)
}
