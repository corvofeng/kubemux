package main

import (
	log "github.com/sirupsen/logrus"

	"gmux/cmd/gmux/internal/command"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
		DisableQuote:  true,
	})

	log.SetLevel(log.DebugLevel)
	logger := log.New()

	if err := command.Root(logger).Execute(); err != nil {
		log.Fatal(err)
	}
}
