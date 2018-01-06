package command

import (
    "fmt"
    "github.com/home-assistant/hassio-cli/command/helpers"
    "github.com/urfave/cli"
    "os"
)


// CmdHomeassistant All home-assistant endpoints for hass.io
func CmdHomeassistant(c *cli.Context) {
    const HassioBasePath = "homeassistant"
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
        fmt.Fprintf(os.Stdout, "No valid action detected")
        os.Exit(3)
    }

    if DebugEnabled {
        fmt.Fprintf(os.Stdout, "DEBUG [CmdHomeassistant]: action->'%s', endpoint='%s', serverOverride->'%s', GET->'%t', options->'%s', rawjson->'%t', filter->'%s'\n",
            action, endpoint, serverOverride, get, Options, RawJSON, Filter )
    }
    if endpoint != "" {
        helpers.ExecCommand(HassioBasePath, endpoint, serverOverride, get,  Options, Filter, RawJSON)
    }
}