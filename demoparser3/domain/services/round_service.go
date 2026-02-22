package services

import (
	"github.com/csconfederation/demoparser3/domain"
	"github.com/csconfederation/demoparser3/domain/events"
	"github.com/csconfederation/demoparser3/domain/models"
)

type RoundService struct {
	currentRound *models.Round
	bus          *domain.EventBus
}

func (rs *RoundService) Handle(event events.Event) error {
	switch e := event.(type) {
	case events.RoundStart:
		return rs.OnRoundStart(e)
	case events.SetPistolRound:
		return rs.SetPistolRound()
	case events.BombPlanted:
		return rs.OnBombPlant()
	case events.RoundEnd:
		return rs.OnRoundEnd()
	case events.RoundEndOfficial:
		return rs.OnRoundEndOfficial()
	case events.FrameDone:
		return rs.OnFrameDone(e)
	}
	return nil
}

func (rs *RoundService) OnRoundStart(e events.RoundStart) error {
	rs.currentRound = models.NewRound(e.TotalRoundsPlayed, e.Tick)
	err := rs.bus.Publish(events.NewRound{
		RoundContext:       rs.currentRound.Context(),
		ConnectedCTPlayers: e.ConnectedCTPlayers,
		ConnectedTPlayers:  e.ConnectedTPlayers,
	})
	if err != nil {
		return err
	}
	return nil
}

func (rs *RoundService) SetPistolRound() error {
	rs.currentRound.IsPistolRound = true
	return nil
}

func (rs *RoundService) OnBombPlant() error {
	rs.currentRound.IsPrePlant = false
	rs.currentRound.IsPostPlant = true
	err := rs.bus.Publish(events.UpdateRoundContext{
		RoundContext: rs.currentRound.Context(),
	})
	if err != nil {
		return err
	}
	return nil
}

func (rs *RoundService) OnRoundEnd() error {
	rs.currentRound.IsPostRoundEnd = true
	err := rs.bus.Publish(events.UpdateRoundContext{
		RoundContext: rs.currentRound.Context(),
	})
	if err != nil {
		return err
	}
	return nil
}

func (rs *RoundService) OnRoundEndOfficial() error {
	err := rs.bus.Publish(events.FinalizeRound{
		Round: rs.currentRound,
	})
	rs.currentRound = nil
	if err != nil {
		return err
	}
	return nil
}

func (rs *RoundService) OnFrameDone(e events.FrameDone) error {
	if rs.currentRound != nil {
		return nil
	}

	// in round
	if rs.currentRound.StartingTick+10*models.TickRate < e.Tick {
		// 10 seconds haven't passed
		return nil
	}

	err := rs.bus.Publish(events.CheckLurk{
		Distances: e.Distances,
	})

	if err != nil {
		return err
	}

	return nil
}
