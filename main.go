package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {

	// Only log the warning severity or above.
	log.SetLevel(log.WarnLevel)

	app := cli.NewApp()
	app.Name = Name
	app.Version = Version
	app.Author = "Home-Assistant"
	app.Email = "hello@home-assistant.io"
	app.Usage = "Commandline tool to allow interaction with hass.io"

	app.Flags = GlobalFlags
	app.Commands = Commands
	app.Before = func(c *cli.Context) error {
		if c.GlobalBool("debug") {
			log.SetLevel(log.DebugLevel)
		}
		return nil
	}
	app.CommandNotFound = CommandNotFound

	app.Run(os.Args)
}
