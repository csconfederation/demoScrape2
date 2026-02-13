package demoscrape2

import (
	"io"
	"math"
	"strconv"

	dem "github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs"
	common "github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs/common"
	events "github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs/events"
	"github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs/msg"
	log "github.com/sirupsen/logrus"
)

//TODO
//"Catch up on the score" - dont remember what this is lol

//BUG fix
//MAKE ROUNDENDOFFICIAL Redundant (may have done, make sure it passes validation tho)
//add verification for missed event triggers if someone DCs/Cs after the redundant event that is stale
//csgo bots all have same steamID, need to use something else just in case for bots
//MM bug?

//add support for esea games? need to change validation and how we determine what round it is (gamestats rounds doesnt work)

//FUNCTIONAL CHANGES
//add verification for if a round event has triggered so far in the round (avoid double roundEnds)
//check for Game start without pistol (if we have bad demo)
//Add backend support
//Add anchor stuff
//Add team economy round stats (ecos, forces, etc)
//Add various nil checking

//CLEAN CODE
//TODO: create a outputPlayer function to clean up output.go
//TODO: convert rating calculations to a function
//TODO: actually use killValues lmao

const DEBUG = false

//const suppressNormalOutput = false

// globals
const printChatLog = true
const printDebugLog = true
const FORCE_NEW_STATS_UPLOAD = false
const ENABLE_WPA_DATA_OUTPUT = false
const BACKEND_PUSHING = true
const MR = 12

const tradeCutoff = 4 // in seconds
var multiKillBonus = [...]float64{0, 0, 0.3, 0.7, 1.2, 2}
var clutchBonus = [...]float64{0, 0.2, 0.6, 1.2, 2, 3}
var killValues = map[string]float64{
	"attacking":     1.2, //base values
	"defending":     1.0,
	"bombDefense":   1.0,
	"retake":        1.2,
	"chase":         0.8,
	"exit":          0.6,
	"t_consolation": 0.5,
	"gravy":         0.6,
	"punish":        0.8,
	"entry":         0.8, //multipliers
	"t_opener":      0.3,
	"ct_opener":     0.5,
	"trade":         0.3,
	"flashAssist":   0.2,
	"assist":        0.15,
}

func InitGameObject() *Game {
	g := Game{
		ReconnectedPlayers: make(map[uint64]bool),
	}
	g.Rounds = make([]*round, 0)
	g.PotentialRound = &round{}

	g.Flags.HasGameStarted = false
	g.Flags.IsGameLive = false
	g.Flags.IsGameOver = false
	g.Flags.InRound = false
	g.Flags.PrePlant = true
	g.Flags.PostPlant = false
	g.Flags.PostWinCon = false
	//these three vars to check if we have a complete round
	g.Flags.RoundIntegrityStart = -1
	g.Flags.RoundIntegrityEnd = -1
	g.Flags.RoundIntegrityEndOfficial = -1

	return &g
}

