package command

import (
	"github.com/urfave/cli"
	"fmt"
	"github.com/home-assistant/hassio-cli/command/helpers"
	"os"
)

func CmdHost(c *cli.Context) {
	const HASSIO_BASE_PATH = "host"
	action := ""
	endpoint := ""
	get := false
	if c.NArg() > 0 {
		action = c.Args()[0]
	}

	switch action {
	case "hardware":      // GET
		endpoint = action
		get = true
	case "reboot",     // POST
		 "update",
		 "shutdown":
		endpoint = action
	default:
		fmt.Fprintf(os.Stderr, "No valid action detected")
		os.Exit(3)
	}

	if endpoint != "" {
		response := helpers.RestCall(HASSIO_BASE_PATH, endpoint, get, c.String("options"))
		helpers.DisplayOutput(response, c.Bool("rawjson"))
	}
}
