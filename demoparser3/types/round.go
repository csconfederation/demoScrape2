package types

import (
	"github.com/csconfederation/demoparser3/logger"
	"github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs/common"
	"github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs/events"
)

type Round struct {
	StartingTick    int
	RoundNum        int                     `json:"roundNum"`
	IsPrePlant      bool                    `json:"isPrePlant"`
	IsPostPlant     bool                    `json:"isPostPlant"`
	IsFinalRound    bool                    `json:"isFinalRound"`
	IsPistolRound   bool                    `json:"isPistolRound"`
	TMoney          bool                    `json:"tMoney"`
	Planter         *common.Player          `json:"planter"`
	Defuser         *common.Player          `json:"defuser"`
	BombStartTick   int                     `json:"bombStartTick"`
	RoundEndReason  events.RoundEndReason   `json:"roundEndReason"`
	AllPlayersStats map[uint64]*PlayerStats `json:"allPlayersStats"`
	Winner          *Team                   `json:"winner"`
	Loser           *Team                   `json:"loser"`
	PostWinCon      bool                    `json:"postWinCon"` // Post win condition. The few seconds between round ending and before a new round starts.
	TeamStats       map[string]*TeamStats   `json:"teamStats"`
	IsValid         bool                    `json:"isValid"`
}

func NewRound(tick int) *Round {
	return &Round{
		StartingTick:    tick,
		IsPrePlant:      true,
		AllPlayersStats: make(map[uint64]*PlayerStats),
		TeamStats:       make(map[string]*TeamStats),
	}
}

func (round *Round) BombPlanted(planter *common.Player) {
	round.IsPrePlant = false
	round.IsPostPlant = false
	round.TMoney = true
	if planter == nil {
		logger.Warn("Bomb planted by a nil player, possibly POV demo")
		return
	}
	round.Planter = planter
}

func (round *Round) BombDefused(defuser *common.Player) {

	round.IsPrePlant = false
	round.IsPostPlant = true
	round.PostWinCon = true
	round.RoundEndReason = events.RoundEndReasonBombDefused
	if defuser == nil {
		logger.Warn("Defuser is nil")
		return
	}
	round.Defuser = defuser
	round.AllPlayersStats[defuser.SteamID64].ImpactPoints += 0.5
}

func (round *Round) BombExplode() {

	round.IsPrePlant = false
	round.IsPostPlant = false
	round.PostWinCon = true
	round.RoundEndReason = events.RoundEndReasonTargetBombed
	if round.Planter == nil {
		logger.Warn("Planter is nil")
		return
	}
	round.AllPlayersStats[round.Planter.SteamID64].ImpactPoints += 0.5
}

