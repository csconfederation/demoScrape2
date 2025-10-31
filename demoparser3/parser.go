package main

import (
	"io"

	"github.com/csconfederation/demoparser3/events"
	"github.com/csconfederation/demoparser3/logger"
	"github.com/csconfederation/demoparser3/types"
	"github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs"
	"go.uber.org/zap"
)

func ProcessDemo(demo io.ReadCloser) (*types.Game, error) {

	game := types.NewGame()
	parser := demoinfocs.NewParser(demo)
	defer func() {
		if err := parser.Close(); err != nil {
			logger.Error("Failed to close parser", zap.Error(err))
		}
	}()

	events.RegisterRoundEvents(parser, game)
	events.RegisterBombEvents(parser, game)
	events.RegisterPlayerEvents(parser, game)
	events.RegisterGameEvents(parser, game)

	err := parser.ParseToEnd()

	//TODO: stats processing
	//err := game.EndOfMatchProcessing()
	//
	if err != nil {
		return nil, err
	}

	return game, nil
}
