package models

type Round struct {
	RoundNum       int `json:"roundNum"`
	StartingTick   int
	IsPrePlant     bool   `json:"isPrePlant"`
	IsPostPlant    bool   `json:"isPostPlant"`
	IsPistolRound  bool   `json:"isPistolRound"`
	IsPostRoundEnd bool   // The few seconds between round ending and before a new round starts.
	Winner         string `json:"winner"`
}

type RoundContext struct {
	StartingTick   int
	IsPrePlant     bool
	IsPostPlant    bool
	IsPostRoundEnd bool
	Winner         string
}

func NewRound(totalRoundsPlayed int, tick int) *Round {
	return &Round{
		RoundNum:     totalRoundsPlayed + 1,
		StartingTick: tick,
		IsPrePlant:   true,
	}
}

func (r *Round) Context() *RoundContext {
	return &RoundContext{
		StartingTick:   r.StartingTick,
		IsPrePlant:     r.IsPrePlant,
		IsPostPlant:    r.IsPostPlant,
		IsPostRoundEnd: r.IsPostRoundEnd,
		Winner:         r.Winner,
	}
}
