package command

import (
	"fmt"
	"os"

	"github.com/home-assistant/hassio-cli/command/helpers"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

// CmdHassOS All host endpoints for hass.io
func CmdHassOS(c *cli.Context) {
	const HassioBasePath = "host"
	action := ""
	endpoint := ""
	serverOverride := ""
	get := false
	DebugEnabled := c.GlobalBool("debug")
	helpers.DebugEnabled = DebugEnabled
	Options := c.String("options")
	RawJSON := c.Bool("rawjson")
	Filter := c.String("filter")
	if c.NArg() > 0 {
		action = c.Args()[0]
	}

	switch action {
	case "info": // GET
		get = true
		endpoint = action
	default:
		log.Error("No valid action detected.\n")
	}

	if DebugEnabled {
		infoMessage := fmt.Sprintf(os.Stdout, "DEBUG [CmdHost]: action->'%s', endpoint='%s', serverOverride->'%s', GET->'%t', options->'%s', rawjson->'%t', filter->'%s'\n",
			action, endpoint, serverOverride, get, Options, RawJSON, Filter)
		log.Info(infoMessage)
	}

	if endpoint != "" {
		helpers.ExecCommand(HassioBasePath, endpoint, serverOverride, get, Options, Filter, RawJSON)
	}
}
