package stats

import (
	"github.com/csconfederation/demoparser3/domain"
	"github.com/csconfederation/demoparser3/domain/events"
	"github.com/csconfederation/demoparser3/domain/models"
	"github.com/csconfederation/demoparser3/domain/models/stats"
)

type FinalStatsService struct {
	StatsByPlayer      map[uint64]*stats.PlayerStats
	SidedStatsByPlayer map[uint64]map[models.Side]*stats.SidedPlayerStats
	StatsByTeam        map[string]*stats.TeamStats
	bus                *domain.EventBus
}

func (fss *FinalStatsService) Handle(event events.Event) error {
	switch e := event.(type) {
	case events.MatchStart:
		return fss.OnMatchStart(e)
	case events.RoundEndTeamStats:
		return fss.OnRoundEndTeamStats(e)
	case events.RoundEndPlayerStats:
		return fss.OnRoundEndPlayerStats(e)
	case events.RoundEndSidedPlayerStats:
		return fss.OnRoundEndSidedPlayerStats(e)
	}
	return nil
}

func (fss *FinalStatsService) OnMatchStart(e events.MatchStart) error {
	fss.StatsByPlayer = make(map[uint64]*stats.PlayerStats)
	for steamID, name := range e.CTMembers {
		fss.SidedStatsByPlayer[steamID] = map[models.Side]*stats.SidedPlayerStats{
			models.CounterTerrorists: stats.NewSidedPlayerStats(name),
			models.Terrorists:        stats.NewSidedPlayerStats(name),
		}
	}
	for steamID, name := range e.TMembers {
		fss.SidedStatsByPlayer[steamID] = map[models.Side]*stats.SidedPlayerStats{
			models.CounterTerrorists: stats.NewSidedPlayerStats(name),
			models.Terrorists:        stats.NewSidedPlayerStats(name),
		}
	}
	fss.StatsByTeam = make(map[string]*stats.TeamStats)
	return nil
}

func (fss *FinalStatsService) OnRoundEndTeamStats(e events.RoundEndTeamStats) error {
	for teamName, teamStats := range e.StatsByTeam {
		fss.StatsByTeam[teamName].Aggregate(teamStats)
	}
	return nil
}

func (fss *FinalStatsService) OnRoundEndPlayerStats(e events.RoundEndPlayerStats) error {
	for steamID, playerStats := range e.StatsByPlayer {
		fss.StatsByPlayer[steamID].Aggregate(playerStats)
		err := fss.bus.Publish(events.CalculateDerivedStats{
			PlayerSteamID: steamID,
			PlayerStats:   fss.StatsByPlayer[steamID],
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (fss *FinalStatsService) OnRoundEndSidedPlayerStats(e events.RoundEndSidedPlayerStats) error {
	for steamID, sideMap := range e.SidedStatsByPlayer {
		for side, sidedStats := range sideMap {
			fss.SidedStatsByPlayer[steamID][side].Aggregate(sidedStats)
			err := fss.bus.Publish(events.CalculateSidedDerivedStats{
				PlayerSteamID:    steamID,
				PlayerSide:       side,
				SidedPlayerStats: fss.SidedStatsByPlayer[steamID][side],
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}
