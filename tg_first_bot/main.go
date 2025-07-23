package main

import (
	"example.com/clients"
	"example.com/consumer"
	"example.com/events"
	"example.com/storage"
	"log"
)

const (
	tgBotHost   = "api.telegram.org"
	basePath    = "storage"
	startOffset = 0
)

func main() {
	tgClient := clients.NewClient(tgBotHost)
	internalBasePath := storage.NewInternalBasePath(basePath)
	newConsumer := consumer.NewConsumer(
		events.NewEventProcess(tgClient, startOffset),
		internalBasePath,
	)
	if err := newConsumer.Start(); err != nil {
		log.Printf("Failed start()", err)
		panic("failed Start()")
	}
}
