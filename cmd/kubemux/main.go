package main

import (
	log "github.com/sirupsen/logrus"

	"kubemux/cmd/kubemux/internal/command"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
		DisableQuote:  true,
	})

	logger := log.New()

	if err := command.Root(logger).Execute(); err != nil {
		log.Fatal(err)
	}
}
