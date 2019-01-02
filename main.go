package main

import (
	"github.com/home-assistant/hassio-cli/cmd"
	log "github.com/sirupsen/logrus"
)

func main() {
	// Only log the warning severity or above.
	log.SetLevel(log.WarnLevel)

	cmd.Execute()
}
