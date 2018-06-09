package command

import (
    "fmt"
    "github.com/home-assistant/hassio-cli/command/helpers"
    "github.com/urfave/cli"
    "os"
)


// CmdAddons All addon endpoints for hass.io
func CmdAddons(c *cli.Context) {
    const HassioBasePath = "addons"
    action := ""
    endpoint := ""
    serverOverride := ""
    get := false
    DebugEnabled := c.GlobalBool("debug")
    helpers.DebugEnabled = DebugEnabled
    Options := c.String("options")
    RawJSON := c.Bool("rawjson")
    Filter := c.String("filter")
    AddonName := c.String("name")
    if c.NArg() > 0 {
        action = c.Args()[0]
    }

    switch action {
    case "list":      // GET
        endpoint = ""
        get = true
    case "info",
         "logo",
         "changelog",
         "logs",  // Fix as not JSON format for output
         "stats":
        if AddonName == "" {
            fmt.Fprintf(os.Stderr, "-name is required. See '%s --help'.\n", c.App.Name)
            os.Exit(11)
        }
        endpoint = AddonName + "/" + action
        get = true
    case "reload":     // POST
        endpoint = action
    case "start",
         "stop",
         "install",
         "uninstall",
         "update":
        if AddonName == "" {
            fmt.Fprintf(os.Stderr, "-name is required. See '%s --help'.\n", c.App.Name)
            os.Exit(11)
        }
        endpoint = AddonName + "/" + action
    default:
        fmt.Fprintf(os.Stdout, "No valid action detected.\n")
        os.Exit(3)
    }

    if DebugEnabled {
        fmt.Fprintf(os.Stdout, "DEBUG [addons]: action->'%s', endpoint='%s', serverOverride->'%s', GET->'%t', options->'%s', rawjson->'%t', filter->'%s'\n",
            action, endpoint, serverOverride, get, Options, RawJSON, Filter )
    }
    if endpoint != "" || action == "list" {
        helpers.ExecCommand(HassioBasePath, endpoint, serverOverride, get,  Options, Filter, RawJSON)
    }
}
