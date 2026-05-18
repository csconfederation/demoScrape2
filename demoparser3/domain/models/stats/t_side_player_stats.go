package stats

type TSidePlayerStats struct {
	TDamage                int     `json:"ctDamage"`
	TImpactPoints          float64 `json:"ctImpactPoints"`
	TWinPoints             float64 `json:"ctWinPoints"`
	TOpeningKills          int     `json:"ctOK"`
	TOpeningDeaths         int     `json:"ctOL"`
	TKills                 int     `json:"ctKills"`
	TDeaths                int     `json:"ctDeaths"`
	TKAST                  float64 `json:"ctKAST"`
	TKASTRounds            int     `json:"ctKASTRounds"`
	TADR                   float64 `json:"ctADR"`
	TTeamsWinPoints        float64 `json:"ctTeamsWinPoints"`
	TWinPointsNormalizer   int     `json:"ctWinPointsNormalizer"`
	TRounds                int     `json:"ctRounds"`
	TRating                float64 `json:"ctRating"`
	TImpactRating          float64 `json:"ctImpactRating"`
	TAverageDeathPlacement int     `json:"ctADP"`
	TAWPKills              int
}

func (ps *TSidePlayerStats) Aggregate(newStats *PlayerStats) {
	// win points normalizer
	ps.TImpactPoints += newStats.ImpactPoints
	ps.TWinPoints += newStats.WinPoints
	ps.TOpeningKills += newStats.OpeningKills
	ps.TOpeningDeaths += newStats.OpeningDeaths
	ps.TKills += newStats.Kills
	ps.TDeaths += newStats.Deaths
	ps.TKASTRounds += newStats.KASTRounds
	ps.TDamage += newStats.DamageDone
	ps.TAverageDeathPlacement += newStats.DeathPlacement
	//ps.TWinPointsNormalizer +=
	ps.TRounds += newStats.Rounds
	ps.TAWPKills += newStats.AWPKills
	// Lurk Rounds
}
