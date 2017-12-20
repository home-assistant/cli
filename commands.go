package main

import (
	"fmt"
	"os"

	"github.com/home-assistant/hassio-cli/command"
	"github.com/urfave/cli"
)

var GlobalFlags = []cli.Flag{}

var Commands = []cli.Command{
	{
		Name:   "homeassistant",
		Usage:  "",
		Action: command.CmdHomeassistant,
		Flags:  []cli.Flag{},
	},
	{
		Name:   "supervisor",
		Usage:  "",
		Action: command.CmdSupervisor,
		Flags:  []cli.Flag{},
	},
	{
		Name:   "host",
		Usage:  "",
		Action: command.CmdHost,
		Flags:  []cli.Flag{},
	},
}

func CommandNotFound(c *cli.Context, command string) {
	fmt.Fprintf(os.Stderr, "%s: '%s' is not a %s command. See '%s --help'.", c.App.Name, command, c.App.Name, c.App.Name)
	os.Exit(2)
}
