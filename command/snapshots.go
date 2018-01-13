package command

import (
    "github.com/urfave/cli"
    "fmt"
    "github.com/home-assistant/hassio-cli/command/helpers"
    "os"
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
    SnapName := c.String("snapname")
    if c.NArg() > 0 {
        action = c.Args()[0]
    }

    switch action {
    case "list":       // GET
        get = true
    case "info":
        if SnapName == "" {
            fmt.Fprintf(os.Stderr, "-snapname is required. See '%s --help'.", c.App.Name)
            os.Exit(11)
        }
        get = true
        endpoint = SnapName + "/info"
    case "reload":     // POST
        endpoint = action
    case "new":
        endpoint = "new/full"
        if SnapName != "" {
            Options = "name=" + SnapName
        }
    case "restore":
        if SnapName == "" {
            fmt.Fprintf(os.Stderr, "-snapname is required. See '%s --help'.", c.App.Name)
            os.Exit(11)
        }
        endpoint = SnapName + "/restore/full"
    case "remove":
        if SnapName == "" {
            fmt.Fprintf(os.Stderr, "-snapname is required. See '%s --help'.", c.App.Name)
            os.Exit(11)
        }
        endpoint = SnapName + "/remove"
    default:
        fmt.Fprintf(os.Stderr, "No valid action detected")
        os.Exit(3)
    }

    if endpoint != "" || action == "list" {
        helpers.ExecCommand(HassioBasePath, endpoint, serverOverride, get,  Options, Filter, RawJSON)
    }
}
