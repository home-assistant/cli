package main

import (
	"log/slog"
	"os"

	"github.com/home-assistant/cli/cmd"
)

func main() {
	// Only log the warning severity or above.
	slog.SetLogLoggerLevel(slog.LevelWarn)
	defer func() {
		if cmd.ExitWithError {
			os.Exit(1)
		}
	}()
	cmd.Execute()
}
