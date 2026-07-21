package events

type GameStart struct {
	MapName string
}

type MatchStart struct {
	CTSide    string
	TSide     string
	CTMembers map[uint64]string
	TMembers  map[uint64]string
}

type GameHalfEnded struct{}

type FrameDone struct {
	Tick      int
	Distances map[uint64]map[uint64]float64
}

type ScoreUpdated struct {
	TeamName string
	OldScore int
	NewScore int
}

type UtilityThrown struct {
	ThrowerID   uint64
	UtilityType string
}

type Death struct {
	VictimID       uint64
	KillerID       uint64
	VictimTeamName string
	Tick           int
}

type Kill struct {
	KillerID             uint64
	VictimID             uint64
	IsAWPKill            bool
	IsHeadshot           bool
	Tick                 int
	KillerTeamName       string
	IsAssisted           bool
	FlashAssisted        bool
	KillerEquipmentValue float64
	VictimEquipmentValue float64
}

type RoundStart struct {
	TotalRoundsPlayed  int
	ConnectedCTPlayers []uint64
	ConnectedTPlayers  []uint64
	Tick               int
}

type RoundEnd struct {
	WinningTeamName string
}

type RoundEndOfficial struct{}

type BombPlanted struct{}

type BombExplode struct {
	PlanterID uint64
}

type BombDefused struct {
	DefuserID uint64
}

type FinalizeStats struct{}

type PlayerHurt struct {
	AttackerID   uint64
	VictimID     uint64
	Weapon       string
	HealthDamage int
	IsUtility    bool
	IsFireDamage bool
	IsHEDamage   bool
}

type PlayerFlashed struct {
	AttackerID    uint64
	VictimID      uint64
	FlashDuration float64
	Tick          int
}

type Assist struct {
	VictimID        uint64
	AssisterID      uint64
	IsAssistedFlash bool
	Tick            int
}
