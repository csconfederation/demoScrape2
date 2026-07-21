package stats

import (
	"github.com/csconfederation/demoparser3/domain/events"
	"github.com/csconfederation/demoparser3/domain/models"
	"github.com/csconfederation/demoparser3/domain/models/stats"
)

type DerivedPlayerStatsService struct {
	DerivedStatsByPlayer      map[uint64]stats.DerivedPlayerStats
	SidedDerivedStatsByPlayer map[uint64]map[models.Side]stats.DerivedPlayerStats
}

func (dpss *DerivedPlayerStatsService) Handle(event events.Event) error {
	switch e := event.(type) {
	case events.CalculateDerivedStats:
		return dpss.OnCalculateDerivedStats(e)
	case events.CalculateSidedDerivedStats:
		return dpss.OnCalculateSidedDerivedStats(e)
	}
	return nil
}

func (dpss *DerivedPlayerStatsService) NewDerivedPlayerStatsService() *DerivedPlayerStatsService {
	return &DerivedPlayerStatsService{
		DerivedStatsByPlayer: make(map[uint64]stats.DerivedPlayerStats),
	}
}

func (dpss *DerivedPlayerStatsService) OnCalculateDerivedStats(e events.CalculateDerivedStats) error {
	derived := dpss.DerivedStatsByPlayer[e.PlayerSteamID]
	playerStats := e.PlayerStats

	derived.AverageTimeAlive = playerStats.TicksAlive / playerStats.Rounds / models.TickRate
	derived.KAST = float64(playerStats.KASTRounds / playerStats.Rounds)
	if playerStats.Kills > 0 {
		derived.KillPointAverage = playerStats.KillPoints / float64(playerStats.Kills)
	}
	if playerStats.ImpactPoints > 0 {
		derived.ImpactInWinningRounds = playerStats.WinPoints / playerStats.ImpactPoints
	}
	if playerStats.Deaths > 0 {
		derived.AverageDeathPlacement = playerStats.DeathPlacement / float64(playerStats.Deaths)
	} else {
		derived.TradeRatio = 0.5
	}
	derived.ADR = float64(playerStats.DamageDone / playerStats.Rounds)
	derived.DamageRatingDifferential = derived.ADR - float64(playerStats.DamageTaken/playerStats.Rounds)
	derived.TradeRatio = float64(playerStats.Traded / playerStats.Deaths)
	derived.KillRatio = float64(playerStats.Kills / playerStats.Rounds)
	derived.UtilityThrown = playerStats.GetTotalUtilityThrown()
	//playerStats.Rws?
	return nil
}

func (dpss *DerivedPlayerStatsService) OnCalculateSidedDerivedStats(e events.CalculateSidedDerivedStats) error {
	derived := dpss.SidedDerivedStatsByPlayer[e.PlayerSteamID][e.PlayerSide]
	playerStats := e.SidedPlayerStats

	derived.ADR = float64(playerStats.DamageDone / playerStats.Rounds)
	derived.KAST = float64(playerStats.KASTRounds / playerStats.Rounds)
	if playerStats.Deaths > 0 {
		derived.AverageDeathPlacement = playerStats.DeathPlacement / float64(playerStats.Deaths)
	}

	return nil
}
