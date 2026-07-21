package stats

import (
	"github.com/csconfederation/demoparser3/domain"
	"github.com/csconfederation/demoparser3/domain/events"
	"github.com/csconfederation/demoparser3/domain/models/stats"
)

type AveragesService struct {
	Overall *stats.Averages
	CTSide  *stats.Averages
	TSide   *stats.Averages
	bus     *domain.EventBus
}

func (as *AveragesService) Handle(event events.Event) error {
	switch e := event.(type) {
	case events.UpdateOverallAverage:
		return as.OnUpdateOverallAverage(e)
	case events.GameEnd:
		return as.OnGameEnd()
	}
	return nil
}

func (as *AveragesService) OnUpdateOverallAverage(e events.UpdateOverallAverage) error {
	rounds := float64(e.PlayerStats.Rounds)
	impactPoints := e.PlayerStats.ImpactPoints
	kills := float64(e.PlayerStats.Kills)
	deaths := float64(e.PlayerStats.Deaths)
	kastRounds := float64(e.PlayerStats.KASTRounds)
	damage := float64(e.PlayerStats.DamageDone)

	newAverages := stats.NewAverages(rounds, impactPoints, kills, deaths, kastRounds, damage)
	as.Overall.Aggregate(newAverages)
	return nil
}

func (as *AveragesService) OnUpdateSidedAverage(e events.UpdateOverallAverage) error {
	rounds := float64(e.PlayerStats.Rounds)
	impactPoints := e.PlayerStats.ImpactPoints
	kills := float64(e.PlayerStats.Kills)
	deaths := float64(e.PlayerStats.Deaths)
	kastRounds := float64(e.PlayerStats.KASTRounds)
	damage := float64(e.PlayerStats.DamageDone)

	newAverages := stats.NewAverages(rounds, impactPoints, kills, deaths, kastRounds, damage)
	as.Overall.Aggregate(newAverages)
	return nil
}

func (as *AveragesService) OnGameEnd() error {
	as.Overall.Normalize()
	as.CTSide.Normalize()
	as.TSide.Normalize()

	//err := as.bus.Publish(events.CalculatePlayerRatings{
	//	OverallAverages: as.Overall,
	//	CTSideAverages:  as.CTSide,
	//	TSideAverages:   as.TSide,
	//})
	//if err != nil {
	//	return err
	//}

	return nil
}
