package stats

import "golang.org/x/exp/constraints"

type PlayerStats struct {
	Name                         string
	ImpactPoints                 float64 `json:"impactPoints" end_of_match_sum:"true"`
	KillPoints                   float64 `json:"killPoints" end_of_match_sum:"true"`
	DamageDone                   int     `json:"damageDone" end_of_match_sum:"true"`
	DamageTaken                  int     `json:"damageTaken" end_of_match_sum:"true"`
	UtilityDamage                int     `json:"utilityDamage" end_of_match_sum:"true"`
	HEDamage                     int     `json:"heDamage" end_of_match_sum:"true"`
	FireDamage                   int     `json:"fireDamage" end_of_match_sum:"true"`
	EnemiesFlashed               int     `json:"enemiesFlashed" end_of_match_sum:"true"`
	EnemyFlashDuration           float64 `json:"enemyFlashTime" end_of_match_sum:"true"`
	MostRecentFlasherID          uint64  `json:"mostRecentFlasher"`
	MostRecentFlashTickValue     float64
	UtilityThrown                map[string]int `json:"utilityThrown" end_of_match_sum:"true"`
	Kills                        int            `json:"kills" end_of_match_sum:"true"`
	Deaths                       int            `json:"deaths" end_of_match_sum:"true"`
	DeathPlacement               int            `json:"deathOrder" end_of_match_sum:"true"` // Order in which a player died. Ex.: Entry fraggers will have the lowest order (possibly 0)
	Assists                      int            `json:"assists" end_of_match_sum:"true"`
	KASTRounds                   int            `json:"kastRounds" end_of_match_sum:"true"`
	RoundsFor                    int
	RoundsAgainst                int
	RoundsWithKills              int
	AWPKills                     int
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
	OpeningKills                 int
	OpeningDeaths                int
	MostImpactfulPlayer          float64
	Rounds                       int
	Clutches                     map[int]int
	Multikills                   map[int]int
	Saves                        int
	WinPoints                    float64
	PlayerSide                   Side
}

func NewPlayerStats(name string, side Side) *PlayerStats {
	return &PlayerStats{
		Name:       name,
		PlayerSide: side,
	}
}

func (ps *PlayerStats) Aggregate(newStats *PlayerStats) {
	ps.Rounds += 1
	ps.Kills += newStats.Kills
	ps.Deaths += newStats.Deaths
	ps.Assists += newStats.Assists
	ps.DamageDone += newStats.DamageDone
	ps.TicksAlive += newStats.TicksAlive
	ps.DeathPlacement += ps.DeathPlacement
	ps.Trades += newStats.Trades
	ps.Traded += newStats.Traded
	ps.OpeningKills += newStats.OpeningKills
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
	ps.KASTRounds += newStats.KASTRounds
	ps.Saves += newStats.Saves
	ps.Entries += newStats.Entries
	ps.ImpactPoints += newStats.ImpactPoints
	ps.WinPoints += newStats.WinPoints
	ps.AWPKills += newStats.AWPKills
	ps.RoundsFor += newStats.RoundsFor
	ps.RoundsAgainst += newStats.RoundsAgainst
	addMap(ps.UtilityThrown, newStats.UtilityThrown)
	ps.DamageTaken += newStats.DamageTaken
	ps.SupportDamage += newStats.SupportDamage
	ps.SupportRounds += newStats.SupportRounds
	ps.RoundsWithKills += newStats.RoundsWithKills
	ps.MostImpactfulPlayer += newStats.MostImpactfulPlayer
	ps.EffectiveAssistsContribution += newStats.EffectiveAssistsContribution
	ps.PlayerSide = Unknown

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
