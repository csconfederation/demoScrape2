package stats

type SidedPlayerStats struct {
	BasePlayerStats
	KAST                float64 `json:"ctKAST"`
	ADR                 float64 `json:"ctADR"`
	TeamWinPoints       float64 `json:"ctTeamsWinPoints"`
	WinPointsNormalizer int     `json:"ctWinPointsNormalizer"`
}

func NewSidedPlayerStats(name string) *SidedPlayerStats {
	return &SidedPlayerStats{
		BasePlayerStats: *NewBasePlayerStats(name),
	}
}

func (sps *SidedPlayerStats) Aggregate(newStats *SidedPlayerStats) {
	sps.ImpactPoints += newStats.ImpactPoints
	sps.WinPoints += newStats.WinPoints
	sps.OpeningKills += newStats.OpeningKills
	sps.OpeningDeaths += newStats.OpeningDeaths
	sps.Kills += newStats.Kills
	sps.Deaths += newStats.Deaths
	sps.KASTRounds += newStats.KASTRounds
	sps.DamageDone += newStats.DamageDone
	sps.DeathPlacement += newStats.DeathPlacement
	//sps.WinPointsNormalizer +=
	sps.Rounds += newStats.Rounds
	sps.AWPKills += newStats.AWPKills
}
