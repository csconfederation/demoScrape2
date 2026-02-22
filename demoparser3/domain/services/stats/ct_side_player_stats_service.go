package stats

import (
	"github.com/csconfederation/demoparser3/domain"
	"github.com/csconfederation/demoparser3/domain/events"
	"github.com/csconfederation/demoparser3/domain/models/stats"
)

type CTSidePlayerStatsService struct {
	ctSideStatsByPlayer map[uint64]*stats.CTSidePlayerStats
	bus                 *domain.EventBus
}

func (ct *CTSidePlayerStatsService) Handle(event events.Event) error {
	switch e := event.(type) {
	case events.CTAWPKill:
		return ct.OnAWPKill(e)
	}
	return nil
}

func (ct *CTSidePlayerStatsService) OnAWPKill(e events.CTAWPKill) error {
	ct.ctSideStatsByPlayer[e.KillerID].CTAWPKills += 1
	return nil
}
