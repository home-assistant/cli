package command

import (
	"fmt"

	"github.com/home-assistant/hassio-cli/command/helpers"
	"github.com/urfave/cli"
)

func CmdHomeassistant(c *cli.Context) {
	const HASSIO_BASE_PATH = "homeassistant"
	action := ""
	endpoint := ""
	payload := ""
	get := false
	if c.NArg() > 0 {
		action = c.Args()[0]
	}
	if c.NArg() == 2 {
		payload = c.Args()[1]
	}

	switch action {
	case "info",      // GET
		 "logs":
		endpoint = action
		get = true
	case "check",     // POST
	     "restart",
	     "start",
	     "stop",
	     "update":
		endpoint = action
	default:
		fmt.Println("No action detected")
	}

	if endpoint != "" {
		response := helpers.RestCall(HASSIO_BASE_PATH, endpoint, get, payload)
		helpers.DisplayOutput(helpers.MapToJSON(response), c.Bool("json"))
	}
}


