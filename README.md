# hassiocli


## Description

Commandline interface to facilitate interaction with hass.io server

## Usage

`hassiocli <subcommand> <action> [<options>]`

E.g.

- `hassiocli homeassistant --json info`

would return the info from hass.io in JSON format.

To send data to an endpoint, in this case to goto a specific version of hass.io:

- `hassiocli homeassistant update --payload="version=0.60"`

## Install

To install, use `go get`:

```bash
$ go get -d github.com/home-assistant/hassio-cli
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
gox -osarch="linux/arm" -ldflags="-s -w" -output="hassiocli"
upx --brute hassiocli
```

## Author

[home-assistant](https://github.com/home-assistant)
