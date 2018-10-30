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
		switch logformat := c.GlobalString("log-format"); logformat {
		case "json":
			log.SetFormatter(&log.JSONFormatter{})
		case "text":
		default:
			log.SetFormatter(&log.TextFormatter{})
		}
		return nil
	}
	app.CommandNotFound = CommandNotFound

	err := app.Run(os.Args)
	if err != nil {
		log.WithField("error", err).Error("some error occured")
	}
}
