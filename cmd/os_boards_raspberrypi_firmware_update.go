package cmd

import (
	"fmt"
	"log/slog"

	helper "github.com/home-assistant/cli/client"
	"github.com/spf13/cobra"
)

var osBoardsRaspberrypiFirmwareUpdateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"up", "upgrade"},
	Short:   "Apply the bundled Raspberry Pi firmware update",
	Long: `
This command applies the bundled Raspberry Pi firmware update (bootloader EEPROM,
and VL805 where present). A reboot is required for the new firmware to take effect.`,
	Example: `
  ha os boards raspberrypi firmware update`,
	ValidArgsFunction: cobra.NoFileCompletions,
	Args:              cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		slog.Debug("os boards raspberrypi firmware update", "args", args)

		// Fetch the current firmware status first so we can fail gracefully.
		infoResp, err := getRPiFirmwareInfo()
		if err != nil {
			helper.PrintError(err)
			ExitWithError = true
			return
		}
		if !infoResp.IsSuccess() {
			// Other API errors (e.g. a 400) come back with a message to show.
			ExitWithError = true
			helper.ShowJSONResponse(infoResp)
			return
		}

		info := infoResp.Result().(*helper.Response)

		if blocked, _ := info.Data["update_blocked"].(bool); blocked {
			msg := "firmware update is blocked"
			if reason, _ := info.Data["blocked_reason"].(string); reason != "" {
				msg = fmt.Sprintf("%s: %s", msg, reason)
			}
			helper.PrintErrorString(msg)
			ExitWithError = true
			return
		}

		if pending, _ := info.Data["update_pending"].(bool); pending {
			fmt.Println("A firmware update is already pending. Reboot the device using `ha host reboot` to finish the installation.")
			return
		}

		if available, ok := info.Data["update_available"].(bool); ok && !available {
			msg := "The Raspberry Pi firmware is already up to date."
			if current, _ := info.Data["current_version"].(string); current != "" {
				msg = fmt.Sprintf("The Raspberry Pi firmware is already up to date (%s).", humanizeRPiFirmwareVersion(current))
			}
			fmt.Println(msg)
			return
		}

		prompt := "This applies the bundled Raspberry Pi firmware update and requires a reboot to take effect."
		current, _ := info.Data["current_version"].(string)
		latest, _ := info.Data["latest_version"].(string)
		if current != "" && latest != "" {
			prompt = fmt.Sprintf("This will update the Raspberry Pi firmware from %s to %s and requires a reboot to take effect.",
				humanizeRPiFirmwareVersion(current), humanizeRPiFirmwareVersion(latest))
		}

		confirmed, err := helper.AskForConfirmation(prompt+" Continue?", 2)
		if err != nil {
			helper.PrintError(err)
			ExitWithError = true
			return
		}
		if !confirmed {
			fmt.Println("Aborted.")
			return
		}

		fmt.Println("\nApplying firmware update. Please be patient and do not unplug or power off the device.")

		ProgressSpinner.Start()
		resp, err := helper.GenericJSONPost("os", "boards/raspberrypi/firmware/update", make(map[string]any))
		ProgressSpinner.Stop()
		if err != nil {
			helper.PrintError(err)
			ExitWithError = true
		} else if helper.ShowJSONResponse(resp) {
			fmt.Println("\nFirmware update applied. Reboot the device using `ha host reboot` to finish the installation.")
		} else {
			ExitWithError = true
		}
	},
}

func init() {
	osBoardsRaspberrypiFirmwareCmd.AddCommand(osBoardsRaspberrypiFirmwareUpdateCmd)
}
