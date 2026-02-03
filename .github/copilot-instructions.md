# Home Assistant CLI

## Project Overview
This is the official Home Assistant CLI tool written in Go, providing command-line
interface to interact with the Home Assistant Supervisor. The CLI enables users to
manage apps, control the core system, handle audio/network settings, manage
backups, and perform various system operations.

The CLI is communicating with the Supervisor using the Supervisor's HTTP REST API.

## Repository Structure
- **`main.go`** - Entry point of the application
- **`cmd/`** - Contains all CLI command implementations using Cobra framework
- **`client/`** - HTTP client functionality for API communication
- **`spinner/`** - Progress spinner implementation
- **Root files** - Configuration and documentation

## Key Technologies
- **Language**: Go (use modern syntax)
- **CLI Framework**: Cobra (github.com/spf13/cobra)
- **HTTP Client**: Resty (github.com/go-resty/resty/v2)
- **Configuration**: Viper (github.com/spf13/viper)
- **Logging**: Go stdlib log/slog

## Available Commands
The CLI provides the following main command categories:
- `apps` - Install, update, remove and configure Home Assistant apps
- `audio` - Audio plug-in management
- `authentication` - Authentication for Home Assistant users
- `cli` - CLI plug-in management
- `core` - Home Assistant Core control
- `dns` - DNS plug-in management
- `docker` - Docker related configuration
- `hardware` - System hardware information
- `host` - Host OS control
- `info` - General Home Assistant information
- `multicast` - Multicast plug-in configuration
- `network` - Network configuration and management
- `observer` - Observer plug-in management
- `os` - Home Assistant OS specific operations
- `resolution` - Resolution center for issues and solutions
- `backups` - Backup creation, restoration, and management
- `supervisor` - Supervisor monitoring and control

## Development Environment
- **API Endpoint**: Configurable via `SUPERVISOR_ENDPOINT` environment variable
- **Authentication**: Uses API tokens via `SUPERVISOR_API_TOKEN` environment variable
- **Config File**: Optional config file support (default: `$HOME/.homeassistant.yaml`)

## Build Commands
- **Build**: `CGO_ENABLED=0 go build -ldflags="-s -w" -o "ha"`
- **Test**: `go test ./...`
- **Format**: `gofmt -s`

## File Organization Patterns
- Command files follow naming pattern: `<component>_<action>.go`
- Each command typically has its own file in the `cmd/` directory
- Helper functions are in `client/helper.go`
- Main client logic is in `client/client.go`

## Architecture Notes
- Uses Cobra for command structure and flag parsing
- Resty for HTTP API calls to Home Assistant Supervisor
- Viper for configuration management
- Go stdlib log/slog for structured logging
- Custom spinner implementation for progress indication

## Testing
- Unit tests available in `client/helper_test.go`
- Test command: `go test ./...`

## Contributing Guidelines
1. Create feature branch
2. Commit changes
3. Rebase against master
4. Run tests with `go test ./...`
5. Format code with `gofmt -s`
6. Create Pull Request

This CLI is designed to work with Home Assistant Supervisor API and is commonly used in
Home Assistant Operating System environments, SSH apps, and development setups.
