package stats

import (
	"github.com/csconfederation/demoparser3/domain"
	"github.com/csconfederation/demoparser3/domain/events"
	"github.com/csconfederation/demoparser3/domain/models"
	"github.com/csconfederation/demoparser3/domain/models/stats"
)

type TeamStatsService struct {
	CTSide   string
	TSide    string
	CTStats  *stats.TeamStats
	TStats   *stats.TeamStats
	roundCtx *models.RoundContext
	bus      *domain.EventBus
}

func (tss *TeamStatsService) Handle(event events.Event) error {
	switch e := event.(type) {
	case events.MatchStart:
		return tss.OnMatchStart(e)
	case events.NewRound:
		return tss.OnNewRound(e)
	case events.SetPistolRound:
		return tss.SetPistolRound()
	case events.GameHalfEnded:
		return tss.OnGameHalfEnd()
	case events.RoundEndOfficial:
		return tss.OnRoundEndOfficial()
	case events.RoundEnd:
		return tss.OnRoundEnd(e)
	//case events.Kill:
	//	return tss.OnKill(e)
	case events.Death:
		return tss.OnDeath(e)
	case events.UpdateRoundContext:
		return tss.OnUpdateRoundContext(e)
	}
	return nil
}

func (tss *TeamStatsService) OnMatchStart(e events.MatchStart) error {
	tss.CTSide = e.CTSide
	tss.TSide = e.TSide
	return nil
}

func (tss *TeamStatsService) OnNewRound(e events.NewRound) error {
	tss.roundCtx = e.RoundContext
	tss.CTStats = stats.NewTeamStats(tss.CTSide, e.ConnectedCTPlayers)
	tss.TStats = stats.NewTeamStats(tss.TSide, e.ConnectedTPlayers)
	tss.CTStats.CTRounds.Played = 1
	tss.TStats.TRounds.Played = 1
	return nil
}

func (tss *TeamStatsService) SetPistolRound() error {
	tss.CTStats.PistolRounds.Played = 1
	tss.TStats.PistolRounds.Played = 1
	return nil
}

func (tss *TeamStatsService) OnGameHalfEnd() error {
	tss.CTStats, tss.TStats = tss.TStats, tss.CTStats
	tss.CTSide, tss.TSide = tss.TSide, tss.CTSide
	return nil
}

func (tss *TeamStatsService) OnRoundEndOfficial() error {
	statsByTeam := make(map[string]*stats.TeamStats)
	statsByTeam[tss.CTSide] = tss.CTStats
	statsByTeam[tss.TSide] = tss.TStats
	err := tss.bus.Publish(events.RoundEndTeamStats{
		StatsByTeam: statsByTeam,
	})

	if err != nil {
		return err
	}
	return nil
}

func (tss *TeamStatsService) OnRoundEnd(e events.RoundEnd) error {
	if e.WinningTeamName == tss.CTSide {
		if tss.CTStats.ClutchAttempt != (stats.ClutchAttempt{}) {
			tss.CTStats.ClutchAttempt.IsSuccessful = true
		}
		if tss.TStats.ClutchAttempt != (stats.ClutchAttempt{}) {
			tss.TStats.ClutchAttempt.IsSuccessful = false
		}
	} else if e.WinningTeamName == tss.TSide {
		if tss.TStats.ClutchAttempt != (stats.ClutchAttempt{}) {
			tss.TStats.ClutchAttempt.IsSuccessful = true
		}
		if tss.CTStats.ClutchAttempt != (stats.ClutchAttempt{}) {
			tss.CTStats.ClutchAttempt.IsSuccessful = false
		}
	}
	return nil
}

//func (tss *TeamStatsService) OnKill(e events.Kill) error {
//	baseKillValue := tss.getBaseKillValue(e.KillerTeamName)
//	multiplier := tss.getKillMultiplier()
//}

func (tss *TeamStatsService) getBaseKillValue(killerTeam string) float64 {
	if !tss.roundCtx.IsPostPlant && !tss.roundCtx.IsPostRoundEnd {
		if killerTeam == tss.TSide {
			return 1.2
		}
		return 1
	}

	if tss.roundCtx.IsPostPlant && !tss.roundCtx.IsPostRoundEnd {
		if killerTeam == tss.CTSide {
			return 1.2
		}
		return 1
	}

	if tss.roundCtx.Winner == tss.TSide {
		if killerTeam == tss.TSide {
			return 0.8
		}
		return 0.6
	}

	if killerTeam == tss.TSide {
		return 0.5
	}

	if tss.roundCtx.IsPostPlant {
		return 0.6
	}

	return 0.8
}

func (tss *TeamStatsService) getKillMultiplier() {

}

func (tss *TeamStatsService) OnDeath(e events.Death) error {
	isOpeningKill := tss.handleOpeningKill()
	if e.VictimTeamName == tss.CTSide {
		tss.handleCTSideDeath(e.VictimID, isOpeningKill)
		if len(tss.CTStats.MembersAlive) == 1 && !tss.roundCtx.IsPostRoundEnd {
			tss.CTStats.ClutchAttempt = *stats.NewClutchAttempt(tss.CTStats.MembersAlive[0], len(tss.TStats.MembersAlive))
		}
	} else {
		tss.handleTSideDeath(e.VictimID, isOpeningKill)
		if len(tss.TStats.MembersAlive) == 1 && !tss.roundCtx.IsPostRoundEnd {
			tss.TStats.ClutchAttempt = *stats.NewClutchAttempt(tss.TStats.MembersAlive[0], len(tss.CTStats.MembersAlive))
		}
	}

	return nil
}

func (tss *TeamStatsService) handleOpeningKill() bool {
	if !tss.isOpeningKill() {
		return false
	}

	return true
}

func (tss *TeamStatsService) isOpeningKill() bool {
	countAtRoundStart := len(tss.CTStats.TeamMembers) + len(tss.TStats.TeamMembers)
	countNow := len(tss.CTStats.MembersAlive) + len(tss.TStats.MembersAlive)
	return countAtRoundStart == countNow
}

func (tss *TeamStatsService) handleCTSideDeath(victimID uint64, isOpeningKill bool) {
	tss.CTStats.MembersAlive = removePlayerID(tss.CTStats.MembersAlive, victimID)
	tss.CTStats.DeathPlacement[victimID] = float64(len(tss.CTStats.TeamMembers) - len(tss.CTStats.MembersAlive))
	if isOpeningKill {
		tss.CTStats.FourVFive.IsSet = true
		tss.TStats.FiveVFour.IsSet = true
	}
}

func (tss *TeamStatsService) handleTSideDeath(victimID uint64, isOpeningKill bool) {
	tss.TStats.MembersAlive = removePlayerID(tss.TStats.MembersAlive, victimID)
	tss.TStats.DeathPlacement[victimID] = float64(len(tss.TStats.TeamMembers) - len(tss.TStats.MembersAlive))
	if isOpeningKill {
		tss.TStats.FourVFive.IsSet = true
		tss.CTStats.FiveVFour.IsSet = true
	}
}

func removePlayerID(players []uint64, playerID uint64) []uint64 {
	idx := -1
	for i, item := range players {
		if item == playerID {
			idx = i
			break
		}
	}
	return append(players[:idx], players[idx+1:]...)
}

func (tss *TeamStatsService) OnUpdateRoundContext(e events.UpdateRoundContext) error {
	tss.roundCtx = e.RoundContext
	return nil
}
