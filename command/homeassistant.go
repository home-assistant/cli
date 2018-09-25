package command

import (
	"fmt"
	"os"

	"github.com/home-assistant/hassio-cli/command/helpers"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

// CmdHomeassistant All home-assistant endpoints for hass.io
func CmdHomeassistant(c *cli.Context) {
	const HassioBasePath = "homeassistant"
	action := ""
	endpoint := ""
	serverOverride := ""
	get := false
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
	case "check", // POST
		"restart",
		"start",
		"stop",
		"update",
		"options":
		endpoint = action
	default:
		fmt.Fprintf(os.Stdout, "No valid action detected.\n")
		os.Exit(3)
	}

	log.Debugf("[CmdHomeassistant]: action->'%s', endpoint='%s', serverOverride->'%s', GET->'%t', options->'%s', rawjson->'%t', filter->'%s'\n",
		action, endpoint, serverOverride, get, Options, RawJSON, Filter)

	if endpoint != "" {
		helpers.ExecCommand(HassioBasePath, endpoint, serverOverride, get, Options, Filter, RawJSON)
	}
}
