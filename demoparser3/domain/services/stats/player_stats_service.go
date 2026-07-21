package stats

import (
	"github.com/csconfederation/demoparser3/domain"
	"github.com/csconfederation/demoparser3/domain/events"
	"github.com/csconfederation/demoparser3/domain/models"
	"github.com/csconfederation/demoparser3/domain/models/stats"
)

type PlayerStatsService struct {
	players              map[uint64]string
	pendingStatsByPlayer map[uint64]*stats.PlayerStats
	statsByPlayer        map[uint64]*stats.PlayerStats
	roundCtx             *models.RoundContext
	cTSide               string
	tSide                string
	lastLurkCheckTick    int
	isOpeningKill        bool
	bus                  *domain.EventBus
}

func NewPlayerStatsService(bus *domain.EventBus) *PlayerStatsService {
	return &PlayerStatsService{
		statsByPlayer: make(map[uint64]*stats.PlayerStats),
		bus:           bus,
	}
}

func (pss *PlayerStatsService) Handle(event events.Event) error {
	switch e := event.(type) {
	case events.MatchStart:
		return pss.OnMatchStart(e)
	case events.NewRound:
		return pss.OnNewRound(e)
	case events.UtilityThrown:
		return pss.OnUtilityThrown(e)
	case events.BombExplode:
		return pss.OnBombExplode(e)
	case events.BombDefused:
		return pss.OnBombDefused(e)
	case events.PublishPendingStats:
		return pss.OnPublishPendingStats(e)
	case events.PlayerHurt:
		return pss.OnPlayerHurt(e)
	case events.PlayerFlashed:
		return pss.OnPlayerFlashed(e)
	case events.Death:
		return pss.OnDeath(e)
	case events.Kill:
		return pss.OnKill(e)
	case events.Assist:
		return pss.OnAssist(e)
	case events.UpdateRoundContext:
		return pss.OnUpdateRoundContext(e)
		//case events.GameHalfEnded:
		//	return pss.OnGameHalfEnded()
		//case events.CalculatePlayerRatings:
		//	return pss.OnCalculatePlayerRatings(e)
	}
	return nil
}

func (pss *PlayerStatsService) OnMatchStart(e events.MatchStart) error {

	for steamID, name := range e.CTMembers {
		pss.players[steamID] = name
	}

	for steamID, name := range e.TMembers {
		pss.players[steamID] = name
	}

	pss.cTSide = e.CTSide
	pss.tSide = e.TSide

	return nil
}

func (pss *PlayerStatsService) OnNewRound(e events.NewRound) error {
	pss.roundCtx = e.RoundContext
	for steamID, name := range pss.players {
		pss.statsByPlayer[steamID] = stats.NewPlayerStats(name)
	}
	pss.isOpeningKill = true
	return nil
}

func (pss *PlayerStatsService) OnUtilityThrown(e events.UtilityThrown) error {
	pss.statsByPlayer[e.ThrowerID].UtilityThrown[e.UtilityType] += 1
	return nil
}

func (pss *PlayerStatsService) OnBombExplode(e events.BombExplode) error {
	pss.statsByPlayer[e.PlanterID].ImpactPoints += 0.5
	return nil
}

func (pss *PlayerStatsService) OnBombDefused(e events.BombDefused) error {
	pss.statsByPlayer[e.DefuserID].ImpactPoints += 0.5
	return nil
}

func (pss *PlayerStatsService) OnPublishPendingStats(e events.PublishPendingStats) error {
	if e.PublishPending {
		err := pss.bus.Publish(events.RoundEndPlayerStats{
			StatsByPlayer: pss.pendingStatsByPlayer,
		})

		if err != nil {
			return err
		}
	} else {
		pss.pendingStatsByPlayer = pss.statsByPlayer
	}
	return nil
}

func (pss *PlayerStatsService) OnPlayerHurt(e events.PlayerHurt) error {
	pss.statsByPlayer[e.VictimID].DamageTaken += e.HealthDamage
	pss.statsByPlayer[e.AttackerID].DamageDone += e.HealthDamage
	if e.IsUtility {
		pss.statsByPlayer[e.AttackerID].UtilityDamage += e.HealthDamage
	}
	if e.IsFireDamage {
		pss.statsByPlayer[e.AttackerID].FireDamage += e.HealthDamage
	}
	if e.IsHEDamage {
		pss.statsByPlayer[e.AttackerID].HEDamage += e.HealthDamage
	}
	return nil
}

func (pss *PlayerStatsService) OnPlayerFlashed(e events.PlayerFlashed) error {
	pss.statsByPlayer[e.AttackerID].EnemiesFlashed += 1
	pss.statsByPlayer[e.AttackerID].EnemyFlashDuration += e.FlashDuration
	blindTicks := e.FlashDuration * models.TickRate
	if float64(e.Tick)+blindTicks > pss.statsByPlayer[e.VictimID].MostRecentFlashTickValue {
		pss.statsByPlayer[e.VictimID].MostRecentFlashTickValue = float64(e.Tick) + blindTicks
		pss.statsByPlayer[e.AttackerID].MostRecentFlasherID = e.AttackerID
	}
	return nil
}

