package stats

type Averages struct {
	RoundNormalizer    float64
	ImpactRoundAverage float64
	KillRoundAverage   float64
	DeathRoundAverage  float64
	KASTRoundAverage   float64
	ADRAverage         float64
}

func NewAverages(roundNormalizer, impactRoundAvg, killRoundAvg, deathRoundAvg, kastRoundAvg, adrAvg float64) *Averages {
	return &Averages{
		RoundNormalizer:    roundNormalizer,
		ImpactRoundAverage: impactRoundAvg,
		KillRoundAverage:   killRoundAvg,
		DeathRoundAverage:  deathRoundAvg,
		KASTRoundAverage:   kastRoundAvg,
		ADRAverage:         adrAvg,
	}
}

func (as *Averages) Aggregate(new *Averages) {
	as.RoundNormalizer += new.RoundNormalizer
	as.ImpactRoundAverage += new.ImpactRoundAverage
	as.KillRoundAverage += new.KillRoundAverage
	as.DeathRoundAverage += new.DeathRoundAverage
	as.KASTRoundAverage += new.KASTRoundAverage
	as.ADRAverage += new.ADRAverage
}

func (as *Averages) Normalize() {
	as.ImpactRoundAverage /= as.RoundNormalizer
	as.KillRoundAverage /= as.RoundNormalizer
	as.DeathRoundAverage /= as.RoundNormalizer
	as.KASTRoundAverage /= as.RoundNormalizer
	as.ADRAverage /= as.RoundNormalizer
}
