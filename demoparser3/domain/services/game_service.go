package services

import (
	"github.com/csconfederation/demoparser3/domain"
	"github.com/csconfederation/demoparser3/domain/events"
	"github.com/csconfederation/demoparser3/domain/models"
)

type GameService struct {
	game *models.Game
	bus  *domain.EventBus
}

func NewGameService(bus *domain.EventBus) *GameService {
	return &GameService{
		game: models.NewGame(),
		bus:  bus,
	}
}

func (gs *GameService) Handle(event events.Event) error {
	switch e := event.(type) {
	case events.GameStart:
		return gs.OnGameStart(e)
	case events.MatchStart:
		return gs.OnMatchStart()
	case events.FinalizeRound:
		return gs.FinalizeRound(e)
	}
	return nil
}

func (gs *GameService) OnGameStart(e events.GameStart) error {
	gs.game.MapName = e.MapName
	return nil
}

func (gs *GameService) OnMatchStart() error {
	gs.game.IsGameLive = true
	return nil
}

func (gs *GameService) FinalizeRound(e events.FinalizeRound) error {
	// if the prev round has the same round number as the current finished one, the prev round was reset.
	// might be able to optimize it by storing only previous round instead of all so far

	// if prev was reset, just overwrite in store
	// if round num differs, commit previous one, store current temporarily.
	publishPending := false
	if len(gs.game.Rounds) != 0 {
		prevRound := gs.game.Rounds[len(gs.game.Rounds)-1]
		if prevRound.RoundNum == e.Round.RoundNum {
			gs.game.Rounds = gs.game.Rounds[:len(gs.game.Rounds)-1]
		} else {
			publishPending = true
		}
	}
	gs.game.Rounds = append(gs.game.Rounds, e.Round)

	err := gs.bus.Publish(events.PublishPendingStats{
		PublishPending: publishPending,
	})

	if err != nil {
		return err
	}
	return nil
}
