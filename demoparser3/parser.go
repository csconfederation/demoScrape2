package main

import (
	"io"

	"github.com/csconfederation/demoparser3/adapter"
	"github.com/csconfederation/demoparser3/domain"
	"github.com/csconfederation/demoparser3/domain/services"
	"github.com/csconfederation/demoparser3/logger"
)

func ProcessDemo(demo io.ReadCloser) (error, error) {

	bus := domain.NewEventBus()
	cscParser := adapter.NewAdapter(demo, bus)
	gs := services.NewGameService(bus)

	bus.Subscribe("GameStart", gs)
	bus.Subscribe("MatchStart", gs)
	bus.Subscribe("RoundStart", gs)

	err := cscParser.Parse()

	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	return nil, nil
}
