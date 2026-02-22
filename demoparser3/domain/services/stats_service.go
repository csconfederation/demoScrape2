package services

//
//import (
//	"github.com/csconfederation/demoparser3/domain/events"
//	"github.com/csconfederation/demoparser3/domain/models/stats"
//)
//
//type StatsService struct {
//	currentStats    *stats.Stats
//	aggregatedStats *stats.Stats
//}
//
//func (ss *StatsService) Handle(event events.Event) error {
//	switch e := event.(type) {
//	case events.RoundStart:
//		return ss.OnRoundStart(e)
//	case events.SetPistolRound:
//		return ss.SetPistolRound()
//	}
//	return nil
//}
//
//func (ss *StatsService) SetPistolRound() error {
//	for teamName := range ss.currentStats.StatsByTeam {
//		ss.currentStats.StatsByTeam[teamName].PistolRounds.Played = 1
//	}
//
//	return nil
//}
