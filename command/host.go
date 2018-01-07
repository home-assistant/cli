package command

import (
    "github.com/urfave/cli"
    "fmt"
    "github.com/home-assistant/hassio-cli/command/helpers"
    "os"
)

// CmdHost All host endpoints for hass.io
func CmdHost(c *cli.Context) {
    const HassioBasePath = "host"
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

    if DebugEnabled {
        fmt.Fprintf(os.Stdout, "DEBUG [CmdHost]: action->'%s', endpoint='%s', serverOverride->'%s', GET->'%t', options->'%s', rawjson->'%t', filter->'%s'\n",
            action, endpoint, serverOverride, get, Options, RawJSON, Filter )
    }

    if endpoint != "" {
        helpers.ExecCommand(HassioBasePath, endpoint, serverOverride, get,  Options, Filter, RawJSON)
    }
}