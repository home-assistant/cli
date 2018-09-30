package command

import (
	"fmt"

	"github.com/home-assistant/hassio-cli/command/helpers"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

// CmdSnapshots snapshot endpoints for hass.io
func CmdSnapshots(c *cli.Context) {
	const HassioBasePath = "snapshots"
	action := ""
	endpoint := ""
	serverOverride := ""
	get := false
	DebugEnabled := c.GlobalBool("debug")
	helpers.DebugEnabled = DebugEnabled
	Options := c.String("options")
	RawJSON := c.Bool("rawjson")
	Filter := c.String("filter")
	SnapSlug := c.String("slug")
	if c.NArg() > 0 {
		action = c.Args()[0]
	}

	var errorMessage string
	switch action {
	case "list": // GET
		get = true
	case "info":
		if SnapSlug == "" {
			errorMessage = fmt.Sprintf("--slug is required. See '%s --help'.\n", c.App.Name)
			log.Error(errorMessage)
		}
		get = true
		endpoint = SnapSlug + "/info"
	case "reload": // POST
		endpoint = action
	case "new":
		endpoint = "new/full"
		if c.String("name") != "" {
			if Options != "" {
				Options += ","
			}
			Options += "name=" + c.String("name")
		}
		if c.String("password") != "" {
			if Options != "" {
				Options += ","
			}
			Options += "password=" + c.String("password")
		}
	case "restore":
		if SnapSlug == "" {
			errorMessage = fmt.Sprintf("--slug is required. See '%s --help'.\n", c.App.Name)
			log.Error(errorMessage)
		}
		if c.String("password") != "" {
			Options = "password=" + c.String("password")
		}
		endpoint = SnapSlug + "/restore/full"
	case "remove":
		if SnapSlug == "" {
			errorMessage = fmt.Sprintf("--slug is required. See '%s --help'.\n", c.App.Name)
			log.Error(errorMessage)
		}
		endpoint = SnapSlug + "/remove"
	default:
		log.Error("No valid action detected.\n")
	}

	if endpoint != "" || action == "list" {
		helpers.ExecCommand(HassioBasePath, endpoint, serverOverride, get, Options, Filter, RawJSON)
	}
}
