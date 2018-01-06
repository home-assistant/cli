package command

import (
    "fmt"

    "github.com/home-assistant/hassio-cli/command/helpers"
    "github.com/urfave/cli"
    "os"
    "strings"
)

func CmdHomeassistant(c *cli.Context) {
    const HassioBasePath = "homeassistant"
    action := ""
    endpoint := ""
    serverOverride := ""
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
    default:
        fmt.Fprintf(os.Stderr, "No valid action detected")
        os.Exit(3)
    }

    if endpoint != "" {
        uri := helpers.GenerateUri(HassioBasePath, endpoint, serverOverride)
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
