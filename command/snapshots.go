package command

import (
    "fmt"
    "os"

    "github.com/home-assistant/hassio-cli/command/helpers"
    "github.com/urfave/cli"
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
    SnapSlug := c.String("slug")
    if c.NArg() > 0 {
        action = c.Args()[0]
    }

    switch action {
        case "list": // GET
            get = true
        case "info":
            if SnapSlug == "" {
                fmt.Fprintf(os.Stderr, "-snapname is required. See '%s --help'.", c.App.Name)
                os.Exit(11)
            }
            get = true
            endpoint = SnapSlug + "/info"
        case "reload": // POST
            endpoint = action
        case "new":
            endpoint = "new/full"
            if c.String("name") != "" {
                if Options != "" {
                    Options += ","
                }
                Options += "name=" + c.String("name")
            }
            if c.String("password") != "" {
                if Options != "" {
                    Options += ","
                }
                Options += "password=" + c.String("password")
            }
        case "restore":
            if SnapSlug == "" {
                fmt.Fprintf(os.Stderr, "-slug is required. See '%s --help'.", c.App.Name)
                os.Exit(11)
            }
            if c.String("password") != "" {
                Options = "password=" + c.String("password")
            }
            endpoint = SnapSlug + "/restore/full"
        case "remove":
            if SnapSlug == "" {
                fmt.Fprintf(os.Stderr, "-slug is required. See '%s --help'.", c.App.Name)
                os.Exit(11)
            }
            endpoint = SnapSlug + "/remove"
        default:
            fmt.Fprintf(os.Stderr, "No valid action detected")
            os.Exit(3)
    }

    if endpoint != "" || action == "list" {
        helpers.ExecCommand(HassioBasePath, endpoint, serverOverride, get, Options, Filter, RawJSON)
    }
}
