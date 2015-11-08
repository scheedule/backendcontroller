package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/scheedule/backend_controller/commands"
)

func main() {
	log.SetLevel(log.DebugLevel)
	commands.Execute()
}
