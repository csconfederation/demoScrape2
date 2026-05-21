package stats

import (
	"github.com/csconfederation/demoparser3/domain"
	"github.com/csconfederation/demoparser3/domain/events"
	"github.com/csconfederation/demoparser3/domain/models"
	"github.com/csconfederation/demoparser3/domain/models/stats"
)

type SidedPlayerStatsService struct {
	players                   map[uint64]string
	pendingSidedStatsByPlayer map[uint64]map[models.Side]*stats.SidedPlayerStats
	sidedStatsByPlayer        map[uint64]map[models.Side]*stats.SidedPlayerStats
	sideByPlayer              map[uint64]models.Side
	isOpeningKill             bool
	bus                       *domain.EventBus
}

func (spss *SidedPlayerStatsService) Handle(event events.Event) error {
	switch e := event.(type) {
	case events.MatchStart:
		return spss.OnMatchStart(e)
	case events.NewRound:
		return spss.OnNewRound(e)
	//case events.UtilityThrown:
	//	return pss.OnUtilityThrown(e)
	case events.BombExplode:
		return spss.OnBombExplode(e)
	case events.BombDefused:
		return spss.OnBombDefused(e)
	case events.PublishPendingStats:
		return spss.OnPublishPendingStats(e)
	case events.PlayerHurt:
		return spss.OnPlayerHurt(e)
	//case events.PlayerFlashed:
	//	return pss.OnPlayerFlashed(e)
	case events.Death:
		return spss.OnDeath(e)
	case events.Kill:
		return spss.OnKill(e)
	case events.Assist:
		return spss.OnAssist(e)
	//case events.UpdateRoundContext:
	//	return pss.OnUpdateRoundContext(e)
	case events.GameHalfEnded:
		return spss.OnGameHalfEnded()
		//case events.CalculatePlayerRatings:
		//	return pss.OnCalculatePlayerRatings(e)
	}
	return nil
}

func (spss *SidedPlayerStatsService) OnMatchStart(e events.MatchStart) error {
	for steamID, name := range e.CTMembers {
		spss.players[steamID] = name
		spss.sidedStatsByPlayer[steamID] = make(map[models.Side]*stats.SidedPlayerStats)
		spss.sideByPlayer[steamID] = models.CounterTerrorists
	}

	for steamID, name := range e.TMembers {
		spss.players[steamID] = name
		spss.sidedStatsByPlayer[steamID] = make(map[models.Side]*stats.SidedPlayerStats)
		spss.sideByPlayer[steamID] = models.Terrorists
	}

	return nil
}

func (spss *SidedPlayerStatsService) OnNewRound(e events.NewRound) error {
	for steamID, name := range spss.players {
		spss.sidedStatsByPlayer[steamID] = map[models.Side]*stats.SidedPlayerStats{
			models.CounterTerrorists: stats.NewSidedPlayerStats(name),
			models.Terrorists:        stats.NewSidedPlayerStats(name),
		}
	}
	spss.isOpeningKill = true
	return nil
}

func (spss *SidedPlayerStatsService) OnBombExplode(e events.BombExplode) error {
	spss.sidedStatsByPlayer[e.PlanterID][models.Terrorists].ImpactPoints += 0.5
	return nil
}

func (spss *SidedPlayerStatsService) OnBombDefused(e events.BombDefused) error {
	spss.sidedStatsByPlayer[e.DefuserID][models.CounterTerrorists].ImpactPoints += 0.5
	return nil
}

func (spss *SidedPlayerStatsService) OnPublishPendingStats(e events.PublishPendingStats) error {
	if e.PublishPending {
		err := spss.bus.Publish(events.RoundEndSidedPlayerStats{
			SidedStatsByPlayer: spss.pendingSidedStatsByPlayer,
		})

		if err != nil {
			return err
		}
	} else {
		spss.pendingSidedStatsByPlayer = spss.sidedStatsByPlayer
	}
	return nil
}

func (spss *SidedPlayerStatsService) OnPlayerHurt(e events.PlayerHurt) error {
	attackerSide := spss.sideByPlayer[e.AttackerID]
	spss.sidedStatsByPlayer[e.AttackerID][attackerSide].DamageDone += e.HealthDamage
	return nil
}

func (spss *SidedPlayerStatsService) OnDeath(e events.Death) error {
	victimSide := spss.sideByPlayer[e.VictimID]
	spss.sidedStatsByPlayer[e.VictimID][victimSide].Deaths += 1
	return nil
}

func (spss *SidedPlayerStatsService) OnKill(e events.Kill) error {
	killerSide := spss.sideByPlayer[e.KillerID]
	victimSide := spss.sideByPlayer[e.VictimID]

	spss.sidedStatsByPlayer[e.KillerID][killerSide].KASTRounds = 1

	if spss.isOpeningKill {
		spss.sidedStatsByPlayer[e.KillerID][killerSide].OpeningKills = 1
		spss.sidedStatsByPlayer[e.VictimID][victimSide].OpeningDeaths = 1
		spss.isOpeningKill = false
	}

	if e.IsAWPKill {
		spss.sidedStatsByPlayer[e.KillerID][killerSide].AWPKills += 1
	}
	return nil
}

func (spss *SidedPlayerStatsService) OnAssist(e events.Assist) error {
	assisterSide := spss.sideByPlayer[e.AssisterID]

	spss.sidedStatsByPlayer[e.AssisterID][assisterSide].KASTRounds = 1
	spss.sidedStatsByPlayer[e.AssisterID][assisterSide].ImpactPoints += 0.15

	return nil
}

func (spss *SidedPlayerStatsService) OnGameHalfEnded() error {
	for steamID, side := range spss.sideByPlayer {
		if side == models.CounterTerrorists {
			spss.sideByPlayer[steamID] = models.Terrorists
		} else {
			spss.sideByPlayer[steamID] = models.CounterTerrorists
		}
	}
	return nil
}
