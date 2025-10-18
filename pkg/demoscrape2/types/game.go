package types

import (
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
	log "github.com/sirupsen/logrus"
)

const MR = 12

var clutchBonus = [...]float64{0, 0.2, 0.6, 1.2, 2, 3}
var multiKillBonus = [...]float64{0, 0, 0.3, 0.7, 1.2, 2}

type Game struct {
	//winnerID         int
	CoreID                   string                  `json:"coreID"`
	MapNum                   int                     `json:"mapNum"`
	WinnerClanName           string                  `json:"winnerClanName"`
	Result                   string                  `json:"result"`
	Rounds                   []*Round                `json:"rounds"`
	PotentialRound           *Round                  `json:"potentialRound"`
	Teams                    map[string]*Team        `json:"teams"`
	Flags                    Flag                    `json:"flags"`
	MapName                  string                  `json:"mapName"`
	TickRate                 int                     `json:"tickRate"`
	TickLength               int                     `json:"tickLength"`
	RoundsToWin              int                     `json:"roundsToWin"` //30 or 16
	TotalPlayerStats         map[uint64]*PlayerStats `json:"totalPlayerStats"`
	CtPlayerStats            map[uint64]*PlayerStats `json:"ctPlayerStats"`
	TPlayerStats             map[uint64]*PlayerStats `json:"TPlayerStats"`
	TotalTeamStats           map[string]*TeamStats   `json:"totalTeamStats"`
	ReconnectedPlayers       map[uint64]bool         `json:"reconnectedPlayers"`       // Map of SteamID to reconnection status
	ConnectedAfterRoundStart map[uint64]bool         `json:"ConnectedAfterRoundStart"` // Map of SteamID to reconnection status
	PlayerOrder              []uint64                `json:"playerOrder"`
	TeamOrder                []string                `json:"teamOrder"`
	TotalRounds              int                     `json:"totalRounds"`
	TotalWPAlog              []*WPALog               `json:"totalWPAlog"`
}

func NewGame() *Game {
	return &Game{
		Rounds:                   make([]*Round, 0),
		Teams:                    make(map[string]*Team),
		TotalPlayerStats:         make(map[uint64]*PlayerStats),
		CtPlayerStats:            make(map[uint64]*PlayerStats),
		TPlayerStats:             make(map[uint64]*PlayerStats),
		TotalTeamStats:           make(map[string]*TeamStats),
		ReconnectedPlayers:       make(map[uint64]bool),
		ConnectedAfterRoundStart: make(map[uint64]bool),
		PlayerOrder:              make([]uint64, 0),
		TeamOrder:                make([]string, 0),
		TotalWPAlog:              make([]*WPALog, 0),
		TickRate:                 64,
	}
}

func (game *Game) SetTickLength(tickLength int) {
	game.TickLength = tickLength
}

func (game *Game) Start(tickRate int, teamTerrorists, teamCounterTerrorists *common.TeamState) {
	game.Flags.HasGameStarted = true
	game.Flags.IsGameLive = true
	log.Debug("GAME HAS STARTED!!!")

	// In case the tickRate is 0 we want to re-set it based on the tickInterval now that the Game has hasGameStarted
	if game.TickRate == 0 {
		game.TickRate = tickRate
	}

	terrorists := NewTeamFromTeamState(teamTerrorists, game.Teams, game.PotentialRound)
	counterTerrorists := NewTeamFromTeamState(teamCounterTerrorists, game.Teams, game.PotentialRound)

	game.Teams[terrorists.Name] = terrorists
	game.Teams[counterTerrorists.Name] = counterTerrorists

	//only handling normal length matches
	game.RoundsToWin = MR + 1
}

