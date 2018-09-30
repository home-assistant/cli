package command

import (
	"fmt"

	"github.com/home-assistant/hassio-cli/command/helpers"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

// CmdSupervisor All supervisor endpoints for hass.io
func CmdSupervisor(c *cli.Context) {
	const HassioBasePath = "supervisor"
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
	case "info", // GET
		"logs":
		endpoint = action
		get = true
	case "reload", // POST
		"update",
		"options":
		endpoint = action
	default:
		log.Error("No valid action detected.\n")
	}

	if DebugEnabled {
		infoMessage := fmt.Sprintf("DEBUG [CmdSupervisor]: action->'%s', endpoint='%s', serverOverride->'%s', GET->'%t', options->'%s', rawjson->'%t', filter->'%s'\n",
			action, endpoint, serverOverride, get, Options, RawJSON, Filter)
		log.Info(infoMessage)
	}

	if endpoint != "" {
		helpers.ExecCommand(HassioBasePath, endpoint, serverOverride, get, Options, Filter, RawJSON)
	}
}
