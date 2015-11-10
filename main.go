package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/scheedule/backendcontroller/commands"
)

func main() {
	log.SetLevel(log.DebugLevel)
	commands.Execute()
}