func (game *Game) ResetFlags() {
	game.Flags.PrePlant = true
	game.Flags.PostPlant = false
	game.Flags.PostWinCon = false
	game.Flags.TClutchVal = 0
	game.Flags.CtClutchVal = 0
	game.Flags.TClutchSteam = 0
	game.Flags.CtClutchSteam = 0
	game.Flags.TMoney = false
	game.Flags.OpeningKill = true
	game.Flags.LastTickProcessed = 0
	game.Flags.TicksProcessed = 0
	game.Flags.DidRoundEndFire = false
	game.Flags.RoundStartedAt = 0
	game.Flags.HaveInitRound = false
}

func (game *Game) ProcessRoundOnWinCon(winnerState *common.TeamState, totalRoundsPlayed int) {
	game.Flags.RoundIntegrityEnd = totalRoundsPlayed
	log.Debug("We are processing round win con stuff", game.Flags.RoundIntegrityEnd)

	game.TotalRounds = game.Flags.RoundIntegrityEnd

	game.Flags.PrePlant = false
	game.Flags.PostPlant = false
	game.Flags.PostWinCon = true

	//set winner
	game.PotentialRound.WinnerClanName = ValidateTeamName(winnerState, game.Teams, game.PotentialRound)
	//log.Debug("We think this team won", winnerClanName)
	if !game.PotentialRound.KnifeRound {
		game.Teams[game.PotentialRound.WinnerClanName].Score += 1
	}
	//go through and set our WPAlog output to the winner
	for _, wpalog := range game.PotentialRound.WPAlog {
		wpalog.Winner = game.PotentialRound.WinnerENUM - 2
	}
}

// TODO: rewrite this better

