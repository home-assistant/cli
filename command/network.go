package command

import (
    "github.com/urfave/cli"
    "fmt"
    "github.com/home-assistant/hassio-cli/command/helpers"
    "os"
)

// CmdNetwork All network endpoints for hass.io
func CmdNetwork(c *cli.Context) {
    const HassioBasePath = "network"
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
    case "info":        // GET
        endpoint = action
        get = true
    case "options":     // POST
        endpoint = action
    default:
        fmt.Fprintf(os.Stderr, "No valid action detected")
        os.Exit(3)
    }

    if DebugEnabled {
        fmt.Fprintf(os.Stdout, "DEBUG [CmdNetwork]: action->'%s', endpoint='%s', serverOverride->'%s', GET->'%t', options->'%s', rawjson->'%t', filter->'%s'\n",
            action, endpoint, serverOverride, get, Options, RawJSON, Filter )
    }

    if endpoint != "" {
        helpers.ExecCommand(HassioBasePath, endpoint, serverOverride, get,  Options, Filter, RawJSON)
    }
}