func (round *Round) RoundEnd(winner, loser *Team, reason events.RoundEndReason) {

	if &winner == nil {
		logger.Warn("Winner is nil - RoundEnd")
		return
	}

	round.Winner = winner.Clone()
	round.Loser = loser.Clone()
	round.RoundEndReason = reason
	round.IsPrePlant = false
	round.IsPostPlant = false
	round.PostWinCon = true
	round.IsFinalRound = isFinalRound(winner.Score, loser.Score)
	if round.IsPistolRound {
		round.TeamStats[round.Winner.Name].PistolRounds.Won += 1
	}

	// TODO: Handle this processing at the end

	// Assuming MR + 1 == 13
	//if winner.Score == 13 && loser.Score < 12 {
	//
	//} else if winner.Score > 12 { //check for OT win
	//	overtime := ((winner.Score+loser.Score)-24-1)/6 + 1
	//	//OT win
	//	if (winner.Score-12-1)/3 == overtime {
	//		processRoundFinal(true)
	//	}
	//}
	//
	//game.Flags.InRound = false
	//game.PotentialRound.EndingTick = p.GameState().IngameTick()
	//game.Flags.RoundIntegrityEndOfficial = p.GameState().TotalRoundsPlayed()
	//
	//log.Debug("We are processing round final stuff", game.Flags.RoundIntegrityEndOfficial)
	//log.Debug(len(game.Rounds))
	//
	////we have the entire round uninterrupted
	//if game.Flags.RoundIntegrityStart == game.Flags.RoundIntegrityEnd && game.Flags.RoundIntegrityEnd == game.Flags.RoundIntegrityEndOfficial {
	//	game.PotentialRound.IntegrityCheck = true
	//
	//	//check team stats
	//	if game.PotentialRound.TeamStats[game.PotentialRound.WinnerClanName].Pistols == 1 {
	//		game.PotentialRound.TeamStats[game.PotentialRound.WinnerClanName].PistolsW = 1
	//	}
	//	if game.PotentialRound.TeamStats[game.PotentialRound.WinnerClanName].FourVFiveS == 1 {
	//		game.PotentialRound.TeamStats[game.PotentialRound.WinnerClanName].FourVFiveW = 1
	//	} else if game.PotentialRound.TeamStats[game.PotentialRound.WinnerClanName].FiveVFourS == 1 {
	//		game.PotentialRound.TeamStats[game.PotentialRound.WinnerClanName].FiveVFourW = 1
	//	}
	//	if game.PotentialRound.TeamStats[game.PotentialRound.WinnerClanName].TR == 1 {
	//		game.PotentialRound.TeamStats[game.PotentialRound.WinnerClanName].TRW = 1
	//	} else if game.PotentialRound.TeamStats[game.PotentialRound.WinnerClanName].CtR == 1 {
	//		game.PotentialRound.TeamStats[game.PotentialRound.WinnerClanName].CtRW = 1
	//	}
	//
	//	//set the clutch
	//	if game.PotentialRound.WinnerENUM == 2 && game.Flags.TClutchSteam != 0 {
	//		game.PotentialRound.TeamStats[game.PotentialRound.WinnerClanName].Clutches = 1
	//		game.PotentialRound.PlayerStats[game.Flags.TClutchSteam].ImpactPoints += clutchBonus[game.Flags.TClutchVal]
	//		switch game.Flags.TClutchVal {
	//		case 1:
	//			game.PotentialRound.PlayerStats[game.Flags.TClutchSteam].Cl_1 = 1
	//		case 2:
	//			game.PotentialRound.PlayerStats[game.Flags.TClutchSteam].Cl_2 = 1
	//		case 3:
	//			game.PotentialRound.PlayerStats[game.Flags.TClutchSteam].Cl_3 = 1
	//		case 4:
	//			game.PotentialRound.PlayerStats[game.Flags.TClutchSteam].Cl_4 = 1
	//		case 5:
	//			game.PotentialRound.PlayerStats[game.Flags.TClutchSteam].Cl_5 = 1
	//		}
	//	} else if game.PotentialRound.WinnerENUM == 3 && game.Flags.CtClutchSteam != 0 {
	//		game.PotentialRound.TeamStats[game.PotentialRound.WinnerClanName].Clutches = 1
	//		game.PotentialRound.PlayerStats[game.Flags.CtClutchSteam].ImpactPoints += clutchBonus[game.Flags.CtClutchVal]
	//		switch game.Flags.CtClutchVal {
	//		case 1:
	//			game.PotentialRound.PlayerStats[game.Flags.CtClutchSteam].Cl_1 = 1
	//		case 2:
	//			game.PotentialRound.PlayerStats[game.Flags.CtClutchSteam].Cl_2 = 1
	//		case 3:
	//			game.PotentialRound.PlayerStats[game.Flags.CtClutchSteam].Cl_3 = 1
	//		case 4:
	//			game.PotentialRound.PlayerStats[game.Flags.CtClutchSteam].Cl_4 = 1
	//		case 5:
	//			game.PotentialRound.PlayerStats[game.Flags.CtClutchSteam].Cl_5 = 1
	//		}
	//	}
	//
	//	//add multikills & saves & misc
	//	highestImpactPoints := 0.0
	//	mipPlayers := 0
	//	for _, player := range (game.PotentialRound).PlayerStats {
	//		if player.Deaths == 0 {
	//			player.KastRounds = 1
	//			if player.TeamENUM != game.PotentialRound.WinnerENUM {
	//				player.Saves = 1
	//				game.PotentialRound.TeamStats[player.TeamClanName].Saves = 1
	//			}
	//		}
	//		steamId64, _ := strconv.ParseUint(player.SteamID, 10, 64)
	//		game.PotentialRound.PlayerStats[steamId64].ImpactPoints += player.KillPoints
	//		game.PotentialRound.PlayerStats[steamId64].ImpactPoints += float64(player.Damage) / float64(250)
	//		game.PotentialRound.PlayerStats[steamId64].ImpactPoints += multiKillBonus[player.Kills]
	//
	//		switch player.Kills {
	//		case 2:
	//			player.TwoK = 1
	//		case 3:
	//			player.ThreeK = 1
	//		case 4:
	//			player.FourK = 1
	//		case 5:
	//			player.FiveK = 1
	//		}
	//
	//		if player.ImpactPoints > highestImpactPoints {
	//			highestImpactPoints = player.ImpactPoints
	//		}
	//
	//		if player.TeamENUM == game.PotentialRound.WinnerENUM {
	//			player.WinPoints = player.ImpactPoints
	//
	//			player.RF = 1
	//		} else {
	//			player.RA = 1
	//		}
	//	}
	//
	//	for _, player := range (game.PotentialRound).PlayerStats {
	//		if player.ImpactPoints == highestImpactPoints {
	//			mipPlayers += 1
	//		}
	//	}
	//	for _, player := range (game.PotentialRound).PlayerStats {
	//		if player.ImpactPoints == highestImpactPoints {
	//			player.Mip = 1.0 / float64(mipPlayers)
	//		}
	//	}
	//
	//	//check the lurk
	//	var susLurker uint64
	//	susLurkBlips := 0
	//	invalidLurk := false
	//	for _, player := range game.PotentialRound.PlayerStats {
	//		if player.Side == 2 {
	//			if player.LurkerBlips > susLurkBlips {
	//				susLurkBlips = player.LurkerBlips
	//				steamId64, _ := strconv.ParseUint(player.SteamID, 10, 64)
	//				susLurker = steamId64
	//			}
	//		}
	//	}
	//	for _, player := range game.PotentialRound.PlayerStats {
	//		if player.Side == 2 {
	//			steamId64, _ := strconv.ParseUint(player.SteamID, 10, 64)
	//			if player.LurkerBlips == susLurkBlips && steamId64 != susLurker {
	//				invalidLurk = true
	//			}
	//		}
	//	}
	//	if !invalidLurk && susLurkBlips > 3 {
	//		game.PotentialRound.PlayerStats[susLurker].LurkRounds = 1
	//	}
	//
	//	//add our valid round
	//	game.Rounds = append(game.Rounds, game.PotentialRound)
	//}
	//if lastRound {
	//	//game.Flags.RoundIntegrityEndOfficial += 1
	//	game.TotalRounds = game.Flags.RoundIntegrityEndOfficial
	//	game.Flags.IsGameLive = false
	//}

	//endRound function functionality

}

