package command

import (
	"fmt"

	"github.com/home-assistant/hassio-cli/command/helpers"
	"github.com/urfave/cli"
	"os"
	"bytes"
)

func CmdHomeassistant(c *cli.Context) {
	const HASSIO_BASE_PATH = "homeassistant"
	action := ""
	endpoint := ""
	get := false
	if c.NArg() > 0 {
		action = c.Args()[0]
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
	case "rollback":
		rollback(HASSIO_BASE_PATH, c.Bool("json"))
	default:
		fmt.Fprintf(os.Stderr, "No valid action detected")
		os.Exit(3)
	}

	if endpoint != "" {
		response := helpers.RestCall(HASSIO_BASE_PATH, endpoint, get, c.String("payload"))
		helpers.DisplayOutput(response, c.Bool("json"))
	}
}

func rollback(basepath string, jsonout bool) {
	info := helpers.StrToMap(helpers.RestCall(basepath, "info", true, ""))

	data := info["data"].(map[string]string)
	currentVersion := data["version"]
	lastVersion := data["last_version"]

	var payload bytes.Buffer
	payload.WriteString("version=")
	payload.WriteString(lastVersion)

	if currentVersion != lastVersion {
		response := helpers.RestCall(basepath, "update", false, payload.String())
		helpers.DisplayOutput(response, jsonout)
	} else {
		fmt.Fprintf(os.Stderr, "Current Version Matches last version, %s : %s", currentVersion, lastVersion)
		os.Exit(3)
	}
}
