package main

import (
	"os"

	"github.com/home-assistant/cli/cmd"
	log "github.com/sirupsen/logrus"
)

func main() {
	// Only log the warning severity or above.
	log.SetLevel(log.WarnLevel)
	defer func() {
		if cmd.ExitWithError {
			os.Exit(1)
		}
	}()
	cmd.Execute()
}