func ProcessDemo(demo io.ReadCloser) (*Game, error) {

	game := InitGameObject()

	p := dem.NewParser(demo)
	defer p.Close()

	//set tick rate
	game.TickRate = 64
	log.Debug("Tick rate is", game.TickRate)

	//---------------FUNCTIONS---------------

	initGameStart := func() {
		game.Flags.HasGameStarted = true
		game.Flags.IsGameLive = true
		log.Debug("GAME HAS STARTED!!!")

		// In case the tickRate is 0 we want to re-set it based on the tickInterval now that the Game has hasGameStarted
		if game.TickRate == 0 {
			game.TickRate = int(math.Round(p.TickRate()))
		}

		game.Teams = make(map[string]*team)

		teamTemp := p.GameState().TeamTerrorists()
		game.Teams[validateTeamName(game, teamTemp.ClanName(), teamTemp.Team())] = &team{Name: validateTeamName(game, teamTemp.ClanName(), teamTemp.Team())}
		teamTemp = p.GameState().TeamCounterTerrorists()
		game.Teams[validateTeamName(game, teamTemp.ClanName(), teamTemp.Team())] = &team{Name: validateTeamName(game, teamTemp.ClanName(), teamTemp.Team())}

		//only handling normal length matches
		game.RoundsToWin = MR + 1
		game.Result = ""

	}

	//reset various flags
	resetRoundFlags := func() {
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
		game.Flags.DidScoreUpdateFire = false
		game.Flags.RoundStartedAt = 0
		game.Flags.HaveInitRound = false
	}

	initTeamPlayer := func(team *common.TeamState, currRoundObj *round) {
		for _, teamMember := range getTeamMembers(team, game, p) {
			player := &playerStats{Name: teamMember.Name, SteamID: strconv.FormatUint(teamMember.SteamID64, 10), IsBot: teamMember.IsBot, Side: int(team.Team()), TeamENUM: team.ID(), TeamClanName: validateTeamName(game, team.ClanName(), team.Team()), Health: 100, TradeList: make(map[uint64]int), DamageList: make(map[uint64]int)}
			currRoundObj.PlayerStats[teamMember.SteamID64] = player
		}
	}

	initRound := func() {
		// Reset the connectedAfterRoundStart
		game.ConnectedAfterRoundStart = make(map[uint64]bool)

		game.Flags.RoundIntegrityStart = p.GameState().TotalRoundsPlayed() + 1
		log.Debug("We are starting round", game.Flags.RoundIntegrityStart)

		newRound := &round{RoundNum: int8(game.Flags.RoundIntegrityStart), StartingTick: p.GameState().IngameTick()}
		newRound.PlayerStats = make(map[uint64]*playerStats)
		newRound.TeamStats = make(map[string]*teamStats)

		//set players in playerStats for the round
		terrorists := p.GameState().TeamTerrorists()
		counterTerrorists := p.GameState().TeamCounterTerrorists()

		initTeamPlayer(terrorists, newRound)
		initTeamPlayer(counterTerrorists, newRound)

		//set teams in teamStats for the round
		newRound.TeamStats[validateTeamName(game, p.GameState().TeamTerrorists().ClanName(), p.GameState().TeamTerrorists().Team())] = &teamStats{TR: 1}
		newRound.TeamStats[validateTeamName(game, p.GameState().TeamCounterTerrorists().ClanName(), p.GameState().TeamCounterTerrorists().Team())] = &teamStats{CtR: 1}

		//create empty WPAlog
		newRound.WPAlog = make([]*wpalog, 0)

		// Reset round
		game.PotentialRound = newRound

		//track the number of people alive for clutch checking and record keeping
		game.Flags.TAlive = len(getTeamMembers(terrorists, game, p))
		game.Flags.CtAlive = len(getTeamMembers(counterTerrorists, game, p))
		game.PotentialRound.InitTerroristCount = game.Flags.TAlive
		game.PotentialRound.InitCTerroristCount = game.Flags.CtAlive

		resetRoundFlags()
	}

	processRoundOnWinCon := func(winnerClanName string) {
		game.Flags.RoundIntegrityEnd = p.GameState().TotalRoundsPlayed()
		log.Debug("We are processing round win con stuff", game.Flags.RoundIntegrityEnd)

		game.TotalRounds = game.Flags.RoundIntegrityEnd

		game.Flags.PrePlant = false
		game.Flags.PostPlant = false
		game.Flags.PostWinCon = true

		//set winner
		game.PotentialRound.WinnerClanName = winnerClanName
		//log.Debug("We think this team won", winnerClanName)
		if !game.PotentialRound.KnifeRound && !game.Flags.DidScoreUpdateFire {
			game.Teams[game.PotentialRound.WinnerClanName].Score += 1
		}
		//go through and set our WPAlog output to the winner
		for _, log := range game.PotentialRound.WPAlog {
			log.Winner = game.PotentialRound.WinnerENUM - 2
		}
	}

	processRoundFinal := func(lastRound bool) {
		game.Flags.InRound = false
		game.PotentialRound.EndingTick = p.GameState().IngameTick()
		game.Flags.RoundIntegrityEndOfficial = p.GameState().TotalRoundsPlayed()

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
				steamId64, _ := strconv.ParseUint(player.SteamID, 10, 64)
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
						steamId64, _ := strconv.ParseUint(player.SteamID, 10, 64)
						susLurker = steamId64
					}
				}
			}
			for _, player := range game.PotentialRound.PlayerStats {
				if player.Side == 2 {
					steamId64, _ := strconv.ParseUint(player.SteamID, 10, 64)
					if player.LurkerBlips == susLurkBlips && steamId64 != susLurker {
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
		if lastRound {
			//game.Flags.RoundIntegrityEndOfficial += 1
			game.TotalRounds = game.Flags.RoundIntegrityEndOfficial
			game.Flags.IsGameLive = false
			game.Result = "Ended"
		}

		//endRound function functionality

	}

	//-------------ALL OUR EVENTS---------------------

	p.RegisterNetMessageHandler(func(m *msg.CDemoFileHeader) {
		game.MapName = m.GetMapName()
	})

	p.RegisterEventHandler(func(e events.PlayerInfo) {
		player := p.GameState().Participants().AllByUserID()[e.Index]
		if player != nil {
			// if game.potentialRound.playerStats[player.SteamID64] == nil {
			// 	game.potentialRound.playerStats[player.SteamID64] = &playerStats{name: player.Name, steamID: player.SteamID64, isBot: player.IsBot, side: int(player.Team), teamENUM: int(player.Team), teamClanName: validateTeamName(game, player.TeamState.ClanName(), player.TeamState.Team()), health: 100, tradeList: make(map[uint64]int), damageList: make(map[uint64]int)}
			// }
			game.ReconnectedPlayers[player.SteamID64] = true
			if game.Flags.InRound && game.Flags.IsGameLive {
				game.ConnectedAfterRoundStart[player.SteamID64] = true
			}
		}
	})

	p.RegisterEventHandler(func(e events.FrameDone) {
		//log.Debug("DIBES ", Game.flags.isGameLive)
		if game.Flags.RoundStartedAt > 0 && game.Flags.RoundStartedAt+(1*game.TickRate) > p.GameState().IngameTick() && !game.Flags.HaveInitRound {
			pistol := false

			//we are going to check to see if the first pistol is actually starting
			membersT := getTeamMembers(p.GameState().TeamTerrorists(), game, p)
			membersCT := getTeamMembers(p.GameState().TeamCounterTerrorists(), game, p)
			if len(membersT) != 0 && len(membersCT) != 0 {
				if membersT[0].Money()+membersT[0].MoneySpentThisRound() == 800 && membersCT[0].Money()+membersCT[0].MoneySpentThisRound() == 800 {
					//start the Game
					if !game.Flags.HasGameStarted {
						initGameStart()
					}

					//track the pistol
					pistol = true
				}
			}
			//log.Debug("Has the Game Started?", Game.flags.hasGameStarted)

			if game.Flags.IsGameLive {
				//init round stats
				initRound()
				game.Flags.HaveInitRound = true
				if pistol {
					for _, team := range game.PotentialRound.TeamStats {
						team.Pistols = 1
					}
				}

			}
		}

		//add to WPAlog
		if game.Flags.IsGameLive && game.Flags.InRound && !game.Flags.PostWinCon && ENABLE_WPA_DATA_OUTPUT {
			//hits every new frame (typically each 1-4 ticks)
			logSize := len(game.PotentialRound.WPAlog)
			clock := 0
			planted := 0
			if game.PotentialRound.Planter != 0 {
				planted = 1
				bombTime, boo := p.GameState().Rules().BombTime()
				bombClock := 0
				if boo != nil {
					bombClock = 40
				} else {
					bombClock = int(bombTime.Seconds())
				}
				clock = bombClock - ((p.GameState().IngameTick() - game.PotentialRound.BombStartTick) / game.TickRate)
			} else {
				roundTime, boo := p.GameState().Rules().RoundTime()
				if boo != nil {
					log.Debug("RUROO RAGGY")
				}
				clock = int(roundTime.Seconds()) - ((p.GameState().IngameTick() - game.PotentialRound.StartingTick) / game.TickRate)
			}

			if logSize == 0 || game.PotentialRound.WPAlog[logSize-1].Tick+(game.TickRate) < p.GameState().IngameTick() {
				newWPAentry := &wpalog{
					Round:               int(game.PotentialRound.RoundNum),
					Clock:               clock,
					Planted:             planted,
					Tick:                p.GameState().IngameTick(),
					CtAlive:             game.Flags.CtAlive,
					TAlive:              game.Flags.TAlive,
					CtEquipVal:          calculateTeamEquipmentValue(game, p.GameState().TeamCounterTerrorists(), p),
					TEquipVal:           calculateTeamEquipmentValue(game, p.GameState().TeamTerrorists(), p),
					CtFlashes:           calculateTeamEquipmentNum(game, p.GameState().TeamCounterTerrorists(), 15, p),
					CtSmokes:            calculateTeamEquipmentNum(game, p.GameState().TeamCounterTerrorists(), 16, p),
					CtMolys:             calculateTeamEquipmentNum(game, p.GameState().TeamCounterTerrorists(), 17, p),
					CtFrags:             calculateTeamEquipmentNum(game, p.GameState().TeamCounterTerrorists(), 14, p),
					TFlashes:            calculateTeamEquipmentNum(game, p.GameState().TeamTerrorists(), 15, p),
					TSmokes:             calculateTeamEquipmentNum(game, p.GameState().TeamTerrorists(), 16, p),
					TMolys:              calculateTeamEquipmentNum(game, p.GameState().TeamTerrorists(), 17, p),
					TFrags:              calculateTeamEquipmentNum(game, p.GameState().TeamTerrorists(), 14, p),
					ClosestCTDisttoBomb: closestCTDisttoBomb(game, p.GameState().TeamCounterTerrorists(), p.GameState().Bomb(), p),
					Kits:                numOfKits(game, p.GameState().TeamCounterTerrorists(), p),
					CtArmor:             playersWithArmor(game, p.GameState().TeamCounterTerrorists(), p),
					TArmor:              playersWithArmor(game, p.GameState().TeamTerrorists(), p),
				}
				game.PotentialRound.WPAlog = append(game.PotentialRound.WPAlog, newWPAentry)
			}
		}

		if game.Flags.IsGameLive && game.Flags.InRound && game.Flags.LastTickProcessed+(4*game.TickRate) < p.GameState().IngameTick() {
			game.Flags.LastTickProcessed = p.GameState().IngameTick()
			game.Flags.TicksProcessed += 1

			//this will be triggered every 4 seconds of in round time after the first 10 seconds

			//check for lurker
			if game.Flags.TAlive > 2 && !game.Flags.PostWinCon && p.GameState().IngameTick() > (18*game.TickRate)+game.PotentialRound.StartingTick {
				membersT := getTeamMembers(p.GameState().TeamTerrorists(), game, p)
				for _, terrorist := range membersT {
					if terrorist.IsAlive() {
						for _, teammate := range membersT {
							if terrorist.SteamID64 != teammate.SteamID64 && teammate.IsAlive() {
								dist := int(terrorist.Position().Distance(teammate.Position()))
								if dist < 500 {
									//invalidate the lurk blip b/c we have a close teammate
									game.PotentialRound.PlayerStats[terrorist.SteamID64].DistanceToTeammates = -999999
								}
								if game.PotentialRound.PlayerStats[terrorist.SteamID64] != nil {
									game.PotentialRound.PlayerStats[terrorist.SteamID64].DistanceToTeammates += dist
								} else {
									log.Debug("THIS IS WHERE WE BROKE_______________________________---------------------------------------------------")
								}
							}
						}
					}
				}
				var lurkerSteam uint64
				lurkerDist := 999999
				for _, terrorist := range membersT {
					if terrorist.IsAlive() {
						if game.PotentialRound.PlayerStats[terrorist.SteamID64] == nil {
							log.Debug(terrorist.Name)
						} else {
							dist := game.PotentialRound.PlayerStats[terrorist.SteamID64].DistanceToTeammates
							if dist < lurkerDist && dist > 0 {
								lurkerDist = dist
								lurkerSteam = terrorist.SteamID64
							}
						}

					}
				}
				if lurkerSteam != 0 {
					game.PotentialRound.PlayerStats[lurkerSteam].LurkerBlips += 1
				}
			}
		}
	})

	p.RegisterEventHandler(func(e events.RoundStart) {
		log.Debug("Round Start", p.GameState().TotalRoundsPlayed())

		game.Flags.RoundStartedAt = p.GameState().IngameTick()

	})

	p.RegisterEventHandler(func(e events.RoundFreezetimeEnd) {
		log.Debug("Round Freeze Time End\n")
		pistol := false

		//we are going to check to see if the first pistol is actually starting
		membersT := getTeamMembers(p.GameState().TeamTerrorists(), game, p)
		membersCT := getTeamMembers(p.GameState().TeamCounterTerrorists(), game, p)
		if len(membersT) != 0 && len(membersCT) != 0 {
			if membersT[0].Money()+membersT[0].MoneySpentThisRound() == 800 && membersCT[0].Money()+membersCT[0].MoneySpentThisRound() == 800 {
				//start the Game
				if !game.Flags.HasGameStarted {
					initGameStart()
				}

				//track the pistol
				pistol = true
			} else if membersT[0].Money()+membersT[0].MoneySpentThisRound() == 0 && membersCT[0].Money()+membersCT[0].MoneySpentThisRound() == 0 {
				game.PotentialRound.KnifeRound = true
				log.Debug("------------KNIFEROUND-----------")
				game.Flags.HasGameStarted = false
			}
		}
		log.Debug("Has the Game Started?", game.Flags.HasGameStarted)

		if game.Flags.IsGameLive {
			//init round stats
			game.Flags.InRound = true
			initRound()
			if pistol {
				for _, team := range game.PotentialRound.TeamStats {
					team.Pistols = 1
				}
			}

		}

	})

	p.RegisterEventHandler(func(e events.RoundEnd) {
		if game.Flags.IsGameLive {

			game.Flags.DidRoundEndFire = true

			log.Debug("Round", p.GameState().TotalRoundsPlayed(), "End", e.WinnerState.ClanName(), "won", "this determined from e.WinnerState.ClanName()")

			log.Debug("e.WinnerState.ID()", e.WinnerState.ID(), "and", "e.Winner", e.Winner, "and", "e.WinnerState.Team()", e.WinnerState.Team())

			validWinner := true
			if e.Winner < 2 {
				validWinner = false
				//and set the integrity flag to false

			} else if e.Winner == 2 {
				game.Flags.TMoney = true
			} else {
				//we need to check if the Game is over

			}

			//we want to actually process the round
			if game.Flags.IsGameLive && validWinner && game.Flags.RoundIntegrityStart == p.GameState().TotalRoundsPlayed() {
				game.PotentialRound.WinnerENUM = int(e.Winner)
				game.PotentialRound.RoundEndReason = roundEndReasons[int(e.Reason)]
				processRoundOnWinCon(validateTeamName(game, e.WinnerState.ClanName(), e.WinnerState.Team()))

				//check last round
				roundWinnerScore := game.Teams[validateTeamName(game, e.WinnerState.ClanName(), e.WinnerState.Team())].Score
				roundLoserScore := game.Teams[validateTeamName(game, e.LoserState.ClanName(), e.LoserState.Team())].Score
				log.Debug("winner Rounds", roundWinnerScore)
				log.Debug("loser Rounds", roundLoserScore)

				if game.RoundsToWin == 16 {
					//check for normal win
					if roundWinnerScore == 16 && roundLoserScore < 15 {
						//normal win
						game.WinnerClanName = game.PotentialRound.WinnerClanName
						processRoundFinal(true)
					} else if roundWinnerScore > 15 { //check for OT win
						overtime := ((roundWinnerScore+roundLoserScore)-30-1)/6 + 1
						//OT win
						if (roundWinnerScore-15-1)/3 == overtime {
							game.WinnerClanName = game.PotentialRound.WinnerClanName
							processRoundFinal(true)
						}
					}
				} else if game.RoundsToWin == 9 {
					//check for normal win
					if roundWinnerScore == 9 && roundLoserScore < 8 {
						//normal win
						game.WinnerClanName = game.PotentialRound.WinnerClanName
						processRoundFinal(true)
					} else if roundWinnerScore == 8 && roundLoserScore == 8 { //check for tie
						//tie
						game.WinnerClanName = game.PotentialRound.WinnerClanName
						processRoundFinal(true)
					}
				} else if game.RoundsToWin == 13 {
					//check for normal win
					if roundWinnerScore == 13 && roundLoserScore < 12 {
						//normal win
						game.WinnerClanName = game.PotentialRound.WinnerClanName
						processRoundFinal(true)
					} else if roundWinnerScore > 12 { //check for OT win
						overtime := ((roundWinnerScore+roundLoserScore)-24-1)/6 + 1
						//OT win
						if (roundWinnerScore-12-1)/3 == overtime {
							game.WinnerClanName = game.PotentialRound.WinnerClanName
							processRoundFinal(true)
						}
					}
				}
			}

			//check last round
			//or check overtime win

		}
	})

	//round end official doesnt fire on the last round
	p.RegisterEventHandler(func(e events.ScoreUpdated) {
		//CS2 swapped this event to be before RoundEnd
		//We have relied on this as a back up for failed RoundEnd events
		//may revisit depending on event reliability

		//added to ensure that a bad round that gets finished does not premuturely finish the game since we track score separately
		if game.Flags.IsGameLive {
			game.Flags.DidScoreUpdateFire = true
			//we take the existing preupdate score of the updating team score
			updatedTeam := game.Teams[validateTeamName(game, e.TeamState.ClanName(), e.TeamState.Team())]
			//and compare to the old score from scoreboard
			updatedTeam.Score = e.NewScore
			game.Teams[validateTeamName(game, e.TeamState.Opponent.ClanName(), e.TeamState.Opponent.Team())].Score = e.TeamState.Opponent.Score()
		}
	})

	//round end official doesnt fire on the last round
	p.RegisterEventHandler(func(e events.RoundEndOfficial) {

		log.Debug("Round End Official\n")

		if !game.Flags.DidRoundEndFire {
			game.Flags.RoundIntegrityEnd -= 1
		}

		log.Debug("isGameLive", game.Flags.IsGameLive, "roundIntegrityEnd", game.Flags.RoundIntegrityEnd, "pTotalRoundsPlayed", p.GameState().TotalRoundsPlayed())

		if game.Flags.IsGameLive && game.Flags.RoundIntegrityEnd == p.GameState().TotalRoundsPlayed() {
			processRoundFinal(false)
		}
	})

	// Register handler on kill events
	p.RegisterEventHandler(func(e events.Kill) {
		flashAssister := ""
		if game.Flags.IsGameLive && isDuringExpectedRound(game, p) {
			pS := game.PotentialRound.PlayerStats
			tick := p.GameState().IngameTick()

			killerExists := false
			victimExists := false
			assisterExists := false
			if e.Killer != nil && pS[e.Killer.SteamID64] != nil {
				killerExists = true
			}
			if e.Victim != nil && pS[e.Victim.SteamID64] != nil {
				victimExists = true
			}
			if e.Assister != nil && pS[e.Assister.SteamID64] != nil {
				assisterExists = true
			}
			if e.Weapon.Type == 404 && isRoundFinalInHalf(game.PotentialRound.RoundNum) {
				killerExists = false
				victimExists = false
				assisterExists = false
			}

			killValue := 1.0
			multiplier := 1.0
			traded := false
			assisted := false
			flashAssisted := false

			//death logic (traded here)
			if victimExists {
				pS[e.Victim.SteamID64].Deaths += 1
				pS[e.Victim.SteamID64].DeathTick = tick
				if e.Victim.Team == 2 {
					game.Flags.TAlive -= 1
					pS[e.Victim.SteamID64].DeathPlacement = float64(game.PotentialRound.InitTerroristCount - game.Flags.TAlive)
					//pS[e.Victim.SteamID64].tADP = float64(Game.potentialRound.initTerroristCount - Game.flags.tAlive)
				} else if e.Victim.Team == 3 {
					game.Flags.CtAlive -= 1
					pS[e.Victim.SteamID64].DeathPlacement = float64(game.PotentialRound.InitCTerroristCount - game.Flags.CtAlive)
					//pS[e.Victim.SteamID64].ctADP = float64(Game.potentialRound.initCTerroristCount - Game.flags.ctAlive)
				} else {
					//else log an error
				}

				//do 4v5 calc
				if game.Flags.OpeningKill && game.PotentialRound.InitCTerroristCount+game.PotentialRound.InitTerroristCount == 10 {
					//the 10th player died
					_4v5Team := pS[e.Victim.SteamID64].TeamClanName
					game.PotentialRound.TeamStats[_4v5Team].FourVFiveS = 1
					for teamName, team := range game.PotentialRound.TeamStats {
						if teamName != _4v5Team {
							team.FiveVFourS = 1
						}
					}
				}

				//add support damage
				for suppSteam, suppDMG := range pS[e.Victim.SteamID64].DamageList {
					if killerExists && suppSteam != e.Killer.SteamID64 {
						pS[suppSteam].SuppDamage += suppDMG
						if pS[suppSteam].SuppDamage > 60 {
							pS[suppSteam].SuppRounds = 1
						}
					} else if !killerExists {
						pS[suppSteam].SuppDamage += suppDMG
						if pS[suppSteam].SuppDamage > 60 {
							pS[suppSteam].SuppRounds = 1
						}
					}

				}

				//check clutch start

				if !game.Flags.PostWinCon {
					if game.Flags.TAlive == 1 && game.Flags.TClutchVal == 0 {
						game.Flags.TClutchVal = game.Flags.CtAlive
						membersT := getTeamMembers(p.GameState().TeamTerrorists(), game, p)
						for _, terrorist := range membersT {
							if terrorist.IsAlive() && e.Victim.SteamID64 != terrorist.SteamID64 {
								game.Flags.TClutchSteam = terrorist.SteamID64
								log.Debug("Clutch opportunity:", terrorist.Name, game.Flags.TClutchVal)
							}
						}
					}
					if game.Flags.CtAlive == 1 && game.Flags.CtClutchVal == 0 {
						game.Flags.CtClutchVal = game.Flags.TAlive
						membersCT := getTeamMembers(p.GameState().TeamCounterTerrorists(), game, p)
						for _, counterTerrorist := range membersCT {
							if counterTerrorist.IsAlive() && e.Victim.SteamID64 != counterTerrorist.SteamID64 {
								game.Flags.CtClutchSteam = counterTerrorist.SteamID64
								log.Debug("Clutch opportunity:", counterTerrorist.Name, game.Flags.CtClutchVal)
							}
						}
					}
				}

				pS[e.Victim.SteamID64].TicksAlive = tick - game.PotentialRound.StartingTick
				for deadGuySteam, deadTick := range (*game.PotentialRound).PlayerStats[e.Victim.SteamID64].TradeList {
					if tick-deadTick < tradeCutoff*game.TickRate {
						pS[deadGuySteam].Traded = 1
						pS[deadGuySteam].Eac += 1
						pS[deadGuySteam].KastRounds = 1
					}
				}
			}

			//assist logic
			if assisterExists && victimExists && e.Assister.TeamState.ID() != e.Victim.TeamState.ID() {
				//this logic needs to be replaced -yeti does not remember why he wrote this
				pS[e.Assister.SteamID64].Assists += 1
				pS[e.Assister.SteamID64].Eac += 1
				pS[e.Assister.SteamID64].KastRounds = 1
				pS[e.Assister.SteamID64].SuppRounds = 1
				assisted = true
				if e.AssistedFlash {
					pS[e.Assister.SteamID64].FAss += 1
					flashAssisted = true
					flashAssister = e.Assister.Name
					log.Debug("VALVE FLASH ASSIST")
				} else if float64(p.GameState().IngameTick()) < pS[e.Victim.SteamID64].MostRecentFlashVal {
					//this will trigger if there is both a flash assist and a damage assist
					pS[pS[e.Victim.SteamID64].MostRecentFlasher].FAss += 1
					pS[pS[e.Victim.SteamID64].MostRecentFlasher].Eac += 1
					pS[pS[e.Victim.SteamID64].MostRecentFlasher].SuppRounds = 1
					flashAssisted = true
					flashAssister = pS[pS[e.Victim.SteamID64].MostRecentFlasher].Name
				}

			}

			//kill logic (trades here)
			if killerExists && victimExists && e.Killer.TeamState.ID() != e.Victim.TeamState.ID() {
				pS[e.Killer.SteamID64].Kills += 1
				pS[e.Killer.SteamID64].KastRounds = 1
				pS[e.Killer.SteamID64].Rwk = 1
				pS[e.Killer.SteamID64].TradeList[e.Victim.SteamID64] = tick
				if e.Weapon.Type == 309 {
					pS[e.Killer.SteamID64].AwpKills += 1
					if e.Killer.Team == 3 {
						pS[e.Killer.SteamID64].CtAWP += 1
					}
				}
				if e.IsHeadshot {
					pS[e.Killer.SteamID64].Hs += 1
				}
				for _, deadTick := range (*game.PotentialRound).PlayerStats[e.Victim.SteamID64].TradeList {
					if tick-deadTick < tradeCutoff*game.TickRate {
						pS[e.Killer.SteamID64].Trades += 1
						traded = true
						break
					}
				}

				killerTeam := e.Killer.Team
				if game.Flags.PrePlant {
					//normal base value
					if killerTeam == 2 {
						//taking site by T
						killValue = 1.2
					} else if killerTeam == 3 {
						//site Defense by CT
						killValue = 1
					}
				} else if game.Flags.PostPlant {
					//site D or retake
					if killerTeam == 2 {
						//site Defense by T
						killValue = 1
					} else if killerTeam == 3 {
						//retake
						killValue = 1.2
					}
				} else if game.Flags.PostWinCon {
					//exit or chase
					if game.PotentialRound.WinnerENUM == 2 { //Ts win
						if killerTeam == 2 { //chase
							killValue = 0.8
						}
						if killerTeam == 3 { //exit
							killValue = 0.6
						}
					} else if game.PotentialRound.WinnerENUM == 3 { //CTs win
						if killerTeam == 2 { //T kill in lost round
							killValue = 0.5
						}
						if killerTeam == 3 { //CT kill in won round
							if game.Flags.TMoney {
								killValue = 0.6
							} else {
								killValue = 0.8
							}
						}
					}
				}

				if game.Flags.OpeningKill {
					game.Flags.OpeningKill = false

					pS[e.Killer.SteamID64].Ok = 1
					pS[e.Victim.SteamID64].Ol = 1

					if killerTeam == 2 { //T entry/opener {
						if game.Flags.PrePlant {
							multiplier += 0.8
							pS[e.Killer.SteamID64].Entries = 1
						} else {
							multiplier += 0.3
						}
					} else if killerTeam == 3 { //CT opener
						multiplier += 0.5
					}

				} else if traded {
					multiplier += 0.3
				}

				if flashAssisted { //flash assisted kill
					multiplier += 0.2
				}
				if assisted { //assisted kill
					killValue -= 0.15
					pS[e.Assister.SteamID64].ImpactPoints += 0.15
				}

				killValue *= multiplier

				ecoRatio := float64(e.Victim.EquipmentValueCurrent()) / float64(e.Killer.EquipmentValueCurrent())
				ecoMod := 1.0
				if ecoRatio > 4 {
					ecoMod += 0.25
				} else if ecoRatio > 2 {
					ecoMod += 0.14
				} else if ecoRatio < 0.25 {
					ecoMod -= 0.25
				} else if ecoRatio < 0.5 {
					ecoMod -= 0.14
				}
				killValue *= ecoMod

				pS[e.Killer.SteamID64].KillPoints += killValue
			}

		}
		var hs string
		if e.IsHeadshot {
			hs = " (HS)"
		}
		var wallBang string
		if e.PenetratedObjects > 0 {
			wallBang = " (WB)"
		}
		log.Debug("%s <%v%s%s> %s at %d flash assist by %s\n", e.Killer, e.Weapon, hs, wallBang, e.Victim, p.GameState().IngameTick(), flashAssister)
	})

	p.RegisterEventHandler(func(e events.PlayerHurt) {
		//log.Debug("Player Hurt\n")
		if game.Flags.IsGameLive {
			var equipment common.EquipmentType
			if e.Weapon == nil {
				equipment = -999
			} else {
				equipment = e.Weapon.Type
			}
			validDmg := e.Player != nil && game.PotentialRound.PlayerStats[e.Player.SteamID64] != nil && (equipment != 404 || !isRoundFinalInHalf(game.PotentialRound.RoundNum))
			if validDmg {
				game.PotentialRound.PlayerStats[e.Player.SteamID64].DamageTaken += e.HealthDamageTaken
			} else if e.Player != nil && e.Player.IsConnected && !(equipment == 404 && isRoundFinalInHalf(game.PotentialRound.RoundNum)) {
				//blow up if we aren't in freeze time
				if !p.GameState().IsFreezetimePeriod() {
					panic("We have a connected player who is not nil but no playerstats!")
				}
			}
			if e.Player != nil && game.PotentialRound.PlayerStats[e.Player.SteamID64] != nil && e.Attacker != nil && e.Player.Team != e.Attacker.Team {
				game.PotentialRound.PlayerStats[e.Attacker.SteamID64].Damage += e.HealthDamageTaken

				//add to damage list for supp damage calc
				game.PotentialRound.PlayerStats[e.Player.SteamID64].DamageList[e.Attacker.SteamID64] += e.HealthDamageTaken

				if equipment >= 500 && equipment <= 506 {
					game.PotentialRound.PlayerStats[e.Attacker.SteamID64].UtilDmg += e.HealthDamageTaken
					if equipment == 506 {
						game.PotentialRound.PlayerStats[e.Attacker.SteamID64].NadeDmg += e.HealthDamageTaken
					}
					if equipment == 502 || equipment == 503 {
						game.PotentialRound.PlayerStats[e.Attacker.SteamID64].InfernoDmg += e.HealthDamageTaken
					}
				}
			}
		}
	})

	p.RegisterEventHandler(func(e events.PlayerFlashed) {
		//log.Debug("Player Flashed")
		if game.Flags.IsGameLive && e.Player != nil && e.Attacker != nil {
			tick := float64(p.GameState().IngameTick())
			blindTicks := e.FlashDuration().Seconds() * float64(game.TickRate)
			victim := e.Player
			flasher := e.Attacker
			if flasher.Team != victim.Team && blindTicks > float64(game.TickRate) && victim.IsAlive() && (float64(victim.FlashDuration) < (blindTicks/float64(game.TickRate) + 1)) {
				game.PotentialRound.PlayerStats[flasher.SteamID64].Ef += 1
				game.PotentialRound.PlayerStats[flasher.SteamID64].EnemyFlashTime += (blindTicks / float64(game.TickRate))
				if tick+blindTicks > game.PotentialRound.PlayerStats[victim.SteamID64].MostRecentFlashVal {
					game.PotentialRound.PlayerStats[victim.SteamID64].MostRecentFlashVal = tick + blindTicks
					game.PotentialRound.PlayerStats[victim.SteamID64].MostRecentFlasher = flasher.SteamID64
				}

			}
			if flasher.Name != "" {
				//debugMsg := fmt.Sprintf("%s flashed %s for %.2f at %d. He was %f blind.\n", flasher, victim, blindTicks/128, int(tick), victim.FlashDuration)
			}

		}
		//log.Debug("Player Flashed", blindTicks, e.Attacker)
	})

	p.RegisterEventHandler(func(e events.BombPlanted) {
		log.Debug("Bomb Planted\n")
		if game.Flags.IsGameLive && !game.Flags.PostWinCon {
			game.Flags.PrePlant = false
			game.Flags.PostPlant = true
			game.Flags.TMoney = true
			game.PotentialRound.Planter = e.BombEvent.Player.SteamID64
			game.PotentialRound.BombStartTick = p.GameState().IngameTick()
		}
	})

	p.RegisterEventHandler(func(e events.BombDefused) {
		log.Debug("Bomb Defused by", e.BombEvent.Player.Name)
		if game.Flags.IsGameLive && !game.Flags.PostWinCon {
			game.Flags.PrePlant = false
			game.Flags.PostPlant = false
			game.Flags.PostWinCon = true
			game.PotentialRound.EndDueToBombEvent = true
			game.PotentialRound.Defuser = e.Player.SteamID64
			game.PotentialRound.PlayerStats[e.BombEvent.Player.SteamID64].ImpactPoints += 0.5
		}
	})

	p.RegisterEventHandler(func(e events.BombExplode) {
		log.Debug("Bomb Exploded\n")
		if game.Flags.IsGameLive && !game.Flags.PostWinCon {
			game.Flags.PrePlant = false
			game.Flags.PostPlant = false
			game.Flags.PostWinCon = true
			game.PotentialRound.EndDueToBombEvent = true
			if game.PotentialRound.Planter != 0 {
				game.PotentialRound.PlayerStats[game.PotentialRound.Planter].ImpactPoints += 0.5
			}
		}
	})

	p.RegisterEventHandler(func(e events.GrenadeProjectileThrow) {
		//log.Debug("Grenade Thrown", e.Projectile.WeaponInstance.Type)
		if game.Flags.IsGameLive {
			if e.Projectile.WeaponInstance.Type == 506 {
				game.PotentialRound.PlayerStats[e.Projectile.Thrower.SteamID64].NadesThrown += 1
			} else if e.Projectile.WeaponInstance.Type == 505 {
				game.PotentialRound.PlayerStats[e.Projectile.Thrower.SteamID64].SmokeThrown += 1
			} else if e.Projectile.WeaponInstance.Type == 504 {
				game.PotentialRound.PlayerStats[e.Projectile.Thrower.SteamID64].FlashThrown += 1
			} else if e.Projectile.WeaponInstance.Type == 502 || e.Projectile.WeaponInstance.Type == 503 {
				game.PotentialRound.PlayerStats[e.Projectile.Thrower.SteamID64].FiresThrown += 1
			}

		}
	})

	p.RegisterEventHandler(func(e events.PlayerTeamChange) {
		log.Debug("Player Changed Team:", e.Player, e.OldTeam, e.NewTeam)

		if game.Flags.IsGameLive && game.Flags.InRound {
			if e.NewTeam > 1 {
				//we are joining an actual team
				if game.PotentialRound.PlayerStats[e.Player.SteamID64] == nil && e.Player.IsAlive() {
					//get team
					team := e.NewTeamState
					player := &playerStats{Name: e.Player.Name, SteamID: strconv.FormatUint(e.Player.SteamID64, 10), IsBot: e.Player.IsBot, Side: int(team.Team()), TeamENUM: team.ID(), TeamClanName: validateTeamName(game, team.ClanName(), team.Team()), Health: 100, TradeList: make(map[uint64]int), DamageList: make(map[uint64]int)}
					steamId64, _ := strconv.ParseUint(player.SteamID, 10, 64)
					game.PotentialRound.PlayerStats[steamId64] = player
				}
			}
		}
	})

	p.RegisterEventHandler(func(e events.PlayerDisconnected) {
		log.Debug("Player DC", e.Player)

		if game.ReconnectedPlayers[e.Player.SteamID64] {
			game.ReconnectedPlayers[e.Player.SteamID64] = false
		}

		//update alive players
		if game.Flags.IsGameLive {
			game.Flags.TAlive = 0
			game.Flags.CtAlive = 0

			membersT := getTeamMembers(p.GameState().TeamTerrorists(), game, p)
			for _, terrorist := range membersT {
				if terrorist.IsAlive() {
					game.Flags.TAlive += 1
				}
			}
			membersCT := getTeamMembers(p.GameState().TeamCounterTerrorists(), game, p)
			for _, counterTerrorist := range membersCT {
				if counterTerrorist.IsAlive() {
					game.Flags.CtAlive += 1
				}
			}
		}

	})

	p.RegisterEventHandler(func(e events.Footstep) {
		if game.Flags.IsGameLive {
			game.Flags.InRound = true
		}

	})

	// Parse to end
	err := p.ParseToEnd()

	endOfMatchProcessing(game)

	return game, err

}
