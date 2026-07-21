package events

import (
	"github.com/csconfederation/demoparser3/domain/models"
	"github.com/csconfederation/demoparser3/domain/models/stats"
)

type NewRound struct {
	RoundContext       *models.RoundContext
	ConnectedCTPlayers []uint64
	ConnectedTPlayers  []uint64
}

type SetPistolRound struct{}

type RoundReset struct{}

type RoundEndPlayerStats struct {
	StatsByPlayer map[uint64]*stats.PlayerStats
}

type RoundEndTeamStats struct {
	StatsByTeam map[string]*stats.TeamStats
}

type RoundEndSidedPlayerStats struct {
	SidedStatsByPlayer map[uint64]map[models.Side]*stats.SidedPlayerStats
}

type CTAWPKill struct {
	KillerID uint64
}

type UpdateRoundContext struct {
	RoundContext *models.RoundContext
}

type FinalizeRound struct {
	Round *models.Round
}

type CheckLurk struct {
	Tick      int
	Distances map[uint64]map[uint64]float64
}

type PublishPendingStats struct {
	PublishPending bool
}

type CalculateDerivedStats struct {
	PlayerSteamID uint64
	PlayerStats   *stats.PlayerStats
}

type CalculateSidedDerivedStats struct {
	PlayerSteamID    uint64
	PlayerSide       models.Side
	SidedPlayerStats *stats.SidedPlayerStats
}

type UpdateOverallAverage struct {
	PlayerStats *stats.PlayerStats
}

type UpdateSidedAverage struct {
	PlayerSide       models.Side
	SidedPlayerStats *stats.SidedPlayerStats
}

type GameEnd struct{}

type CalculatePlayerRatings struct {
	OverallAverages *stats.Averages
	CTSideAverages  *stats.Averages
	TSideAverages   *stats.Averages
}
