package main

// https://github.com/home-assistant/hassio/blob/dev/API.md
import (
	"fmt"
	"os"

	"github.com/home-assistant/hassio-cli/command"
	"github.com/urfave/cli"
)

// GlobalFlags Used to hold global flags
var GlobalFlags = []cli.Flag{
	cli.BoolFlag{
		Name:  "debug, d",
		Usage: "Prints Debug information",
	},
	cli.StringFlag{
		Name:  "log-format",
		Usage: "log format to use, valid options are text and json. Default is text",
	},
}

// Commands holds the commands that are supported by the CLI
var Commands = []cli.Command{
	{
		Name:    "homeassistant",
		Aliases: []string{"ha"},
		Usage:   "info, logs, check, restart, start, stop, update",
		Action:  command.CmdHomeassistant,
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "rawjson, j",
				Usage: "Returns the output in JSON format",
			},
			cli.StringFlag{
				Name:  "options, o",
				Usage: "holds data for POST in format `key=val,key2=val2`",
			},
			cli.StringFlag{
				Name:  "filter, f",
				Usage: "properties to extract from returned data `prop1,prop2`",
			},
		},
	},
	{
		Name:    "supervisor",
		Usage:   "info, logs, reload, update",
		Aliases: []string{"su"},
		Action:  command.CmdSupervisor,
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "rawjson, j",
				Usage: "Returns the output in JSON format",
			},
			cli.StringFlag{
				Name:  "options, o",
				Usage: "holds data for POST in format `key=val,key2=val2`",
			},
			cli.StringFlag{
				Name:  "filter, f",
				Usage: "properties to extract from returned data `prop1,prop2`",
			},
		},
	},
	{
		Name:    "host",
		Usage:   "reboot, shutdown",
		Aliases: []string{"ho"},
		Action:  command.CmdHost,
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "rawjson, j",
				Usage: "Returns the output in JSON format",
			},
			cli.StringFlag{
				Name:  "options, o",
				Usage: "holds data for POST in format `key=val,key2=val2`",
			},
			cli.StringFlag{
				Name:  "filter, f",
				Usage: "properties to extract from returned data `prop1,prop2`",
			},
		},
	},
	{
		Name:    "hassos",
		Aliases: []string{"os"},
		Usage:   "info, update",
		Action:  command.CmdHassOS,
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "rawjson, j",
				Usage: "Returns the output in JSON format",
			},
			cli.StringFlag{
				Name:  "options, o",
				Usage: "holds data for POST in format `key=val,key2=val2`",
			},
			cli.StringFlag{
				Name:  "filter, f",
				Usage: "properties to extract from returned data `prop1,prop2`",
			},
		},
	},
	{
		Name:    "hardware",
		Usage:   "info, audio",
		Aliases: []string{"hw"},
		Action:  command.CmdHardware,
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "rawjson, j",
				Usage: "Returns the output in JSON format",
			},
			cli.StringFlag{
				Name:  "filter, f",
				Usage: "properties to extract from returned data `prop1,prop2`",
			},
		},
	},
	{
		Name:    "snapshots",
		Usage:   "list, info, reload, new, restore, remove",
		Aliases: []string{"sn"},
		Action:  command.CmdSnapshots,
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "rawjson, j",
				Usage: "Returns the output in JSON format",
			},
			cli.StringFlag{
				Name:  "options, o",
				Usage: "holds data for POST in format `key=val,key2=val2`",
			},
			cli.StringFlag{
				Name:  "filter, f",
				Usage: "properties to extract from returned data `prop1,prop2`",
			},
			cli.StringFlag{
				Name:  "slug",
				Usage: "used with 'info|remove|restore|new' actions to return info on a specific snapshot `slugofsnapshot`",
			},
			cli.StringFlag{
				Name:  "name",
				Usage: "used with 'restore|new' actions to set a name for a snapshot",
			},
			cli.StringFlag{
				Name:  "password",
				Usage: "used with 'restore|new' actions to set a password on a snapshot",
			},
		},
	},
	{
		Name:    "addons",
		Usage:   "list, info, logo, changelog, logs, stats,\n reload, start, stop, install, uninstall, update",
		Aliases: []string{"ad"},
		Action:  command.CmdAddons,
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "rawjson, j",
				Usage: "Returns the output in JSON format",
			},
			cli.StringFlag{
				Name:  "options, o",
				Usage: "holds data for POST in format `key=val,key2=val2`",
			},
			cli.StringFlag{
				Name:  "filter, f",
				Usage: "properties to extract from returned data `prop1,prop2`",
			},
			cli.StringFlag{
				Name:  "name",
				Usage: "used with 'info' actions to return info on a specific addon `nameofaddon`",
			},
		},
	},
}

// CommandNotFound used to display if a user enters a non-existant command
func CommandNotFound(c *cli.Context, command string) {
	fmt.Fprintf(os.Stderr, "%s: '%s' is not a %s command. See '%s --help'.\n", c.App.Name, command, c.App.Name, c.App.Name)
	os.Exit(2)
}
