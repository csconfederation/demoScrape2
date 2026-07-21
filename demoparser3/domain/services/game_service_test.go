package services

import (
	"testing"

	"github.com/csconfederation/demoparser3/domain"
	"github.com/csconfederation/demoparser3/domain/events"
	"github.com/stretchr/testify/assert"
)

func TestNewGameService(t *testing.T) {
	bus := domain.NewEventBus()
	gs := NewGameService(bus)

	assert.NotNil(t, gs.game)
	assert.Equal(t, gs.bus, bus)
}

func TestGameService_Handle(t *testing.T) {
	bus := domain.NewEventBus()
	gs := NewGameService(bus)
	event := events.GameStart{
		MapName: "Test",
	}

	err := gs.Handle(event)

	assert.NoError(t, err)
	assert.Equal(t, gs.game.MapName, event.MapName)
}

func TestGameService_OnGameStart(t *testing.T) {
	bus := domain.NewEventBus()
	gs := NewGameService(bus)
	event := events.GameStart{
		MapName: "Test",
	}

	err := gs.OnGameStart(event)

	assert.NoError(t, err)
	assert.Equal(t, gs.game.MapName, event.MapName)
}
