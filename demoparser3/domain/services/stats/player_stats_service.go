package stats

import (
	"github.com/csconfederation/demoparser3/domain"
	"github.com/csconfederation/demoparser3/domain/events"
	"github.com/csconfederation/demoparser3/domain/models"
	"github.com/csconfederation/demoparser3/domain/models/stats"
)

type PlayerStatsService struct {
	players           map[uint64]string
	statsByPlayer     map[uint64]*stats.PlayerStats
	roundCtx          *models.RoundContext
	cTSide            string
	tSide             string
	lastLurkCheckTick int
	bus               *domain.EventBus
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
	case events.RoundEndOfficial:
		return pss.OnRoundEndOfficial()
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
	case events.GameHalfEnded:
		return pss.OnGameHalfEnded()
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

func (pss *PlayerStatsService) OnRoundEndOfficial() error {
	err := pss.bus.Publish(events.RoundEndPlayerStats{
		StatsByPlayer: pss.statsByPlayer,
	})
	if err != nil {
		return err
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
	pss.statsByPlayer[e.KillerID].KASTRound = 1
	pss.statsByPlayer[e.KillerID].RoundsWithKills = 1
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
	// TODO: Idk how to handle this
	//pss.statsByPlayer[e.KillerID].ImpactPoints += pss.calculateKillValue(e.KillerTeamName)
	return nil
}

func (pss *PlayerStatsService) addSupportDamage(victimID uint64, killerID uint64) {
	for steam, damage := range pss.statsByPlayer[victimID].DamageList {
		if killerID != 0 && steam != killerID {
			pss.statsByPlayer[steam].SupportDamage += damage
			if pss.statsByPlayer[steam].SupportDamage > 60 {
				pss.statsByPlayer[steam].SupportRound = 1
			}
		} else if killerID == 0 {
			pss.statsByPlayer[steam].SupportDamage += damage
			if pss.statsByPlayer[steam].SupportDamage > 60 {
				pss.statsByPlayer[steam].SupportRound = 1
			}
		}

	}
}

func (pss *PlayerStatsService) checkForTraded(victimID uint64, deathTick int) {
	for killed, tick := range pss.statsByPlayer[victimID].KillsList {
		if deathTick-tick < 4*models.TickRate {
			pss.statsByPlayer[killed].Traded = 1
			pss.statsByPlayer[killed].EAC += 1
			pss.statsByPlayer[killed].KASTRound = 1
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
	pss.statsByPlayer[e.AssisterID].EAC += 1
	pss.statsByPlayer[e.AssisterID].KASTRound = 1
	pss.statsByPlayer[e.AssisterID].SupportRound = 1
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
	pss.statsByPlayer[flasherID].EAC += 1
	pss.statsByPlayer[flasherID].SupportRound = 1
}

func (pss *PlayerStatsService) calculateKillValue(killerTeam string, isOpeningKill bool) float64 {
	//baseValue := pss.getBaseKillValue(killerTeam)
	//multiplier := pss.getKillMultiplier(isOpeningKill)
	//ecoModifier := pss.getEcoModifier()
	//return baseValue * multiplier * ecoModifier
	return 0
}

//func (pss *PlayerStatsService) getBaseKillValue(killerTeam string) float64 {
//	if pss.roundCtx.IsPrePlant {
//		if killerTeam == pss.tSide {
//			// Taking site by T
//			return 1.2
//		}
//
//		// Site defense by CT
//		return 1.0
//	}
//
//	if pss.roundCtx.IsPostPlant {
//		if killerTeam == pss.tSide {
//			// Site defense by T
//			return 1.0
//		}
//
//		// Retake by CT
//		return 1.2
//	}
//
//	if pss.roundCtx.Winner == pss.tSide {
//		if killerTeam == pss.tSide {
//			// Chase by T
//			return 0.8
//		}
//
//		// Exit by CT
//		return 0.6
//	}
//
//	// CT won
//	if killerTeam == pss.tSide {
//		// T kill in lost round
//		return 0.5
//	}
//
//	// CT kill in won round
//	if pss.roundCtx.IsPostPlant {
//		// Ts got money for planting
//		return 0.6
//	}
//	return 0.8
//}
//
//func (pss *PlayerStatsService) getKillMultiplier(isOpeningKill bool) float64 {
//	return nil
//}

func (pss *PlayerStatsService) OnUpdateRoundContext(e events.UpdateRoundContext) error {
	pss.roundCtx = e.RoundContext
	return nil
}

func (pss *PlayerStatsService) OnGameHalfEnded() error {
	pss.cTSide, pss.tSide = pss.tSide, pss.cTSide
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
