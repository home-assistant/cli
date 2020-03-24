# Home Assistant CLI

[![Build Status](https://travis-ci.org/home-assistant/cli.svg?branch=master)](https://travis-ci.org/home-assistant/cli)

## Description

Commandline interface to facilitate interaction with Home Assistant

## Usage

- `ha help`
- `ha <subcommand> <action> [<options>]`

E.g.:

- `ha core info --raw-json`

### Modifiers

#### Global

- --log-level debug -> will set the log level to debug
- --api-token string   Hass.io API token
- --config string      config file (default is $HOME/.homeassistant.yaml)
- --endpoint string    Endpoint for the Home Assistant Supervisor (default is 'supervisor')
- --log-level string   Log level, defaults to Warn
- --raw-json           Output raw JSON from the API

all options are also available as `SUPERVISOR_` prefixed environment variables like `SUPERVISOR_LOG_LEVEL`

#### SubCommand

Available Commands:

- addons         Install, update, remove and configure Home Assistant add-ons
- audio          Audio device handling.
- authentication Authentication for Home Assistant users.
- core           Provides control of the Home Assistant Core
- dns            Get information, update or configure the Home Assistant DNS server
- hardware       Provides hardware information about your system
- help           Help about any command
- host           Control the host/system that Home Assistant is running on
- info           Provides a general Home Assistant information overview
- os             Operating System specific for updating, info and configuration imports
- snapshots      Create, restore and remove snapshot backups
- supervisor     Monitor, control and configure the Home Assistant Supervisor

## Install

Download the pre-build executable from the releases page:

<https://github.com/home-assistant/cli/releases>

If running on the Home Assistant host just run `ha`, but if on a remote host you'll need to specify token and endpoint:

```shell
ha --endpoint $HA_SERVER/api/hassio --api-token $SUPERVISOR_TOKEN <cmd>
```

or if you prefer to use environment variables to avoid repetition:

```shell
export SUPERVISOR_ENDPOINT=https://homeassistant.local:8123/api/hassio
export SUPERVISOR_API_TOKEN=longandsafesecret
ha
```

## Contribution

1. Fork ([https://github.com/home-assistant/cli/fork](https://github.com/home-assistant/cli/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the `master` branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
1. Create a new Pull Request

## Building

```bash
go test ./...
gox -osarch="linux/arm" -ldflags="-s -w" -output="ha"
```
