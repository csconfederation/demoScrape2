package stats

type playedWon struct {
	Played int
	Won    int
}

func (pw *playedWon) add(other playedWon) {
	pw.Played += other.Played
	pw.Won += other.Won
}

type ClutchAttempt struct {
	PlayerID     uint64
	ClutchValue  int
	IsSuccessful bool
}

type TeamStats struct {
	TeamName       string
	TeamMembers    []uint64
	MembersAlive   []uint64
	DeathPlacement map[uint64]float64 `json:"deathPlacement"`
	FiveVFour      playedWon          `json:"5v4" end_of_match_sum:"true"`
	FourVFive      playedWon          `json:"4v5" end_of_match_sum:"true"`
	PistolRounds   playedWon          `json:"pistolRounds" end_of_match_sum:"true"`
	CTRounds       playedWon          `json:"ctRounds" end_of_match_sum:"true"`
	TRounds        playedWon          `json:"tRounds" end_of_match_sum:"true"`
	ClutchAttempt  ClutchAttempt
	Saves          int
}

func NewTeamStats(teamName string, connectedTeamPlayers []uint64) *TeamStats {
	return &TeamStats{
		TeamName:     teamName,
		TeamMembers:  connectedTeamPlayers,
		MembersAlive: connectedTeamPlayers,
	}
}

func NewClutchAttempt(playerID uint64, clutchValue int) *ClutchAttempt {
	return &ClutchAttempt{
		PlayerID:    playerID,
		ClutchValue: clutchValue,
	}
}

func (ts *TeamStats) Aggregate(newStats *TeamStats) {
	ts.PistolRounds.add(newStats.PistolRounds)
	ts.FiveVFour.add(newStats.FiveVFour)
	ts.FourVFive.add(newStats.FourVFive)
	ts.Saves += newStats.Saves
	// ts.Clutches
	ts.CTRounds.add(newStats.CTRounds)
	ts.TRounds.add(newStats.TRounds)
}
