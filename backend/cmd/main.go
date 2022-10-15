package main

import (
	"log"

	"github.com/jwnpoh/njcreaderapp/backend/cmd/broker"
	"github.com/jwnpoh/njcreaderapp/backend/cmd/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("unable to load config, failed to start app")
	}

	broker := broker.NewBrokerService(cfg)

	log.Panic(broker.Start())
}
