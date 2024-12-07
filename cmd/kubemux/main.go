package main

import (
	log "github.com/sirupsen/logrus"

	"kubemux/cmd/kubemux/internal/command"
)

func main() {

	if err := command.Root().Execute(); err != nil {
		log.Fatal(err)
	}
}
