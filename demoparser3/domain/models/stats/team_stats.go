package stats

type playedWon struct {
	Played int
	Won    int
}

type conditionOutcome struct {
	IsSet bool
	Won   bool
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
	FiveVFour      conditionOutcome   `json:"5v4" end_of_match_sum:"true"`
	FourVFive      conditionOutcome   `json:"4v5" end_of_match_sum:"true"`
	PistolRounds   playedWon          `json:"pistolRounds" end_of_match_sum:"true"`
	CTRounds       playedWon          `json:"ctRounds" end_of_match_sum:"true"`
	TRounds        playedWon          `json:"tRounds" end_of_match_sum:"true"`
	ClutchAttempt  ClutchAttempt
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
