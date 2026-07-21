package services

import (
	"github.com/csconfederation/demoparser3/domain"
	"github.com/csconfederation/demoparser3/domain/events"
	"github.com/csconfederation/demoparser3/domain/models"
)

type TeamService struct {
	CTSide *models.Team
	TSide  *models.Team
	bus    *domain.EventBus
}

func (ts *TeamService) Handle(event events.Event) error {
	switch e := event.(type) {
	case events.MatchStart:
		return ts.OnMatchStart(e)
	case events.GameHalfEnded:
		return ts.OnGameHalfEnd()
	case events.ScoreUpdated:
		return ts.OnScoreUpdate(e)
	case events.RoundStart:
		return ts.OnRoundStart(e)
	}
	return nil
}

func (ts *TeamService) OnMatchStart(e events.MatchStart) error {
	ts.CTSide = models.NewTeam(e.CTSide)
	ts.TSide = models.NewTeam(e.TSide)
	return nil
}

func (ts *TeamService) OnGameHalfEnd() error {
	ts.CTSide, ts.TSide = ts.TSide, ts.CTSide
	return nil
}

func (ts *TeamService) OnRoundStart(e events.RoundStart) error {
	if (e.TotalRoundsPlayed%12)+1 == 1 && e.TotalRoundsPlayed <= models.MR {
		err := ts.bus.Publish(events.SetPistolRound{})
		if err != nil {
			return err
		}
	}
	return nil
}

func (ts *TeamService) OnScoreUpdate(e events.ScoreUpdated) error {
	if e.OldScore > e.NewScore {
		// round was reset
		err := ts.bus.Publish(events.RoundReset{})
		if err != nil {
			return err
		}
	}
	switch e.TeamName {
	case ts.CTSide.Name:
		ts.CTSide.Score = e.NewScore
	case ts.TSide.Name:
		ts.TSide.Score = e.NewScore
	}
	return nil
}
