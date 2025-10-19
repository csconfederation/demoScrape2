package types

import (
	"math"
	"reflect"

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

func (game *Game) removeInvalidRounds() {
	//we want to remove bad rounds (knife/veto rounds, incomplete rounds, redo rounds)
	validRoundsMap := make(map[int8]bool)
	validRounds := make([]*Round, 0)
	lastProcessedRoundNum := game.Rounds[len(game.Rounds)-1].RoundNum + 1
	for i := len(game.Rounds) - 1; i >= 0; i-- {
		_, validRoundExists := validRoundsMap[game.Rounds[i].RoundNum]
		if game.Rounds[i].IntegrityCheck && !game.Rounds[i].KnifeRound && !validRoundExists {
			if game.Rounds[i].RoundNum == lastProcessedRoundNum-1 {
				//this i-th round is good to add
				validRoundsMap[game.Rounds[i].RoundNum] = true
				validRounds = append(validRounds, game.Rounds[i])
				lastProcessedRoundNum = game.Rounds[i].RoundNum
			}
		} else {
			//this i-th round is bad and we will remove it
		}
	}
	for i, j := 0, len(validRounds)-1; i < j; i, j = i+1, j-1 {
		validRounds[i], validRounds[j] = validRounds[j], validRounds[i]
	}
	game.Rounds = validRounds
}

func (game *Game) EndOfMatchProcessing() {
	game.removeInvalidRounds()
	for i := len(game.Rounds) - 1; i >= 0; i-- {
		game.TotalWPAlog = append(game.TotalWPAlog, game.Rounds[i].WPAlog...)

		for teamName, team := range game.Rounds[i].TeamStats {
			if game.TotalTeamStats[teamName] == nil && teamName != "" {
				game.TotalTeamStats[teamName] = &TeamStats{}
			}

			endOfMatchSummation(game.TotalTeamStats[teamName], &team)
		}

		//add to round master stats
		log.Debug(game.Rounds[i].RoundNum)
		for steam, player := range (*game.Rounds[i]).PlayerStats {
			if game.TotalPlayerStats[steam] == nil {
				game.TotalPlayerStats[steam] = &PlayerStats{Name: player.Name, SteamID: player.SteamID, TeamClanName: player.TeamClanName}
			}
			game.TotalPlayerStats[steam].Rounds += 1
			game.TotalPlayerStats[steam].Side = 4

			endOfMatchSummation(game.TotalPlayerStats[steam], &player)

			if player.IsBot {
				game.TotalPlayerStats[steam].IsBot = true
			}

			if player.RF == 1 {
				game.Rounds[i].WinTeamDmg += player.Damage
			}

			// TODO: have T/CT side type within PlayerStats
			if player.Side == int(common.TeamTerrorists) {
				game.TotalPlayerStats[steam].WinPointsNormalizer += game.Rounds[i].InitTerroristCount
				game.TotalPlayerStats[steam].TImpactPoints += player.ImpactPoints
				game.TotalPlayerStats[steam].TWinPoints += player.WinPoints
				game.TotalPlayerStats[steam].TOK += player.Ok
				game.TotalPlayerStats[steam].TOL += player.Ol
				game.TotalPlayerStats[steam].TKills += player.Kills
				game.TotalPlayerStats[steam].TDeaths += player.Deaths
				game.TotalPlayerStats[steam].TKASTRounds += player.KastRounds
				game.TotalPlayerStats[steam].TDamage += player.Damage
				game.TotalPlayerStats[steam].TADP += player.DeathPlacement
				//Game.TotalPlayerStats[steam].tTeamsWinPoints +=
				game.TotalPlayerStats[steam].TWinPointsNormalizer += game.Rounds[i].InitTerroristCount
				game.TotalPlayerStats[steam].TRounds += 1
				game.TotalPlayerStats[steam].TRF += player.RF
				game.TotalPlayerStats[steam].LurkRounds += player.LurkRounds
				if player.LurkRounds != 0 {
					game.TotalPlayerStats[steam].Wlp += player.WinPoints
				}

				game.Rounds[i].TeamStats[player.TeamClanName].TWinPoints += player.WinPoints
				game.Rounds[i].TeamStats[player.TeamClanName].TImpactPoints += player.ImpactPoints
			} else if player.Side == int(common.TeamCounterTerrorists) {
				game.TotalPlayerStats[steam].WinPointsNormalizer += game.Rounds[i].InitCTerroristCount
				game.TotalPlayerStats[steam].CtImpactPoints += player.ImpactPoints
				game.TotalPlayerStats[steam].CtWinPoints += player.WinPoints
				game.TotalPlayerStats[steam].CtOK += player.Ok
				game.TotalPlayerStats[steam].CtOL += player.Ol
				game.TotalPlayerStats[steam].CtKills += player.Kills
				game.TotalPlayerStats[steam].CtDeaths += player.Deaths
				game.TotalPlayerStats[steam].CtKASTRounds += player.KastRounds
				game.TotalPlayerStats[steam].CtDamage += player.Damage
				game.TotalPlayerStats[steam].CtADP += player.DeathPlacement
				//Game.TotalPlayerStats[steam].tTeamsWinPoints +=
				game.TotalPlayerStats[steam].CtWinPointsNormalizer += game.Rounds[i].InitCTerroristCount
				game.TotalPlayerStats[steam].CtRounds += 1
				game.TotalPlayerStats[steam].CtAWP += player.CtAWP

				game.Rounds[i].TeamStats[player.TeamClanName].CtWinPoints += player.WinPoints
				game.Rounds[i].TeamStats[player.TeamClanName].CtImpactPoints += player.ImpactPoints
			}

			game.Rounds[i].TeamStats[player.TeamClanName].WinPoints += player.WinPoints
			game.Rounds[i].TeamStats[player.TeamClanName].ImpactPoints += player.ImpactPoints

		}
		for steam, player := range (*game.Rounds[i]).PlayerStats {
			game.TotalPlayerStats[steam].TeamsWinPoints += game.Rounds[i].TeamStats[player.TeamClanName].WinPoints
			game.TotalPlayerStats[steam].TTeamsWinPoints += game.Rounds[i].TeamStats[player.TeamClanName].TWinPoints
			game.TotalPlayerStats[steam].CtTeamsWinPoints += game.Rounds[i].TeamStats[player.TeamClanName].CtWinPoints

			//give players rws
			if player.RF != 0 {
				if game.Rounds[i].EndDueToBombEvent {
					player.Rws = 70 * (float64(player.Damage) / float64(game.Rounds[i].WinTeamDmg))
					steamId64 := player.SteamID
					if player.Side == 2 && game.Rounds[i].Planter == steamId64 {
						player.Rws += 30
					} else if player.Side == 3 && game.Rounds[i].Defuser == steamId64 {
						player.Rws += 30
					}
				} else { //round ended due to damage/time
					player.Rws = 100 * (float64(player.Damage) / float64(game.Rounds[i].WinTeamDmg))
				}
				if math.IsNaN(player.Rws) {
					player.Rws = 0.0
				}
				game.TotalPlayerStats[steam].Rws += player.Rws
			}
		}
	}

	for _, player := range game.TotalPlayerStats {
		game.TotalTeamStats[player.TeamClanName].Util += player.SmokeThrown + player.FlashThrown + player.NadesThrown + player.FiresThrown
		game.TotalTeamStats[player.TeamClanName].Ud += player.UtilDmg
		game.TotalTeamStats[player.TeamClanName].Ef += player.Ef
		game.TotalTeamStats[player.TeamClanName].Fass += player.FAss
		game.TotalTeamStats[player.TeamClanName].Traded += player.Traded
		game.TotalTeamStats[player.TeamClanName].Deaths += int(player.Deaths)
	}

	game.calculateDerivedFields()
	return
}

func endOfMatchSummation(currentStats, newStats interface{}) {
	if currentStats == nil || newStats == nil {
		return
	}

	vCurrentStats := reflect.ValueOf(currentStats)
	vNewStats := reflect.ValueOf(newStats)

	if vCurrentStats.Kind() != reflect.Ptr || vNewStats.Kind() != reflect.Ptr {
		return
	}

	vCurrentStats = vCurrentStats.Elem()
	vNewStats = vNewStats.Elem()
	t := vCurrentStats.Type()

	for i := 0; i < vCurrentStats.NumField(); i++ {
		field := t.Field(i)

		if field.Tag.Get("end_of_match_sum") != "sum" {
			continue
		}

		fCurrentStats := vCurrentStats.Field(i)
		fNewStats := vNewStats.Field(i)

		if !fCurrentStats.CanSet() {
			continue
		}

		switch fCurrentStats.Kind() {
		case reflect.Int, reflect.Uint8:
			fCurrentStats.SetInt(fCurrentStats.Int() + fNewStats.Int())
			break
		case reflect.Float64:
			fCurrentStats.SetFloat(fCurrentStats.Float() + fNewStats.Float())
			break
		default:
			panic("Unknown stat point")
		}
	}
}

func (game *Game) calculateDerivedFields() {

	impactRoundAvg := 0.0
	killRoundAvg := 0.0
	deathRoundAvg := 0.0
	kastRoundAvg := 0.0
	adrAvg := 0.0
	roundNormalizer := 0

	tImpactRoundAvg := 0.0
	tKillRoundAvg := 0.0
	tDeathRoundAvg := 0.0
	tKastRoundAvg := 0.0
	tAdrAvg := 0.0
	tRoundNormalizer := 0

	ctImpactRoundAvg := 0.0
	ctKillRoundAvg := 0.0
	ctDeathRoundAvg := 0.0
	ctKastRoundAvg := 0.0
	ctAdrAvg := 0.0
	ctRoundNormalizer := 0

	//check our shit
	for _, player := range game.TotalPlayerStats {

		player.Atd = player.TicksAlive / player.Rounds / game.TickRate
		player.DeathPlacement = player.DeathPlacement / float64(player.Deaths)
		player.Kast = player.KastRounds / float64(player.Rounds)
		player.KillPointAvg = player.KillPoints / float64(player.Kills)
		if player.Kills == 0 {
			player.KillPointAvg = 0
		}
		player.Iiwr = player.WinPoints / player.ImpactPoints
		player.Adr = float64(player.Damage) / float64(player.Rounds)
		player.DrDiff = player.Adr - (float64(player.DamageTaken) / float64(player.Rounds))
		player.Tr = float64(player.Traded) / float64(player.Deaths)
		player.KR = float64(player.Kills) / float64(player.Rounds)
		player.UtilThrown = player.SmokeThrown + player.FlashThrown + player.NadesThrown + player.FiresThrown
		player.Rws = player.Rws / float64(player.Rounds)

		if player.CtRounds > 0 {
			player.CtADR = float64(player.CtDamage) / float64(player.CtRounds)
			player.CtKAST = player.CtKASTRounds / float64(player.CtRounds)
			player.CtADP = player.CtADP / float64(player.CtDeaths)
			if player.CtDeaths == 0 {
				player.CtADP = 0
			}
			ctImpactRoundAvg += player.CtImpactPoints
			ctKillRoundAvg += float64(player.CtKills)
			ctDeathRoundAvg += float64(player.CtDeaths)
			ctKastRoundAvg += player.CtKASTRounds
			ctAdrAvg += float64(player.CtDamage)
			ctRoundNormalizer += player.CtRounds
		}

		if player.TRounds > 0 {
			player.TADR = float64(player.TDamage) / float64(player.TRounds)
			player.TKAST = player.TKASTRounds / float64(player.TRounds)
			player.TADP = player.TADP / float64(player.TDeaths)
			if player.TDeaths == 0 {
				player.TADP = 0
			}
			tImpactRoundAvg += player.TImpactPoints
			tKillRoundAvg += float64(player.TKills)
			tDeathRoundAvg += float64(player.TDeaths)
			tKastRoundAvg += player.TKASTRounds
			tAdrAvg += float64(player.TDamage)
			tRoundNormalizer += player.TRounds
		}

		if math.IsNaN(player.Rws) {
			player.Rws = 0.0
		}
		if player.ImpactPoints == 0 {
			player.Iiwr = 0
		}
		if player.Deaths == 0 {
			player.DeathPlacement = 0
			player.Tr = .50
		}

		roundNormalizer += player.Rounds
		impactRoundAvg += player.ImpactPoints
		killRoundAvg += float64(player.Kills)
		deathRoundAvg += float64(player.Deaths)
		kastRoundAvg += player.KastRounds
		adrAvg += float64(player.Damage)
	}

	impactRoundAvg /= float64(roundNormalizer)
	killRoundAvg /= float64(roundNormalizer)
	deathRoundAvg /= float64(roundNormalizer)
	kastRoundAvg /= float64(roundNormalizer)
	adrAvg /= float64(roundNormalizer)

	tImpactRoundAvg /= float64(tRoundNormalizer)
	tKillRoundAvg /= float64(tRoundNormalizer)
	tDeathRoundAvg /= float64(tRoundNormalizer)
	tKastRoundAvg /= float64(tRoundNormalizer)
	tAdrAvg /= float64(tRoundNormalizer)

	ctImpactRoundAvg /= float64(ctRoundNormalizer)
	ctKillRoundAvg /= float64(ctRoundNormalizer)
	ctDeathRoundAvg /= float64(ctRoundNormalizer)
	ctKastRoundAvg /= float64(ctRoundNormalizer)
	ctAdrAvg /= float64(ctRoundNormalizer)

	for _, player := range game.TotalPlayerStats {
		openingFactor := (float64(player.Ok-player.Ol) / 13.0) + 1 //move from 13 to (rounds / 5)
		playerIPR := player.ImpactPoints / float64(player.Rounds)
		playerWPR := player.WinPoints / float64(player.Rounds)

		if player.TeamsWinPoints != 0 {
			player.ImpactRating = (0.1 * float64(openingFactor)) + (0.6 * (playerIPR / impactRoundAvg)) + (0.3 * (playerWPR / (player.TeamsWinPoints / float64(player.WinPointsNormalizer))))
		} else {
			log.Debug("UH 16-0?")
			player.ImpactRating = (0.1 * float64(openingFactor)) + (0.6 * (playerIPR / impactRoundAvg))
		}
		playerDR := float64(player.Deaths) / float64(player.Rounds)
		playerRatingDeathComponent := 0.07 * (deathRoundAvg / playerDR)
		if player.Deaths == 0 || playerRatingDeathComponent > 0.21 {
			playerRatingDeathComponent = 0.21
		}
		player.Rating = (0.3 * player.ImpactRating) + (0.35 * (player.KR / killRoundAvg)) + playerRatingDeathComponent + (0.08 * (player.Kast / kastRoundAvg)) + (0.2 * (player.Adr / adrAvg))

		//ctRating
		if player.CtRounds > 0 {
			openingFactor = (float64(player.CtOK-player.CtOL) / 13.0) + 1
			playerIPR = player.CtImpactPoints / float64(player.CtRounds)
			playerWPR = player.CtWinPoints / float64(player.CtRounds)

			if player.CtTeamsWinPoints != 0 {
				player.CtImpactRating = (0.1 * float64(openingFactor)) + (0.6 * (playerIPR / ctImpactRoundAvg)) + (0.3 * (playerWPR / (player.CtTeamsWinPoints / float64(player.CtWinPointsNormalizer))))
			} else {
				log.Debug("UH 16-0?")
				player.CtImpactRating = (0.1 * float64(openingFactor)) + (0.6 * (playerIPR / ctImpactRoundAvg))
			}
			playerDR = float64(player.CtDeaths) / float64(player.CtRounds)
			playerRatingDeathComponent = 0.07 * (ctDeathRoundAvg / playerDR)
			if player.CtDeaths == 0 || playerRatingDeathComponent > 0.21 {
				playerRatingDeathComponent = 0.21
			}
			player.CtRating = (0.3 * player.CtImpactRating) + (0.35 * ((float64(player.CtKills) / float64(player.CtRounds)) / ctKillRoundAvg)) + playerRatingDeathComponent + (0.08 * (player.CtKAST / ctKastRoundAvg)) + (0.2 * (player.CtADR / ctAdrAvg))
		}

		//tRating
		if player.TRounds > 0 {
			openingFactor = (float64(player.TOK-player.TOL) / 13.0) + 1
			playerIPR = player.TImpactPoints / float64(player.TRounds)
			playerWPR = player.TWinPoints / float64(player.TRounds)

			if player.TTeamsWinPoints != 0 {
				player.TImpactRating = (0.1 * float64(openingFactor)) + (0.6 * (playerIPR / tImpactRoundAvg)) + (0.3 * (playerWPR / (player.TTeamsWinPoints / float64(player.TWinPointsNormalizer))))
			} else {
				log.Debug("UH 16-0?")
				player.TImpactRating = (0.1 * float64(openingFactor)) + (0.6 * (playerIPR / tImpactRoundAvg))
			}
			playerDR = float64(player.TDeaths) / float64(player.TRounds)
			playerRatingDeathComponent = 0.07 * (tDeathRoundAvg / playerDR)
			if player.TDeaths == 0 || playerRatingDeathComponent > 0.21 {
				playerRatingDeathComponent = 0.21
			}
			player.TRating = (0.3 * player.TImpactRating) + (0.35 * ((float64(player.TKills) / float64(player.TRounds)) / tKillRoundAvg)) + playerRatingDeathComponent + (0.08 * (player.TKAST / tKastRoundAvg)) + (0.2 * (player.TADR / tAdrAvg))
		}

		log.Debug("openingFactor", 0.1*float64(openingFactor))
		log.Debug("playerIPR", 0.6*(playerIPR/impactRoundAvg))
		log.Debug("playerWPR", 0.3*(playerWPR/(player.TeamsWinPoints/float64(player.WinPointsNormalizer))))
		log.Debug("player.teamsWinPoints", player.TeamsWinPoints)
		log.Debug("player.winPointsNormalizer", player.WinPointsNormalizer)

		log.Debug("%+v\n\n", player)
	}
	log.Debug("impactRoundAvg", impactRoundAvg)
	log.Debug("killRoundAvg", killRoundAvg)
	log.Debug("deathRoundAvg", deathRoundAvg)
	log.Debug("kastRoundAvg", kastRoundAvg)
	log.Debug("adrAvg", adrAvg)

	game.calculateSidedStats()
	return
}

func (game *Game) calculateSidedStats() {

	game.CtPlayerStats = make(map[uint64]*PlayerStats)
	game.TPlayerStats = make(map[uint64]*PlayerStats)

	for i := len(game.Rounds) - 1; i >= 0; i-- {
		//add to round master stats
		for steam, player := range (*game.Rounds[i]).PlayerStats {
			//sidedStats := make(map[uint64]*playerStats)
			sidedStats := game.CtPlayerStats
			if player.Side == 2 {
				sidedStats = game.TPlayerStats
			}
			if sidedStats[steam] == nil {
				sidedStats[steam] = &PlayerStats{Name: player.Name, SteamID: player.SteamID, TeamClanName: player.TeamClanName}
			}
			sidedStats[steam].Rounds += 1
			sidedStats[steam].Kills += player.Kills
			sidedStats[steam].Assists += player.Assists
			sidedStats[steam].Deaths += player.Deaths
			sidedStats[steam].Damage += player.Damage
			sidedStats[steam].TicksAlive += player.TicksAlive
			sidedStats[steam].DeathPlacement += player.DeathPlacement
			sidedStats[steam].Trades += player.Trades
			sidedStats[steam].Traded += player.Traded
			sidedStats[steam].Ok += player.Ok
			sidedStats[steam].Ol += player.Ol
			sidedStats[steam].KillPoints += player.KillPoints
			sidedStats[steam].Cl_1 += player.Cl_1
			sidedStats[steam].Cl_2 += player.Cl_2
			sidedStats[steam].Cl_3 += player.Cl_3
			sidedStats[steam].Cl_4 += player.Cl_4
			sidedStats[steam].Cl_5 += player.Cl_5
			sidedStats[steam].TwoK += player.TwoK
			sidedStats[steam].ThreeK += player.ThreeK
			sidedStats[steam].FourK += player.FourK
			sidedStats[steam].FiveK += player.FiveK
			sidedStats[steam].NadeDmg += player.NadeDmg
			sidedStats[steam].InfernoDmg += player.InfernoDmg
			sidedStats[steam].UtilDmg += player.UtilDmg
			sidedStats[steam].Ef += player.Ef
			sidedStats[steam].FAss += player.FAss
			sidedStats[steam].EnemyFlashTime += player.EnemyFlashTime
			sidedStats[steam].Hs += player.Hs
			sidedStats[steam].KastRounds += player.KastRounds
			sidedStats[steam].Saves += player.Saves
			sidedStats[steam].Entries += player.Entries
			sidedStats[steam].ImpactPoints += player.ImpactPoints
			sidedStats[steam].WinPoints += player.WinPoints
			sidedStats[steam].AwpKills += player.AwpKills
			sidedStats[steam].RF += player.RF
			sidedStats[steam].RA += player.RA
			sidedStats[steam].NadesThrown += player.NadesThrown
			sidedStats[steam].SmokeThrown += player.SmokeThrown
			sidedStats[steam].FlashThrown += player.FlashThrown
			sidedStats[steam].FiresThrown += player.FiresThrown
			sidedStats[steam].DamageTaken += player.DamageTaken
			sidedStats[steam].SuppDamage += player.SuppDamage
			sidedStats[steam].SuppRounds += player.SuppRounds
			sidedStats[steam].Rwk += player.Rwk
			sidedStats[steam].Mip += player.Mip
			sidedStats[steam].Eac += player.Eac
			sidedStats[steam].Side = player.Side

			if player.IsBot {
				sidedStats[steam].IsBot = true
			}

			sidedStats[steam].LurkRounds += player.LurkRounds
			if player.LurkRounds != 0 {
				sidedStats[steam].Wlp += player.WinPoints
			}

			if math.IsNaN(player.Rws) {
				player.Rws = 0.0
			}
			sidedStats[steam].Rws += player.Rws

			if player.Side == 2 {
				sidedStats[steam].Rating = game.TotalPlayerStats[steam].TRating
				sidedStats[steam].ImpactRating = game.TotalPlayerStats[steam].TImpactRating
			} else {
				sidedStats[steam].Rating = game.TotalPlayerStats[steam].CtRating
				sidedStats[steam].ImpactRating = game.TotalPlayerStats[steam].CtImpactRating
			}

		}
	}

	for _, player := range game.CtPlayerStats {
		player.Atd = player.TicksAlive / player.Rounds / game.TickRate
		player.DeathPlacement = player.DeathPlacement / float64(player.Deaths)
		player.Kast = player.KastRounds / float64(player.Rounds)
		player.KillPointAvg = player.KillPoints / float64(player.Kills)
		if player.Kills == 0 {
			player.KillPointAvg = 0
		}
		player.Iiwr = player.WinPoints / player.ImpactPoints
		player.Adr = float64(player.Damage) / float64(player.Rounds)
		player.DrDiff = player.Adr - (float64(player.DamageTaken) / float64(player.Rounds))
		player.Tr = float64(player.Traded) / float64(player.Deaths)
		player.KR = float64(player.Kills) / float64(player.Rounds)
		player.UtilThrown = player.SmokeThrown + player.FlashThrown + player.NadesThrown + player.FiresThrown
		player.Rws = player.Rws / float64(player.Rounds)
		if math.IsNaN(player.Rws) {
			player.Rws = 0.0
		}
		if player.ImpactPoints == 0 {
			player.Iiwr = 0
		}
		if player.Deaths == 0 {
			player.DeathPlacement = 0
			player.Tr = .50
		}
		if player.TDeaths == 0 {
			player.TADP = 0
		}
		if player.CtDeaths == 0 {
			player.CtADP = 0
		}
	}
	for _, player := range game.TPlayerStats {
		player.Atd = player.TicksAlive / player.Rounds / game.TickRate
		player.DeathPlacement = player.DeathPlacement / float64(player.Deaths)
		player.Kast = player.KastRounds / float64(player.Rounds)
		player.KillPointAvg = player.KillPoints / float64(player.Kills)
		if player.Kills == 0 {
			player.KillPointAvg = 0
		}
		player.Iiwr = player.WinPoints / player.ImpactPoints
		player.Adr = float64(player.Damage) / float64(player.Rounds)
		player.DrDiff = player.Adr - (float64(player.DamageTaken) / float64(player.Rounds))
		player.Tr = float64(player.Traded) / float64(player.Deaths)
		player.KR = float64(player.Kills) / float64(player.Rounds)
		player.UtilThrown = player.SmokeThrown + player.FlashThrown + player.NadesThrown + player.FiresThrown
		player.Rws = player.Rws / float64(player.Rounds)
		if math.IsNaN(player.Rws) {
			player.Rws = 0.0
		}
		if player.ImpactPoints == 0 {
			player.Iiwr = 0
		}
		if player.Deaths == 0 {
			player.DeathPlacement = 0
			player.Tr = .50
		}
		if player.TDeaths == 0 {
			player.TADP = 0
		}
		if player.CtDeaths == 0 {
			player.CtADP = 0
		}
	}

	return
}
