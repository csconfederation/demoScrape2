package domain

import (
	"reflect"

	"github.com/csconfederation/demoparser3/domain/events"
	"github.com/csconfederation/demoparser3/logger"
)

type EventHandler interface {
	Handle(event events.Event) error
}

type EventBus struct {
	handlers map[string][]EventHandler
}

func NewEventBus() *EventBus {
	return &EventBus{
		handlers: make(map[string][]EventHandler),
	}
}

func (bus *EventBus) Subscribe(eventType string, handler EventHandler) {
	bus.handlers[eventType] = append(bus.handlers[eventType], handler)
}

func (bus *EventBus) Publish(event events.Event) error {
	eventType := reflect.TypeOf(event).Name()
	handlers := bus.handlers[eventType]

	for _, handler := range handlers {
		if err := handler.Handle(event); err != nil {
			logger.Error("handler failed",
				"event", eventType,
				"error", err.Error(),
			)
			return err
		}
	}

	return nil
}
