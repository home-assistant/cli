package command

import (
    "github.com/urfave/cli"
    "fmt"
    "github.com/home-assistant/hassio-cli/command/helpers"
    "os"
    "strings"
)

func CmdNetwork(c *cli.Context) {
    const HASSIO_BASE_PATH = "network"
    action := ""
    endpoint := ""
    serverOverride := ""
    get := false
    if c.NArg() > 0 {
        action = c.Args()[0]
    }

    switch action {
    case "info":        // GET
        endpoint = action
        get = true
    case "options":     // POST
        serverOverride = "http://172.17.0.2"
        endpoint = action
    default:
        fmt.Fprintf(os.Stderr, "No valid action detected")
        os.Exit(3)
    }

    if endpoint != "" {
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