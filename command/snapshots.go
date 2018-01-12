package command

import (
    "github.com/urfave/cli"
    "fmt"
    "github.com/home-assistant/hassio-cli/command/helpers"
    "os"
)

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
    if c.NArg() > 0 {
        action = c.Args()[0]
    }

    switch action {
    case "list":       // GET
        get = true
    case "info":
        get = true
        endpoint = c.String("snapname") + "/info"
    case "reload":     // POST
        endpoint = action
    case "new":
        endpoint = "new/full"
        if c.String("snapname") != "" {
            Options = "name=" + c.String("snapname")
        }
    case "restore":
        if c.String("snapname") == "" {
            fmt.Fprintf(os.Stderr, "-snapname is required. See '%s --help'.", c.App.Name)
            os.Exit(11)
        }
        endpoint = c.String("snapname") + "/restore/full"
    case "remove":
        if c.String("snapname") == "" {
            fmt.Fprintf(os.Stderr, "-snapname is required. See '%s --help'.", c.App.Name)
            os.Exit(11)
        }
        endpoint = c.String("snapname") + "/remove"
    default:
        fmt.Fprintf(os.Stderr, "No valid action detected")
        os.Exit(3)
    }

    if endpoint != "" {
        helpers.ExecCommand(HassioBasePath, endpoint, serverOverride, get,  Options, Filter, RawJSON)
    }
}
