package stats

type DerivedPlayerStats struct {
	AverageTimeAlive         int
	AverageDeathPlacement    float64
	KAST                     float64
	KillPointAverage         float64
	ImpactInWinningRounds    float64
	ADR                      float64
	DamageRatingDifferential float64
	KillRatio                float64
	TradeRatio               float64
	UtilityThrown            int
}
