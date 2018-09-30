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

	log.WithFields(log.Fields{
		"action":         action,
		"endpoint":       endpoint,
		"serverOverride": serverOverride,
		"get":            get,
		"options":        Options,
		"rawJSON":        RawJSON,
		"filter":         Filter,
	}).Debug("[CmdHomeassistant]")

	if endpoint != "" {
		helpers.ExecCommand(HassioBasePath, endpoint, serverOverride, get, Options, Filter, RawJSON)
	}
}