func (pss *PlayerStatsService) OnDeath(e events.Death) error {
	pss.statsByPlayer[e.VictimID].Deaths += 1
	pss.statsByPlayer[e.VictimID].TicksAlive = e.Tick - pss.roundCtx.StartingTick
	pss.addSupportDamage(e.VictimID, e.KillerID)
	pss.checkForTraded(e.VictimID, e.Tick)
	return nil
}

func (pss *PlayerStatsService) OnKill(e events.Kill) error {
	pss.statsByPlayer[e.KillerID].KillsList[e.VictimID] = e.Tick
	pss.statsByPlayer[e.KillerID].KASTRounds = 1
	pss.statsByPlayer[e.KillerID].RoundsWithKills = 1

	if pss.isOpeningKill {
		pss.statsByPlayer[e.KillerID].OpeningKills = 1
		pss.statsByPlayer[e.KillerID].Entries = 1
		pss.statsByPlayer[e.VictimID].OpeningDeaths = 1
		pss.isOpeningKill = false
	}

	if e.IsAWPKill {
		pss.statsByPlayer[e.KillerID].AWPKills += 1
		err := pss.bus.Publish(events.CTAWPKill{
			KillerID: e.KillerID,
		})
		if err != nil {
			return err
		}
	}
	if e.IsHeadshot {
		pss.statsByPlayer[e.KillerID].Headshots += 1
	}
	pss.checkForTrades(e.VictimID, e.Tick, e.KillerID)

	pss.statsByPlayer[e.KillerID].KillPoints += pss.calculateKillValue(e.KillerTeamName, e.IsAssisted, e.KillerID,
		e.FlashAssisted, e.KillerEquipmentValue, e.VictimEquipmentValue)
	return nil
}

func (pss *PlayerStatsService) addSupportDamage(victimID uint64, killerID uint64) {
	for steam, damage := range pss.statsByPlayer[victimID].DamageList {
		if killerID != 0 && steam != killerID {
			pss.statsByPlayer[steam].SupportDamage += damage
			if pss.statsByPlayer[steam].SupportDamage > 60 {
				pss.statsByPlayer[steam].SupportRounds = 1
			}
		} else if killerID == 0 {
			pss.statsByPlayer[steam].SupportDamage += damage
			if pss.statsByPlayer[steam].SupportDamage > 60 {
				pss.statsByPlayer[steam].SupportRounds = 1
			}
		}

	}
}

func (pss *PlayerStatsService) checkForTraded(victimID uint64, deathTick int) {
	for killed, tick := range pss.statsByPlayer[victimID].KillsList {
		if deathTick-tick < 4*models.TickRate {
			pss.statsByPlayer[killed].Traded = 1
			pss.statsByPlayer[killed].EffectiveAssistsContribution += 1
			pss.statsByPlayer[killed].KASTRounds = 1
		}
	}
}

func (pss *PlayerStatsService) checkForTrades(victimID uint64, deathTick int, killerID uint64) {
	for _, tick := range pss.statsByPlayer[victimID].KillsList {
		if deathTick-tick < 4*models.TickRate {
			pss.statsByPlayer[killerID].Trades += 1
			return
		}
	}
}

func (pss *PlayerStatsService) OnAssist(e events.Assist) error {
	pss.statsByPlayer[e.AssisterID].Assists += 1
	pss.statsByPlayer[e.AssisterID].EffectiveAssistsContribution += 1
	pss.statsByPlayer[e.AssisterID].KASTRounds = 1
	pss.statsByPlayer[e.AssisterID].SupportRounds = 1
	pss.statsByPlayer[e.AssisterID].ImpactPoints += 0.15
	if e.IsAssistedFlash {
		pss.statsByPlayer[e.AssisterID].FlashAssists += 1
	} else if float64(e.Tick) < pss.statsByPlayer[e.VictimID].MostRecentFlashTickValue {
		pss.flashAssist(e.VictimID)
	}
	return nil
}

func (pss *PlayerStatsService) flashAssist(victimID uint64) {
	flasherID := pss.statsByPlayer[victimID].MostRecentFlasherID
	pss.statsByPlayer[flasherID].FlashAssists += 1
	pss.statsByPlayer[flasherID].EffectiveAssistsContribution += 1
	pss.statsByPlayer[flasherID].SupportRounds = 1
}

func (pss *PlayerStatsService) calculateKillValue(killerTeam string, isAssisted bool, killerID uint64,
	flashAssisted bool, killerEquipmentValue float64, victimEquipmentValue float64) float64 {
	baseValue := pss.getBaseKillValue(killerTeam, isAssisted)
	multiplier := pss.getKillMultiplier(killerTeam, killerID, flashAssisted)
	ecoModifier := pss.getEcoModifier(killerEquipmentValue, victimEquipmentValue)
	return baseValue * multiplier * ecoModifier
}

