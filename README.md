# hassio-cli



## Description

## Usage

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
