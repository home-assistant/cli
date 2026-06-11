package cmd

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	resty "github.com/go-resty/resty/v2"
	helper "github.com/home-assistant/cli/client"
	"github.com/spf13/cobra"
)

// getRPiFirmwareInfo fetches the Raspberry Pi firmware status. A 404 is turned
// into a friendly error in case the API doesn't return it already.
func getRPiFirmwareInfo() (*resty.Response, error) {
	url, err := helper.URLHelper("os", "boards/raspberrypi/firmware")
	if err != nil {
		return nil, err
	}

	resp, err := helper.GetJSONRequestTimeout(helper.DefaultTimeout).Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() == http.StatusNotFound {
		// Prefer the Supervisor's own message when it returns a JSON body,
		// otherwise (e.g. a bare 404 from an older Supervisor) explain why.
		if data, ok := resp.Error().(*helper.Response); ok && data.Message != "" {
			return nil, errors.New(data.Message)
		}
		return nil, errors.New("firmware information is not available on this system " +
			"(requires a Raspberry Pi 4 / 5 or Home Assistant Yellow with Home Assistant OS 18.0 or newer)")
	}

	return helper.GenericJSONErrorHandling(resp, err)
}

// humanizeRPiFirmwareVersion turns a raw firmware version into a human-readable
// string. The Supervisor reports the bootloader EEPROM build as a Unix timestamp,
// optionally suffixed with the VL805 EEPROM revision (timestamp-hexstring). It
// renders the timestamp as a UTC YYYY-MM-DD date and appends (VL805 hexstring)
// when a VL805 revision is present.
func humanizeRPiFirmwareVersion(version string) string {
	if version == "" {
		return version
	}
	timestamp, vl805, _ := strings.Cut(version, "-")
	secs, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return version
	}
	date := time.Unix(secs, 0).UTC().Format("2006-01-02")
	if vl805 != "" {
		return fmt.Sprintf("%s (VL805 %s)", date, vl805)
	}
	return date
}

// humanizeRPiFirmwareData humanizes the version fields of a firmware info response
// in place, leaving non-string values (e.g. null) untouched.
func humanizeRPiFirmwareData(data map[string]any) {
	for _, key := range []string{"current_version", "latest_version"} {
		if v, ok := data[key].(string); ok {
			data[key] = humanizeRPiFirmwareVersion(v)
		}
	}
}

var osBoardsRaspberrypiFirmwareCmd = &cobra.Command{
	Use:     "firmware",
	Aliases: []string{"fw"},
	Short:   "Show Raspberry Pi firmware information",
	Long: `
This command shows information about the Raspberry Pi bootloader (EEPROM) firmware,
including the currently installed and latest bundled versions and whether an update
is available. Available on Raspberry Pi 4 / 5 and Home Assistant Yellow.`,
	Example: `
  ha os boards raspberrypi firmware`,
	ValidArgsFunction: cobra.NoFileCompletions,
	Args:              cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		slog.Debug("os boards raspberrypi firmware", "args", args)

		resp, err := getRPiFirmwareInfo()
		if err != nil {
			helper.PrintError(err)
			ExitWithError = true
		} else {
			ExitWithError = !helper.ShowJSONResponseTransform(resp, humanizeRPiFirmwareData)
		}
	},
}

func init() {
	osBoardsRaspberrypiCmd.AddCommand(osBoardsRaspberrypiFirmwareCmd)
}
