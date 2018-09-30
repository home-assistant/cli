package command

import (
	"fmt"

	"github.com/home-assistant/hassio-cli/command/helpers"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

// CmdHardware All hardware endpoints for hass.io
func CmdHardware(c *cli.Context) {
	const HassioBasePath = "hardware"
	action := ""
	endpoint := ""
	serverOverride := ""
	get := false
	DebugEnabled := c.GlobalBool("debug")
	helpers.DebugEnabled = DebugEnabled
	Options := ""
	RawJSON := c.Bool("rawjson")
	Filter := c.String("filter")
	if c.NArg() > 0 {
		action = c.Args()[0]
	}

	switch action {
	case "info", // GET
		"audio":
		get = true
		endpoint = action
	default:
		log.Error("No valid action detected.\n")
	}

	if DebugEnabled {
		infoMessage := fmt.Sprintf("DEBUG [CmdHardware]: action->'%s', endpoint='%s', serverOverride->'%s', GET->'%t', options->'%s', rawjson->'%t', filter->'%s'\n",
			action, endpoint, serverOverride, get, Options, RawJSON, Filter)
		log.Info(infoMessage)
	}

	if endpoint != "" {
		helpers.ExecCommand(HassioBasePath, endpoint, serverOverride, get, Options, Filter, RawJSON)
	}
}
