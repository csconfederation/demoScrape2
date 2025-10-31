package events

import (
	"github.com/csconfederation/demoparser3/types"
	"github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs"
	"github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs/events"
	"github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs/msg"
)

const tradeCutoff = 4 // in seconds

func RegisterGameEvents(parser demoinfocs.Parser, game *types.Game) {
	parser.RegisterNetMessageHandler(func(msg *msg.CSVCMsg_ServerInfo) {
		game.MapName = *msg.MapName
	})

	parser.RegisterEventHandler(func(event events.FrameDone) {

		if !parser.GameState().IsMatchStarted() {
			return
		}

		//TODO: Maybe
		//game.Flags.TicksProcessed += 1

		//this will be triggered every 4 seconds of in round time after the first 10 seconds
		//check for lurker
		//TODO: I don't think we're using this
		//if game.Flags.TAlive > 2 && !game.Flags.PostWinCon && parser.GameState().IngameTick() > (18*game.TickRate)+game.PotentialRound.StartingTick {
		//	membersT := types.GetTeamMembers(parser.GameState().TeamTerrorists(), game, parser)
		//	for _, terrorist := range membersT {
		//		if terrorist.IsAlive() {
		//			for _, teammate := range membersT {
		//				if terrorist.SteamID64 != teammate.SteamID64 && teammate.IsAlive() {
		//					dist := int(terrorist.Position().Distance(teammate.Position()))
		//					if dist < 500 {
		//						//invalidate the lurk blip b/c we have a close teammate
		//						game.PotentialRound.PlayerStats[terrorist.SteamID64].DistanceToTeammates = -999999
		//					}
		//					if game.PotentialRound.PlayerStats[terrorist.SteamID64] != nil {
		//						game.PotentialRound.PlayerStats[terrorist.SteamID64].DistanceToTeammates += dist
		//					} else {
		//						log.Debug("THIS IS WHERE WE BROKE_______________________________---------------------------------------------------")
		//					}
		//				}
		//			}
		//		}
		//	}
		//	var lurkerSteam uint64
		//	lurkerDist := 999999
		//	for _, terrorist := range membersT {
		//		if terrorist.IsAlive() {
		//			if game.PotentialRound.PlayerStats[terrorist.SteamID64] == nil {
		//				log.Debug(terrorist.Name)
		//			} else {
		//				dist := game.PotentialRound.PlayerStats[terrorist.SteamID64].DistanceToTeammates
		//				if dist < lurkerDist && dist > 0 {
		//					lurkerDist = dist
		//					lurkerSteam = terrorist.SteamID64
		//				}
		//			}
		//
		//		}
		//	}
		//	if lurkerSteam != 0 {
		//		game.PotentialRound.PlayerStats[lurkerSteam].LurkerBlips += 1
		//	}
		//}
	})
	//
	//round end official doesn't fire on the last round
	// TODO: Do we need this event
	//	parser.RegisterEventHandler(func(event events.ScoreUpdated) {
	//		//CS2 swapped this event to be before RoundEnd
	//		//We have relied on this as a backup for failed RoundEnd events
	//		//may revisit depending on event reliability
	//
	//		if !game.Flags.IsGameLive {
	//			return
	//		}
	//
	//		//added to ensure that a bad round that gets finished does not prematurely finish the game since we track score separately
	//		//we take the existing preupdate score of the updating team score
	//		updatedTeam := game.Teams[types.ValidateTeamName(event.TeamState, game.Teams, game.PotentialRound)]
	//		//and compare to the old score from scoreboard
	//		if event.OldScore != updatedTeam.Score {
	//			updatedTeam.Score = event.OldScore
	//		}
	//	})
	//
	parser.RegisterEventHandler(func(event events.GrenadeProjectileThrow) {
		if !parser.GameState().IsMatchStarted() {
			return
		}
		player := game.CurrentRound.AllPlayersStats[event.Projectile.Thrower.SteamID64]
		player.GrenadeThrown(event.Projectile.WeaponInstance.Type)
	})

	//TODO: Complicated logic I do not want to touch just yet
	//Register handler on kill events
	//parser.RegisterEventHandler(func(event events.Kill) {
	//
	//
	//
	//
	//	flashAssister := ""
	//	if game.Flags.IsGameLive && isDuringExpectedRound(game, parser) {
	//		pS := game.PotentialRound.PlayerStats
	//		tick := parser.GameState().IngameTick()
	//
	//		killerExists := false
	//		victimExists := false
	//		assisterExists := false
	//		if e.Killer != nil && pS[e.Killer.SteamID64] != nil {
	//			killerExists = true
	//		}
	//		if e.Victim != nil && pS[e.Victim.SteamID64] != nil {
	//			victimExists = true
	//		}
	//		if e.Assister != nil && pS[e.Assister.SteamID64] != nil {
	//			assisterExists = true
	//		}
	//		if e.Weapon.Type == common.EqBomb && game.PotentialRound.IsRoundFinalInHalf() {
	//			killerExists = false
	//			victimExists = false
	//			assisterExists = false
	//		}
	//
	//		killValue := 1.0
	//		multiplier := 1.0
	//		traded := false
	//		assisted := false
	//		flashAssisted := false
	//
	//		//death logic (traded here)
	//		if victimExists {
	//			pS[e.Victim.SteamID64].Deaths += 1
	//			pS[e.Victim.SteamID64].DeathTick = tick
	//			if e.Victim.Team == 2 {
	//				game.Flags.TAlive -= 1
	//				pS[e.Victim.SteamID64].DeathPlacement = float64(game.PotentialRound.InitTerroristCount - game.Flags.TAlive)
	//				//pS[e.Victim.SteamID64].tADP = float64(Game.potentialRound.initTerroristCount - Game.flags.tAlive)
	//			} else if e.Victim.Team == 3 {
	//				game.Flags.CtAlive -= 1
	//				pS[e.Victim.SteamID64].DeathPlacement = float64(game.PotentialRound.InitCTerroristCount - game.Flags.CtAlive)
	//				//pS[e.Victim.SteamID64].ctADP = float64(Game.potentialRound.initCTerroristCount - Game.flags.ctAlive)
	//			} else {
	//				//else log an error
	//			}
	//
	//			//do 4v5 calc
	//			if game.Flags.OpeningKill && game.PotentialRound.InitCTerroristCount+game.PotentialRound.InitTerroristCount == 10 {
	//				//the 10th player died
	//				_4v5Team := pS[e.Victim.SteamID64].TeamClanName
	//				game.PotentialRound.TeamStats[_4v5Team].FourVFiveS = 1
	//				for teamName, team := range game.PotentialRound.TeamStats {
	//					if teamName != _4v5Team {
	//						team.FiveVFourS = 1
	//					}
	//				}
	//			}
	//
	//			//add support damage
	//			for suppSteam, suppDMG := range pS[e.Victim.SteamID64].DamageList {
	//				if killerExists && suppSteam != e.Killer.SteamID64 {
	//					pS[suppSteam].SuppDamage += suppDMG
	//					if pS[suppSteam].SuppDamage > 60 {
	//						pS[suppSteam].SuppRounds = 1
	//					}
	//				} else if !killerExists {
	//					pS[suppSteam].SuppDamage += suppDMG
	//					if pS[suppSteam].SuppDamage > 60 {
	//						pS[suppSteam].SuppRounds = 1
	//					}
	//				}
	//
	//			}
	//
	//			//check clutch start
	//
	//			if !game.Flags.PostWinCon {
	//				if game.Flags.TAlive == 1 && game.Flags.TClutchVal == 0 {
	//					game.Flags.TClutchVal = game.Flags.CtAlive
	//					membersT := types.GetTeamMembers(parser.GameState().TeamTerrorists(), game, parser)
	//					for _, terrorist := range membersT {
	//						if terrorist.IsAlive() && e.Victim.SteamID64 != terrorist.SteamID64 {
	//							game.Flags.TClutchSteam = terrorist.SteamID64
	//							log.Debug("Clutch opportunity:", terrorist.Name, game.Flags.TClutchVal)
	//						}
	//					}
	//				}
	//				if game.Flags.CtAlive == 1 && game.Flags.CtClutchVal == 0 {
	//					game.Flags.CtClutchVal = game.Flags.TAlive
	//					membersCT := types.GetTeamMembers(parser.GameState().TeamCounterTerrorists(), game, parser)
	//					for _, counterTerrorist := range membersCT {
	//						if counterTerrorist.IsAlive() && e.Victim.SteamID64 != counterTerrorist.SteamID64 {
	//							game.Flags.CtClutchSteam = counterTerrorist.SteamID64
	//							log.Debug("Clutch opportunity:", counterTerrorist.Name, game.Flags.CtClutchVal)
	//						}
	//					}
	//				}
	//			}
	//
	//			pS[e.Victim.SteamID64].TicksAlive = tick - game.PotentialRound.StartingTick
	//			for deadGuySteam, deadTick := range (*game.PotentialRound).PlayerStats[e.Victim.SteamID64].TradeList {
	//				if tick-deadTick < tradeCutoff*game.TickRate {
	//					pS[deadGuySteam].Traded = 1
	//					pS[deadGuySteam].Eac += 1
	//					pS[deadGuySteam].KastRounds = 1
	//				}
	//			}
	//		}
	//
	//		//assist logic
	//		if assisterExists && victimExists && e.Assister.TeamState.ID() != e.Victim.TeamState.ID() {
	//			//this logic needs to be replaced -yeti does not remember why he wrote this
	//			pS[e.Assister.SteamID64].Assists += 1
	//			pS[e.Assister.SteamID64].Eac += 1
	//			pS[e.Assister.SteamID64].KastRounds = 1
	//			pS[e.Assister.SteamID64].SuppRounds = 1
	//			assisted = true
	//			if e.AssistedFlash {
	//				pS[e.Assister.SteamID64].FAss += 1
	//				flashAssisted = true
	//				flashAssister = e.Assister.Name
	//				log.Debug("VALVE FLASH ASSIST")
	//			} else if float64(parser.GameState().IngameTick()) < pS[e.Victim.SteamID64].MostRecentFlashVal {
	//				//this will trigger if there is both a flash assist and a damage assist
	//				pS[pS[e.Victim.SteamID64].MostRecentFlasher].FAss += 1
	//				pS[pS[e.Victim.SteamID64].MostRecentFlasher].Eac += 1
	//				pS[pS[e.Victim.SteamID64].MostRecentFlasher].SuppRounds = 1
	//				flashAssisted = true
	//				flashAssister = pS[pS[e.Victim.SteamID64].MostRecentFlasher].Name
	//			}
	//
	//		}
	//
	//		//kill logic (trades here)
	//		if killerExists && victimExists && e.Killer.TeamState.ID() != e.Victim.TeamState.ID() {
	//			pS[e.Killer.SteamID64].Kills += 1
	//			pS[e.Killer.SteamID64].KastRounds = 1
	//			pS[e.Killer.SteamID64].Rwk = 1
	//			pS[e.Killer.SteamID64].TradeList[e.Victim.SteamID64] = tick
	//			if e.Weapon.Type == 309 {
	//				pS[e.Killer.SteamID64].AwpKills += 1
	//				if e.Killer.Team == 3 {
	//					pS[e.Killer.SteamID64].CtAWP += 1
	//				}
	//			}
	//			if e.IsHeadshot {
	//				pS[e.Killer.SteamID64].Hs += 1
	//			}
	//			for _, deadTick := range (*game.PotentialRound).PlayerStats[e.Victim.SteamID64].TradeList {
	//				if tick-deadTick < tradeCutoff*game.TickRate {
	//					pS[e.Killer.SteamID64].Trades += 1
	//					traded = true
	//					break
	//				}
	//			}
	//
	//			killerTeam := e.Killer.Team
	//			if game.Flags.PrePlant {
	//				//normal base value
	//				if killerTeam == 2 {
	//					//taking site by T
	//					killValue = 1.2
	//				} else if killerTeam == 3 {
	//					//site Defense by CT
	//					killValue = 1
	//				}
	//			} else if game.Flags.PostPlant {
	//				//site D or retake
	//				if killerTeam == 2 {
	//					//site Defense by T
	//					killValue = 1
	//				} else if killerTeam == 3 {
	//					//retake
	//					killValue = 1.2
	//				}
	//			} else if game.Flags.PostWinCon {
	//				//exit or chase
	//				if game.PotentialRound.WinnerENUM == 2 { //Ts win
	//					if killerTeam == 2 { //chase
	//						killValue = 0.8
	//					}
	//					if killerTeam == 3 { //exit
	//						killValue = 0.6
	//					}
	//				} else if game.PotentialRound.WinnerENUM == 3 { //CTs win
	//					if killerTeam == 2 { //T kill in lost round
	//						killValue = 0.5
	//					}
	//					if killerTeam == 3 { //CT kill in won round
	//						if game.Flags.TMoney {
	//							killValue = 0.6
	//						} else {
	//							killValue = 0.8
	//						}
	//					}
	//				}
	//			}
	//
	//			if game.Flags.OpeningKill {
	//				game.Flags.OpeningKill = false
	//
	//				pS[e.Killer.SteamID64].Ok = 1
	//				pS[e.Victim.SteamID64].Ol = 1
	//
	//				if killerTeam == 2 { //T entry/opener {
	//					if game.Flags.PrePlant {
	//						multiplier += 0.8
	//						pS[e.Killer.SteamID64].Entries = 1
	//					} else {
	//						multiplier += 0.3
	//					}
	//				} else if killerTeam == 3 { //CT opener
	//					multiplier += 0.5
	//				}
	//
	//			} else if traded {
	//				multiplier += 0.3
	//			}
	//
	//			if flashAssisted { //flash assisted kill
	//				multiplier += 0.2
	//			}
	//			if assisted { //assisted kill
	//				killValue -= 0.15
	//				pS[e.Assister.SteamID64].ImpactPoints += 0.15
	//			}
	//
	//			killValue *= multiplier
	//
	//			ecoRatio := float64(e.Victim.EquipmentValueCurrent()) / float64(e.Killer.EquipmentValueCurrent())
	//			ecoMod := 1.0
	//			if ecoRatio > 4 {
	//				ecoMod += 0.25
	//			} else if ecoRatio > 2 {
	//				ecoMod += 0.14
	//			} else if ecoRatio < 0.25 {
	//				ecoMod -= 0.25
	//			} else if ecoRatio < 0.5 {
	//				ecoMod -= 0.14
	//			}
	//			killValue *= ecoMod
	//
	//			pS[e.Killer.SteamID64].KillPoints += killValue
	//		}
	//
	//	}
	//	var hs string
	//	if e.IsHeadshot {
	//		hs = " (HS)"
	//	}
	//	var wallBang string
	//	if e.PenetratedObjects > 0 {
	//		wallBang = " (WB)"
	//	}
	//	log.Debug("%s <%v%s%s> %s at %d flash assist by %s\n", e.Killer, e.Weapon, hs, wallBang, e.Victim, parser.GameState().IngameTick(), flashAssister)
	//})

	//TODO: Might not need this
	//parser.RegisterEventHandler(func(e events.Footstep) {
	//	if game.Flags.IsGameLive {
	//		game.Flags.InRound = true
	//	}
	//
	//})
}
