package stats

type PlayerStats struct {
	Name                     string
	ImpactPoints             float64 `json:"impactPoints" end_of_match_sum:"true"`
	DamageDone               int     `json:"damageDone" end_of_match_sum:"true"`
	DamageTaken              int     `json:"damageTaken" end_of_match_sum:"true"`
	UtilityDamage            int     `json:"utilityDamage" end_of_match_sum:"true"`
	HEDamage                 int     `json:"heDamage" end_of_match_sum:"true"`
	FireDamage               int     `json:"fireDamage" end_of_match_sum:"true"`
	EnemiesFlashed           int     `json:"enemiesFlashed" end_of_match_sum:"true"`
	EnemyFlashDuration       float64 `json:"enemyFlashTime" end_of_match_sum:"true"`
	MostRecentFlasherID      uint64  `json:"mostRecentFlasher"`
	MostRecentFlashTickValue float64
	UtilityThrown            map[string]int `json:"utilityThrown" end_of_match_sum:"true"`
	Kills                    int            `json:"kills" end_of_match_sum:"true"`
	Deaths                   int            `json:"deaths" end_of_match_sum:"true"`
	DeathOrder               int            `json:"deathOrder" end_of_match_sum:"true"` // Order in which a player died. Ex.: Entry fraggers will have the lowest order (possibly 0)
	Assists                  int            `json:"assists" end_of_match_sum:"true"`
	KASTRound                int            `json:"kastRounds" end_of_match_sum:"true"`
	RoundsWithKills          int
	AWPKills                 int
	Headshots                int
	EAC                      int            `json:"eac" end_of_match_sum:"true"` // Effective Assists Contribution
	FlashAssists             int            `json:"flashAssists" end_of_match_sum:"true"`
	DamageList               map[uint64]int `json:"damageList"`
	SupportDamage            int            `json:"supportDamage" end_of_match_sum:"true"`
	SupportRound             int
	TicksAlive               int
	KillsList                map[uint64]int
	Traded                   int
	Trades                   int
	LurkerBlips              int `json:"lurkerBlips" end_of_match_sum:"true"`
}

func NewPlayerStats(name string) *PlayerStats {
	return &PlayerStats{
		Name: name,
	}
}