func (pss *PlayerStatsService) getBaseKillValue(killerTeam string, isAssisted bool) float64 {
	killValue := 1.0

	if pss.roundCtx.IsPrePlant {
		if killerTeam == pss.tSide {
			// Taking site by T
			killValue += 0.2
		}
	} else if pss.roundCtx.IsPostPlant && !pss.roundCtx.IsPostRoundEnd {
		if killerTeam == pss.cTSide {
			// Site retake by CT
			killValue += 0.2
		}
	} else if pss.roundCtx.Winner == pss.tSide {
		if killerTeam == pss.tSide {
			// Chase by T
			killValue -= 0.2
		} else {
			// Exit by CT
			killValue -= 0.4
		}
	} else {
		// CT won
		if killerTeam == pss.tSide {
			// T kill in lost round
			killValue -= 0.5
		} else {
			// CT kill in won round
			if pss.roundCtx.IsPostPlant {
				// Ts got money for planting
				killValue -= 0.4
			} else {
				killValue -= 0.2
			}
		}
	}

	if isAssisted {
		killValue -= 0.15
	}

	return killValue
}

func (pss *PlayerStatsService) getKillMultiplier(killerTeam string, killerID uint64, flashAssisted bool) float64 {
	multiplier := 1.0

	if pss.isOpeningKill {
		if killerTeam == pss.tSide {
			// T entry/opener
			if pss.roundCtx.IsPrePlant {
				multiplier += 0.8
			} else {
				multiplier += 0.3
			}
		} else {
			// CT opener
			multiplier += 0.5
		}
	} else if pss.statsByPlayer[killerID].Trades > 0 {
		multiplier += 0.3
	}

	if flashAssisted {
		multiplier += 0.2
	}

	return multiplier
}

func (pss *PlayerStatsService) getEcoModifier(killerEquipmentValue float64, victimEquipmentValue float64) float64 {
	ratio := victimEquipmentValue / killerEquipmentValue
	modifier := 1.0

	if ratio > 4 {
		modifier += 0.25
	} else if ratio > 2 {
		modifier += 0.14
	} else if ratio < 0.25 {
		modifier -= 0.25
	} else if ratio < 0.5 {
		modifier -= 0.14
	}

	return modifier
}

func (pss *PlayerStatsService) OnUpdateRoundContext(e events.UpdateRoundContext) error {
	pss.roundCtx = e.RoundContext
	return nil
}

func (pss *PlayerStatsService) CheckLurk(e events.CheckLurk) error {
	// check every 4 seconds after round start
	if pss.lastLurkCheckTick+4*models.TickRate < e.Tick {
		return nil
	}

	lurkerID := getLurker(e.Distances)
	if lurkerID != 0 {
		pss.statsByPlayer[lurkerID].LurkerBlips++
	}
	return nil
}

func getLurker(distances map[uint64]map[uint64]float64) uint64 {
	maxDistance := 0.0
	var lurkerID uint64

	for playerID, teammates := range distances {
		totalDist := 0.0
		for _, dist := range teammates {
			if dist >= 500 {
				totalDist += dist
			}
		}

		if totalDist > maxDistance {
			maxDistance = totalDist
			lurkerID = playerID
		}
	}

	return lurkerID
}

//func (pss *PlayerStatsService) OnCalculatePlayerRatings(e events.CalculatePlayerRatings) error {
//	for steamID, stats := range pss.statsByPlayer {
//		//openingFactor := (float64(stats.OpeningKills-stats.OpeningDeaths.Ol) / 13.0) + 1 //move from 13 to (rounds / 5)
//		//playerIPR := stats.ImpactPoints / float64(stats.Rounds)
//		//playerWPR := stats.WinPoints / float64(stats.Rounds)
//
//		//if stats.TeamsWinPoints != 0 {
//		//	player.ImpactRating = (0.1 * float64(openingFactor)) + (0.6 * (playerIPR / impactRoundAvg)) + (0.3 * (playerWPR / (player.TeamsWinPoints / float64(player.WinPointsNormalizer))))
//		//} else {
//		//	player.ImpactRating = (0.1 * float64(openingFactor)) + (0.6 * (playerIPR / impactRoundAvg))
//		//}
//		stats.ImpactRating = 0.0
//		playerDR := float64(stats.Deaths) / float64(stats.Rounds)
//		playerRatingDeathComponent := 0.07 * (e.OverallAverages.DeathRoundAverage / playerDR)
//		if stats.Deaths == 0 || playerRatingDeathComponent > 0.21 {
//			playerRatingDeathComponent = 0.21
//		}
//		// KillRatio is in derived
//		// KAST in derived
//		// ADR in derived
//		stats.Rating = (0.3 * stats.ImpactRating) + (0.35 * (stats.Kills / e.OverallAverages.KillRoundAverage)) +
//			playerRatingDeathComponent + (0.08 * (stats.KASTRounds / e.OverallAverages.KASTRoundAverage)) +
//			(0.2 * (stats.ADR / e.OverallAverages.ADRAverage))
//
//	}
//}
