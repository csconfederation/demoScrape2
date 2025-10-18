package types

import (
	dem "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
	log "github.com/sirupsen/logrus"
)

type Round struct {
	//round value
	RoundNum            int8                    `json:"roundNum"`
	StartingTick        int                     `json:"startingTick"`
	EndingTick          int                     `json:"endingTick"`
	PlayerStats         map[uint64]*PlayerStats `json:"playerStats"`
	TeamStats           map[string]*TeamStats   `json:"teamStats"`
	InitTerroristCount  int                     `json:"initTerroristCount"`
	InitCTerroristCount int                     `json:"initCTerroristCount"`
	WinnerClanName      string                  `json:"winnerClanName"`
	//winnerID            int //this is the unique ID which should not change BUT IT DOES
	WinnerENUM         int     `json:"winnerENUM"` //this effectively represents the side that won: 2 (T) or 3 (CT)
	IntegrityCheck     bool    `json:"integrityCheck"`
	Planter            uint64  `json:"planter"`
	Defuser            uint64  `json:"defuser"`
	EndDueToBombEvent  bool    `json:"endDueToBombEvent"`
	WinTeamDmg         int     `json:"winTeamDmg"`
	ServerNormalizer   int     `json:"serverNormalizer"`
	ServerImpactPoints float64 `json:"serverImpactPoints"`
	KnifeRound         bool    `json:"knifeRound"`
	RoundEndReason     string  `json:"roundEndReason"`

	WPAlog        []*WPALog `json:"WPAlog"`
	BombStartTick int       `json:"bombStartTick"`
}

func NewRound(roundNum int8, startingTick int) *Round {
	return &Round{
		RoundNum:     roundNum,
		StartingTick: startingTick,
		PlayerStats:  make(map[uint64]*PlayerStats),
		TeamStats:    make(map[string]*TeamStats),
		WPAlog:       make([]*WPALog, 0),
	}
}

func (round *Round) IsRoundFinalInHalf() bool {
	return round.RoundNum%MR == 0 || (round.RoundNum > (MR*2) && round.RoundNum%3 == 0)
}

func (round *Round) Start(game *Game, parser dem.Parser) {
	// Reset the connectedAfterRoundStart
	game.ConnectedAfterRoundStart = make(map[uint64]bool)

	game.Flags.RoundIntegrityStart = parser.GameState().TotalRoundsPlayed() + 1
	log.Debug("We are starting round", game.Flags.RoundIntegrityStart)

	newRound := NewRound(int8(game.Flags.RoundIntegrityStart), parser.GameState().IngameTick())

	//set players in playerStats for the round
	terrorists := parser.GameState().TeamTerrorists()
	counterTerrorists := parser.GameState().TeamCounterTerrorists()

	newRound.initTeamPlayer(terrorists, game, parser)
	newRound.initTeamPlayer(counterTerrorists, game, parser)

	//set teams in teamStats for the round
	newRound.TeamStats[ValidateTeamName(terrorists, game.Teams, newRound)] = NewTR()
	newRound.TeamStats[ValidateTeamName(counterTerrorists, game.Teams, newRound)] = NewCtR()

	// Reset round
	game.PotentialRound = newRound

	//track the number of people alive for clutch checking and record keeping
	game.Flags.TAlive = len(GetTeamMembers(terrorists, game, parser))
	game.Flags.CtAlive = len(GetTeamMembers(counterTerrorists, game, parser))
	game.PotentialRound.InitTerroristCount = game.Flags.TAlive
	game.PotentialRound.InitCTerroristCount = game.Flags.CtAlive

	game.ResetFlags()
}

func (round *Round) initTeamPlayer(team *common.TeamState, game *Game, parser dem.Parser) {
	for _, teamMember := range GetTeamMembers(team, game, parser) {
		player := NewPlayerStatsFromPlayer(teamMember, team, game)
		round.PlayerStats[teamMember.SteamID64] = player
	}
}
