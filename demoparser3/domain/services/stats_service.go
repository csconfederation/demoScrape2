package services

import (
	"github.com/csconfederation/demoparser3/domain/events"
	"github.com/csconfederation/demoparser3/domain/models/stats"
)

type StatsService struct {
	StatsByPlayer       map[uint64]*stats.PlayerStats
	CTSideStatsByPlayer map[uint64]*stats.CTSidePlayerStats
	TSideStatsByPlayer  map[uint64]*stats.TSidePlayerStats
	StatsByTeam         map[string]*stats.TeamStats
}

func (ss *StatsService) Handle(event events.Event) error {
	switch e := event.(type) {
	case events.GameStart:
		return ss.OnGameStart()
	case events.RoundEndTeamStats:
		return ss.OnRoundEndTeamStats(e)
	case events.RoundEndPlayerStats:
		return ss.OnRoundEndPlayerStats(e)
	}
	return nil
}

func (ss *StatsService) OnGameStart() error {
	ss.StatsByPlayer = make(map[uint64]*stats.PlayerStats)
	ss.CTSideStatsByPlayer = make(map[uint64]*stats.CTSidePlayerStats)
	ss.TSideStatsByPlayer = make(map[uint64]*stats.TSidePlayerStats)
	ss.StatsByTeam = make(map[string]*stats.TeamStats)
	return nil
}

func (ss *StatsService) OnRoundEndTeamStats(e events.RoundEndTeamStats) error {
	for teamName, teamStats := range e.StatsByTeam {
		ss.StatsByTeam[teamName].Aggregate(teamStats)
	}
	return nil
}

func (ss *StatsService) OnRoundEndPlayerStats(e events.RoundEndPlayerStats) error {
	for steamID, playerStats := range e.StatsByPlayer {
		ss.StatsByPlayer[steamID].Aggregate(playerStats)
		if playerStats.PlayerSide == stats.CounterTerrorists {
			ss.CTSideStatsByPlayer[steamID].Aggregate(playerStats)
		} else {
			ss.TSideStatsByPlayer[steamID].Aggregate(playerStats)
		}
	}
	return nil
}
