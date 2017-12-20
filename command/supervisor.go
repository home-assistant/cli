package command

import (
	"github.com/urfave/cli"
	"fmt"
	"github.com/home-assistant/hassio-cli/command/helpers"
)

func CmdSupervisor(c *cli.Context) {
	const HASSIO_BASE_PATH = "supervisor"
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
	case "reload",     // POST
		"update":
		endpoint = action
	default:
		fmt.Println("No action detected")
	}

	if endpoint != "" {
		response := helpers.RestCall(HASSIO_BASE_PATH, endpoint, get, payload)

		helpers.DisplayOutput(response, c.Bool("json"))
	}
}
