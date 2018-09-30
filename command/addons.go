package command

import (
	"fmt"

	"github.com/home-assistant/hassio-cli/command/helpers"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

// CmdAddons All addon endpoints for hass.io
func CmdAddons(c *cli.Context) {
	const HassioBasePath = "addons"
	action := ""
	endpoint := ""
	serverOverride := ""
	get := false
	DebugEnabled := c.GlobalBool("debug")
	helpers.DebugEnabled = DebugEnabled
	Options := c.String("options")
	RawJSON := c.Bool("rawjson")
	Filter := c.String("filter")
	AddonName := c.String("name")
	if c.NArg() > 0 {
		action = c.Args()[0]
	}

	var errorMessage string
	switch action {
	case "list": // GET
		endpoint = ""
		get = true
	case "info",
		"logo",
		"changelog",
		"logs", // Fix as not JSON format for output
		"stats":
		if AddonName == "" {
			errorMessage = fmt.Sprintf("--name is required. See '%s --help'.\n", c.App.Name)
			log.Error(errorMessage)
		}
		endpoint = AddonName + "/" + action
		get = true
	case "reload": // POST
		endpoint = action
	case "start",
		"stop",
		"install",
		"uninstall",
		"update":
		if AddonName == "" {
			errorMessage = fmt.Sprintf("--name is required. See '%s --help'.\n", c.App.Name)
			log.Error(errorMessage)
		}
		endpoint = AddonName + "/" + action
	default:
		log.Error("No valid action detected.\n")
	}

	if DebugEnabled {
		infoMessage := fmt.Sprintf("DEBUG [addons]: action->'%s', endpoint='%s', serverOverride->'%s', GET->'%t', options->'%s', rawjson->'%t', filter->'%s'\n",
			action, endpoint, serverOverride, get, Options, RawJSON, Filter)
		log.Info(infoMessage)
	}
	if endpoint != "" || action == "list" {
		helpers.ExecCommand(HassioBasePath, endpoint, serverOverride, get, Options, Filter, RawJSON)
	}
}
