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
    get := false
    if c.NArg() > 0 {
        action = c.Args()[0]
    }

    switch action {
    case "list":       // GET
        get = true
        action = ""
    case "reload":     // POST
        endpoint = action
    case "new":
        endpoint = action
    default:
        fmt.Fprintf(os.Stderr, "No valid action detected")
        os.Exit(3)
    }

    if endpoint != "" || action == "list" {
        uri := helpers.GenerateUri(HASSIO_BASE_PATH, endpoint, serverOverride)
        response := helpers.RestCall(uri, get,  c.String("options"))

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
