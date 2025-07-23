package consumer

import (
	"example.com/events"
	"example.com/storage"
)

type Consumer struct {
	client           *events.EventProcessor
	internalBasePath *storage.InternalBasePath
}

func NewConsumer(client *events.EventProcessor, internalBasePath *storage.InternalBasePath) *Consumer {
	return &Consumer{
		client:           client,
		internalBasePath: internalBasePath,
	}
}
