package main

import (
	"os"

	"github.com/urfave/cli"
)

func main() {

	app := cli.NewApp()
	app.Name = Name
	app.Version = Version
	app.Author = "Ryan McLean"
	app.Email = "pvtryan100@gmail.com"
	app.Usage = "Commandline tool to allow interation with hass.io"

	app.Flags = GlobalFlags
	app.Commands = Commands
	app.CommandNotFound = CommandNotFound

	app.Run(os.Args)
}
