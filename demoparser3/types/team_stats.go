package types

import (
	"fmt"

	"github.com/csconfederation/demoparser3/types/stats"
)

type totalWon struct {
	Total int
	Won   int
}

type conditionOutcome struct {
	IsSet bool
	Won   bool
}

type TeamStats struct {
	MembersAliveAtRoundStart int
	MembersAlive             int              `json:"membersAlive"`
	Saves                    int              `json:"saves" end_of_match_sum:"true"`
	Clutches                 int              `json:"clutches" end_of_match_sum:"true"`
	FiveVFour                conditionOutcome `json:"5v4" end_of_match_sum:"true"`
	FourVFive                conditionOutcome `json:"4v5" end_of_match_sum:"true"`
	PistolRounds             totalWon         `json:"pistolRounds" end_of_match_sum:"true"`
	CTRounds                 totalWon         `json:"ctRounds" end_of_match_sum:"true"`
	TRounds                  totalWon         `json:"tRounds" end_of_match_sum:"true"`
}

func NewTeamStats(connectedTeamPlayers int) *TeamStats {
	return &TeamStats{
		MembersAliveAtRoundStart: connectedTeamPlayers,
		MembersAlive:             connectedTeamPlayers,
	}
}

func NewTotalTeamStats(teams map[string]*Team) map[string]*TeamStats {
	totalTeamStats := make(map[string]*TeamStats)
	for _, team := range teams {
		totalTeamStats[team.Name] = &TeamStats{}
	}

	return totalTeamStats
}

func (teamStats *TeamStats) PlayerKilled() int {
	teamStats.MembersAlive -= 1
	//if teamStats.MembersAlive == 4 && killerTeamAliveMembers == 5 {
	//	teamStats.FourVFive.IsSet = true
	//}

	return teamStats.MembersAliveAtRoundStart - teamStats.MembersAlive
}

func (teamStats *TeamStats) Kill(victimTeamAliveMembers int) {
	if teamStats.MembersAlive == 5 && victimTeamAliveMembers == 4 {
		teamStats.FiveVFour.IsSet = true
		return
	}
}

func (teamStats *TeamStats) Aggregate(new stats.Stats) error {
	newStats, ok := new.(*TeamStats)
	if !ok {
		return fmt.Errorf("expected *TeamStats, got %T", new)
	}
	return stats.Sum(teamStats, newStats)
}