func (round *Round) RoundFreezetimeEnd(roundsPlayed int) {
	round.IsPistolRound = (roundsPlayed%12)+1 == 1
	if round.IsFinalRound {
		for team := range round.TeamStats {
			round.TeamStats[team].PistolRounds.Total += 1
		}
	}
}

func isFinalRound(winnerScore, loserScore int) bool {
	if winnerScore == 13 && loserScore < 12 {
		return true
	}

	overtime := ((winnerScore+loserScore)-24-1)/6 + 1
	//OT win
	if (winnerScore-12-1)/3 == overtime && overtime > 0 {
		return true
	}

	return false
}

//type Round struct {
//	RoundNum            int8                    `json:"roundNum"`
//	StartingTick        int                     `json:"startingTick"`
//	EndingTick          int                     `json:"endingTick"`
//	PlayerStats         map[uint64]*PlayerStats `json:"playerStats"`
//	TeamStats           map[string]*TeamStats   `json:"teamStats"`
//	InitTerroristCount  int                     `json:"initTerroristCount"`
//	InitCTerroristCount int                     `json:"initCTerroristCount"`
//	WinnerClanName      string                  `json:"winnerClanName"`
//	WinnerENUM          int                     `json:"winnerENUM"` //this effectively represents the side that won: 2 (T) or 3 (CT)
//	IntegrityCheck      bool                    `json:"integrityCheck"`
//	Planter             uint64                  `json:"planter"`
//	Defuser             uint64                  `json:"defuser"`
//	EndDueToBombEvent   bool                    `json:"endDueToBombEvent"`
//	WinTeamDmg          int                     `json:"winTeamDmg"`
//	KnifeRound          bool                    `json:"knifeRound"`
//	RoundEndReason      string                  `json:"roundEndReason"`
//
//	WPAlog        []*WPALog `json:"WPAlog"`
//	BombStartTick int       `json:"bombStartTick"`
//}
