# Hass.io CLI

[![Build Status](https://travis-ci.org/home-assistant/hassio-cli.svg?branch=master)](https://travis-ci.org/home-assistant/hassio-cli)

## Description

Commandline interface to facilitate interaction with Hass.io server

## Usage

- `hassio help`
- `hassio <subcommand> <action> [<options>]`

E.g.:

- `hassio homeassistant info --raw-json`

### Modifiers

#### Global

- --log-level debug -> will set the log level to debug
- --api-token string   Hass.io API token
- --config string      config file (default is $HOME/.homeassistant.yaml)
- --endpoint string    Endpoint for Hass.io Supervisor ( default is 'hassio' )
- --log-level string   Log level, defaults to Warn
- --raw-json           Output raw JSON from the API

all options are also available as `HASSIO_` prefixed environment variables like `HASSIO_LOG_LEVEL`

#### SubCommand

Available Commands:

- addons
- completion    Generates bash completion scripts
- hardware
- hassos
- homeassistant
- host
- info
- snapshots
- supervisor

## Install

To install, use `go get`:

`go get -d github.com/home-assistant/hassio-cli`

If running on hassio host just run `hassio-cli`, but if on a remote host you'll need to specify token and endpoint:

```
hassio-cli --endpoint $HASS_SERVER/api/hassio --api-token $HASS_TOKEN <cmd>
```

or if you prefer to use environment variables to avoid repetition:

```
export HASSIO_ENDPOINT=https://hassio.local:8123/api/hassio
export HASSIO_API_TOKEN=longandsafesecret
hassio-cli 
```

## Contribution

1. Fork ([https://github.com/home-assistant/hassio-cli/fork](https://github.com/home-assistant/hassio-cli/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
1. Create a new Pull Request

## Building

```bash
go test ./...
gox -osarch="linux/arm" -ldflags="-s -w" -output="hassio"
```
