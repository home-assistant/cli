package command

import (
    "github.com/urfave/cli"
    "fmt"
    "github.com/home-assistant/hassio-cli/command/helpers"
    "os"
    "strings"
)

func CmdSnapshots(c *cli.Context) {
    const HASSIO_BASE_PATH = "snapshots"
    action := ""
    endpoint := ""
    serverOverride := ""
    options := c.String("options")
    get := false
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
            options = "name=" + c.String("snapname")
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

    if endpoint != "" || action == "list" {
        uri := helpers.GenerateUri(HASSIO_BASE_PATH, endpoint, serverOverride)
        response := helpers.RestCall(uri, get,  options)

        if c.String("filter") == "" {
            helpers.DisplayOutput(response, c.Bool("rawjson"))
        } else {
            filter := strings.Split(c.String("filter"), ",")
            data := helpers.FilterProperties(response, filter)
            helpers.DisplayOutput(data, c.Bool("rawjson"))
        }
        responseMap := helpers.ByteArrayToMap(response)
        result := responseMap["result"]
        if result != "ok" {
            os.Exit(10)
        }
    }
}
