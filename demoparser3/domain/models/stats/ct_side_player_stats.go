package stats

type CTSidePlayerStats struct {
	CTDamage                int     `json:"ctDamage"`
	CTImpactPoints          float64 `json:"ctImpactPoints"`
	CTWinPoints             float64 `json:"ctWinPoints"`
	CTOpeningKills          int     `json:"ctOK"`
	CTOpeningDeaths         int     `json:"ctOL"`
	CTKills                 int     `json:"ctKills"`
	CTDeaths                int     `json:"ctDeaths"`
	CTKAST                  float64 `json:"ctKAST"`
	CTKASTRounds            int     `json:"ctKASTRounds"`
	CTADR                   float64 `json:"ctADR"`
	CTTeamsWinPoints        float64 `json:"ctTeamsWinPoints"`
	CTWinPointsNormalizer   int     `json:"ctWinPointsNormalizer"`
	CTRounds                int     `json:"ctRounds"`
	CTRating                float64 `json:"ctRating"`
	CTImpactRating          float64 `json:"ctImpactRating"`
	CTAverageDeathPlacement int     `json:"ctADP"`
	CTAWPKills              int
}

func (ps *CTSidePlayerStats) Aggregate(newStats *PlayerStats) {
	// win points normalizer
	ps.CTImpactPoints += newStats.ImpactPoints
	ps.CTWinPoints += newStats.WinPoints
	ps.CTOpeningKills += newStats.OpeningKills
	ps.CTOpeningDeaths += newStats.OpeningDeaths
	ps.CTKills += newStats.Kills
	ps.CTDeaths += newStats.Deaths
	ps.CTKASTRounds += newStats.KASTRounds
	ps.CTDamage += newStats.DamageDone
	ps.CTAverageDeathPlacement += newStats.DeathPlacement
	//ps.CTWinPointsNormalizer +=
	ps.CTRounds += newStats.Rounds
	ps.CTAWPKills += newStats.AWPKills
}
