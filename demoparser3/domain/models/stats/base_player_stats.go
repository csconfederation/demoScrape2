package stats

type BasePlayerStats struct {
	Name           string
	DamageDone     int
	ImpactPoints   float64
	WinPoints      float64
	OpeningKills   int
	OpeningDeaths  int
	Kills          int
	Deaths         int
	KASTRounds     int
	Rounds         int
	DeathPlacement float64
	AWPKills       int
}

func NewBasePlayerStats(name string) *BasePlayerStats {
	return &BasePlayerStats{
		Name:   name,
		Rounds: 1,
	}
}

func (current *BasePlayerStats) AggregateBasePlayerStats(new *BasePlayerStats) {
	current.DamageDone += new.DamageDone
	current.ImpactPoints += new.ImpactPoints
	current.WinPoints += new.WinPoints
	current.OpeningKills += new.OpeningKills
	current.OpeningDeaths += new.OpeningDeaths
	current.Kills += new.Kills
	current.Deaths += new.Deaths
	current.KASTRounds += new.KASTRounds
	current.Rounds += new.Rounds
	current.DeathPlacement += new.DeathPlacement
	current.AWPKills += new.AWPKills
}
