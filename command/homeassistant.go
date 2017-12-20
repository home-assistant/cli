package command

import (
	"fmt"

	"github.com/home-assistant/hassio-cli/command/helpers"
	"github.com/urfave/cli"
)

const HASSIO_BASE_PATH string = "homeassistant"

func CmdHomeassistant(c *cli.Context) {
	x := helpers.RestCall(HASSIO_BASE_PATH, "info", "")
	fmt.Println(x)
}
