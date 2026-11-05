package client

import "runtime/debug"

// hassSource identifies the CLI in the X-Hass-Source header. It carries no
// version so Supervisor's reports group across CLI versions.
const hassSource = "HomeAssistantCLI"

// Version is the CLI version, set at build time with
// -ldflags "-X github.com/home-assistant/cli/client.Version=<version>".
var Version = ""

// userAgent returns the User-Agent sent with every Supervisor API request.
func userAgent() string {
	return hassSource + "/" + version()
}

// version returns the CLI version, falling back to the module version Go
// records at build time, and to "dev" when that is unknown too.
func version() string {
	if Version != "" {
		return Version
	}
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
		return info.Main.Version
	}
	return "dev"
}
