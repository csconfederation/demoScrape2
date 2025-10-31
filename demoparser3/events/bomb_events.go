package events

import (
	"github.com/csconfederation/demoparser3/logger"
	"github.com/csconfederation/demoparser3/types"
	"github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs"
	"github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs/events"
)

func RegisterBombEvents(parser demoinfocs.Parser, game *types.Game) {
	if parser == nil {
		logger.Error("Parser cannot be nil")

	}
	if game == nil {
		logger.Error("Game cannot be nil")
		return
	}

	parser.RegisterEventHandler(func(event events.BombPlanted) {
		if !parser.GameState().IsMatchStarted() {
			return
		}
		if game.CurrentRound == nil {
			logger.Error("Game has no current round", "event", "BombPlanted")
			return
		}
		if game.CurrentRound.PostWinCon {
			return
		}
		game.CurrentRound.BombPlanted(event.Player)
	})

	parser.RegisterEventHandler(func(event events.BombDefused) {
		if !parser.GameState().IsMatchStarted() {
			return
		}

		if game.CurrentRound == nil {
			logger.Error("Game has no current round", "event", "BombDefused")
			return
		}

		if game.CurrentRound.PostWinCon {
			return
		}

		game.CurrentRound.BombDefused(event.Player)
	})

	parser.RegisterEventHandler(func(event events.BombExplode) {

		if !parser.GameState().IsMatchStarted() {
			return
		}

		if game.CurrentRound == nil {
			logger.Error("Game has no current round", "event", "BombExplode")
			return
		}

		if game.CurrentRound.PostWinCon {
			return
		}

		game.CurrentRound.BombExplode()
	})
}
