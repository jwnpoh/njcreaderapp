package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/jwnpoh/njcreaderapp/backend/broker"
)

func main() {
	godotenv.Load(".env")
	port := os.Getenv("PORT")

	broker := broker.NewBrokerService(port)

	log.Panic(broker.Start())
}
