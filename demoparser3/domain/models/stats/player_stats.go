package stats

import (
	"golang.org/x/exp/constraints"
)

type PlayerStats struct {
	BasePlayerStats
	KillPoints                   float64 `json:"killPoints" end_of_match_sum:"true"`
	DamageTaken                  int     `json:"damageTaken" end_of_match_sum:"true"`
	UtilityDamage                int     `json:"utilityDamage" end_of_match_sum:"true"`
	HEDamage                     int     `json:"heDamage" end_of_match_sum:"true"`
	FireDamage                   int     `json:"fireDamage" end_of_match_sum:"true"`
	EnemiesFlashed               int     `json:"enemiesFlashed" end_of_match_sum:"true"`
	EnemyFlashDuration           float64 `json:"enemyFlashTime" end_of_match_sum:"true"`
	MostRecentFlasherID          uint64  `json:"mostRecentFlasher"`
	MostRecentFlashTickValue     float64
	UtilityThrown                map[string]int `json:"utilityThrown" end_of_match_sum:"true"`
	Assists                      int            `json:"assists" end_of_match_sum:"true"`
	RoundsFor                    int
	RoundsAgainst                int
	RoundsWithKills              int
	Headshots                    int
	EffectiveAssistsContribution int            `json:"eac" end_of_match_sum:"true"`
	FlashAssists                 int            `json:"flashAssists" end_of_match_sum:"true"`
	DamageList                   map[uint64]int `json:"damageList"`
	SupportDamage                int            `json:"supportDamage" end_of_match_sum:"true"`
	SupportRounds                int
	TicksAlive                   int
	KillsList                    map[uint64]int
	Traded                       int
	Trades                       int
	LurkerBlips                  int `json:"lurkerBlips" end_of_match_sum:"true"`
	Entries                      int
	MostImpactfulPlayer          float64
	Clutches                     map[int]int
	Multikills                   map[int]int
	Saves                        int
	ImpactRating                 float64
	Rating                       float64
}

func NewPlayerStats(name string) *PlayerStats {
	return &PlayerStats{
		BasePlayerStats: *NewBasePlayerStats(name),
	}
}

func (ps *PlayerStats) Aggregate(newStats *PlayerStats) {
	ps.AggregateBasePlayerStats(&newStats.BasePlayerStats)
	ps.Assists += newStats.Assists
	ps.TicksAlive += newStats.TicksAlive
	ps.Trades += newStats.Trades
	ps.Traded += newStats.Traded
	ps.KillPoints += newStats.KillPoints
	addMap(ps.Clutches, newStats.Clutches)
	addMap(ps.Multikills, newStats.Multikills)
	ps.HEDamage += newStats.HEDamage
	ps.FireDamage += newStats.FireDamage
	ps.UtilityDamage += newStats.UtilityDamage
	ps.EnemiesFlashed += newStats.EnemiesFlashed
	ps.FlashAssists += newStats.FlashAssists
	ps.EnemyFlashDuration += newStats.EnemyFlashDuration
	ps.Headshots += newStats.Headshots
	ps.Saves += newStats.Saves
	ps.Entries += newStats.Entries
	ps.RoundsFor += newStats.RoundsFor
	ps.RoundsAgainst += newStats.RoundsAgainst
	addMap(ps.UtilityThrown, newStats.UtilityThrown)
	ps.DamageTaken += newStats.DamageTaken
	ps.SupportDamage += newStats.SupportDamage
	ps.SupportRounds += newStats.SupportRounds
	ps.RoundsWithKills += newStats.RoundsWithKills
	ps.MostImpactfulPlayer += newStats.MostImpactfulPlayer
	ps.EffectiveAssistsContribution += newStats.EffectiveAssistsContribution
	//ps.PlayerSide = models.Unknown

	// if RoundsFor == 1 add WinTeamDamage for that round. Lives in round me thinks
	//if newStats.RoundsFor == 1 {
	//	ps.WinTeamDamage += newStats.DamageDone
	//}
}

func addMap[K comparable, V constraints.Integer](dst, src map[K]V) {
	for k, v := range src {
		dst[k] += v
	}
}

func (ps *PlayerStats) GetTotalUtilityThrown() int {
	total := 0
	for _, v := range ps.UtilityThrown {
		total += v
	}

	return total
}