func (game *Game) ProcessRoundFinal(isLastRound bool, endingTick, totalRoundsPlayed int) {
	game.Flags.InRound = false
	game.PotentialRound.EndingTick = endingTick
	game.Flags.RoundIntegrityEndOfficial = totalRoundsPlayed

	log.Debug("We are processing round final stuff", game.Flags.RoundIntegrityEndOfficial)
	log.Debug(len(game.Rounds))

	//we have the entire round uninterrupted
	if game.Flags.RoundIntegrityStart == game.Flags.RoundIntegrityEnd && game.Flags.RoundIntegrityEnd == game.Flags.RoundIntegrityEndOfficial {
		game.PotentialRound.IntegrityCheck = true

		//check team stats
		if game.PotentialRound.TeamStats[game.PotentialRound.WinnerClanName].Pistols == 1 {
			game.PotentialRound.TeamStats[game.PotentialRound.WinnerClanName].PistolsW = 1
		}
		if game.PotentialRound.TeamStats[game.PotentialRound.WinnerClanName].FourVFiveS == 1 {
			game.PotentialRound.TeamStats[game.PotentialRound.WinnerClanName].FourVFiveW = 1
		} else if game.PotentialRound.TeamStats[game.PotentialRound.WinnerClanName].FiveVFourS == 1 {
			game.PotentialRound.TeamStats[game.PotentialRound.WinnerClanName].FiveVFourW = 1
		}
		if game.PotentialRound.TeamStats[game.PotentialRound.WinnerClanName].TR == 1 {
			game.PotentialRound.TeamStats[game.PotentialRound.WinnerClanName].TRW = 1
		} else if game.PotentialRound.TeamStats[game.PotentialRound.WinnerClanName].CtR == 1 {
			game.PotentialRound.TeamStats[game.PotentialRound.WinnerClanName].CtRW = 1
		}

		//set the clutch
		if game.PotentialRound.WinnerENUM == 2 && game.Flags.TClutchSteam != 0 {
			game.PotentialRound.TeamStats[game.PotentialRound.WinnerClanName].Clutches = 1
			game.PotentialRound.PlayerStats[game.Flags.TClutchSteam].ImpactPoints += clutchBonus[game.Flags.TClutchVal]
			switch game.Flags.TClutchVal {
			case 1:
				game.PotentialRound.PlayerStats[game.Flags.TClutchSteam].Cl_1 = 1
			case 2:
				game.PotentialRound.PlayerStats[game.Flags.TClutchSteam].Cl_2 = 1
			case 3:
				game.PotentialRound.PlayerStats[game.Flags.TClutchSteam].Cl_3 = 1
			case 4:
				game.PotentialRound.PlayerStats[game.Flags.TClutchSteam].Cl_4 = 1
			case 5:
				game.PotentialRound.PlayerStats[game.Flags.TClutchSteam].Cl_5 = 1
			}
		} else if game.PotentialRound.WinnerENUM == 3 && game.Flags.CtClutchSteam != 0 {
			game.PotentialRound.TeamStats[game.PotentialRound.WinnerClanName].Clutches = 1
			game.PotentialRound.PlayerStats[game.Flags.CtClutchSteam].ImpactPoints += clutchBonus[game.Flags.CtClutchVal]
			switch game.Flags.CtClutchVal {
			case 1:
				game.PotentialRound.PlayerStats[game.Flags.CtClutchSteam].Cl_1 = 1
			case 2:
				game.PotentialRound.PlayerStats[game.Flags.CtClutchSteam].Cl_2 = 1
			case 3:
				game.PotentialRound.PlayerStats[game.Flags.CtClutchSteam].Cl_3 = 1
			case 4:
				game.PotentialRound.PlayerStats[game.Flags.CtClutchSteam].Cl_4 = 1
			case 5:
				game.PotentialRound.PlayerStats[game.Flags.CtClutchSteam].Cl_5 = 1
			}
		}

		//add multikills & saves & misc
		highestImpactPoints := 0.0
		mipPlayers := 0
		for _, player := range (game.PotentialRound).PlayerStats {
			if player.Deaths == 0 {
				player.KastRounds = 1
				if player.TeamENUM != game.PotentialRound.WinnerENUM {
					player.Saves = 1
					game.PotentialRound.TeamStats[player.TeamClanName].Saves = 1
				}
			}
			steamId64 := player.SteamID
			game.PotentialRound.PlayerStats[steamId64].ImpactPoints += player.KillPoints
			game.PotentialRound.PlayerStats[steamId64].ImpactPoints += float64(player.Damage) / float64(250)
			game.PotentialRound.PlayerStats[steamId64].ImpactPoints += multiKillBonus[player.Kills]

			switch player.Kills {
			case 2:
				player.TwoK = 1
			case 3:
				player.ThreeK = 1
			case 4:
				player.FourK = 1
			case 5:
				player.FiveK = 1
			}

			if player.ImpactPoints > highestImpactPoints {
				highestImpactPoints = player.ImpactPoints
			}

			if player.TeamENUM == game.PotentialRound.WinnerENUM {
				player.WinPoints = player.ImpactPoints

				player.RF = 1
			} else {
				player.RA = 1
			}
		}

		for _, player := range (game.PotentialRound).PlayerStats {
			if player.ImpactPoints == highestImpactPoints {
				mipPlayers += 1
			}
		}
		for _, player := range (game.PotentialRound).PlayerStats {
			if player.ImpactPoints == highestImpactPoints {
				player.Mip = 1.0 / float64(mipPlayers)
			}
		}

		//check the lurk
		var susLurker uint64
		susLurkBlips := 0
		invalidLurk := false
		for _, player := range game.PotentialRound.PlayerStats {
			if player.Side == 2 {
				if player.LurkerBlips > susLurkBlips {
					susLurkBlips = player.LurkerBlips
					susLurker = player.SteamID
				}
			}
		}
		for _, player := range game.PotentialRound.PlayerStats {
			if player.Side == 2 {
				if player.LurkerBlips == susLurkBlips && player.SteamID != susLurker {
					invalidLurk = true
				}
			}
		}
		if !invalidLurk && susLurkBlips > 3 {
			game.PotentialRound.PlayerStats[susLurker].LurkRounds = 1
		}

		//add our valid round
		game.Rounds = append(game.Rounds, game.PotentialRound)
	}
	if isLastRound {
		//game.Flags.RoundIntegrityEndOfficial += 1
		game.TotalRounds = game.Flags.RoundIntegrityEndOfficial
		game.Flags.IsGameLive = false
	}

	//endRound function functionality
}
