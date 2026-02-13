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

	parser.RegisterEventHandler(func(event events.GameHalfEnded) {
		game.TSide, game.CTSide = game.CTSide, game.TSide
	})

	parser.RegisterEventHandler(func(event events.FrameDone) {

		if !parser.GameState().IsMatchStarted() {
			return
		}

		//this will be triggered every 4 seconds of in round time after the first 10 seconds
		//check for lurker
		if game.CurrentRound.TeamStats[game.TSide].MembersAlive <= 2 {
			return
		}

		if game.CurrentRound.PostWinCon {
			return
		}

		currentTick := parser.GameState().IngameTick()

		// 18 seconds should have passed after round start
		if currentTick <= (18*types.TickRate)+game.CurrentRound.StartingTick {
			return
		}

		lurker := game.CheckLurk(currentTick, parser.GameState().TeamTerrorists())
		if lurker != 0 {
			game.CurrentRound.AllPlayersStats[lurker].LurkerBlips += 1
		}
	})

	//round end official doesn't fire on the last round
	parser.RegisterEventHandler(func(event events.ScoreUpdated) {
		//CS2 swapped this event to be before RoundEnd
		//We have relied on this as a backup for failed RoundEnd events
		//may revisit depending on event reliability

		if !parser.GameState().IsMatchStarted() {
			return
		}

		//added to ensure that a bad round that gets finished does not prematurely finish the game since we track score separately
		//we take the existing preupdate score of the updating team score
		updatedTeam := game.Teams[event.TeamState.ClanName()]
		//and compare to the old score from scoreboard
		if event.OldScore != updatedTeam.Score {
			updatedTeam.Score = event.OldScore
		}
	})

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
	//	if !parser.GameState().IsMatchStarted() {
	//		return
	//	}
	//
	//	if event.Weapon.Type == common.EqBomb {
	//		return
	//	}
	//
	//	if event.Victim == nil {
	//		logger.Error("Victim is nil", "event", "Kill")
	//		return
	//	}
	//
	//	// team kill
	//	if event.Victim.TeamState.ClanName() == event.Killer.TeamState.ClanName(){
	//		return
	//	}
	//
	//
	//	if event.Killer != nil {
	//		game.ProcessKiller(event.Killer, event.Weapon.Type)
	//	}
	//
	//	game.ProcessVictim(event.Victim)
	//
	//
	//
	//	// TODO: Might have to redo this
	//	game.PlayerKilled(event.Victim, event.Killer, event.Assister, event.AssistedFlash)
	//
	//	//pS[e.Victim.SteamID64].TicksAlive = tick - game.PotentialRound.StartingTick
	//	//for deadGuySteam, deadTick := range (*game.PotentialRound).PlayerStats[e.Victim.SteamID64].TradeList {
	//	//	if tick-deadTick < tradeCutoff*game.TickRate {
	//	//		pS[deadGuySteam].Traded = 1
	//	//		pS[deadGuySteam].Eac += 1
	//	//		pS[deadGuySteam].KastRounds = 1
	//	//	}
	//	//}
	//
	//	//kill logic (trades here)
	//	if killerExists && victimExists && e.Killer.TeamState.ID() != e.Victim.TeamState.ID() {
	//		pS[e.Killer.SteamID64].Kills += 1
	//		pS[e.Killer.SteamID64].KastRounds = 1
	//		pS[e.Killer.SteamID64].Rwk = 1
	//		pS[e.Killer.SteamID64].TradeList[e.Victim.SteamID64] = tick
	//		if e.Weapon.Type == 309 {
	//			pS[e.Killer.SteamID64].AwpKills += 1
	//			if e.Killer.Team == 3 {
	//				pS[e.Killer.SteamID64].CtAWP += 1
	//			}
	//		}
	//		if e.IsHeadshot {
	//			pS[e.Killer.SteamID64].Hs += 1
	//		}
	//		for _, deadTick := range (*game.PotentialRound).PlayerStats[e.Victim.SteamID64].TradeList {
	//			if tick-deadTick < tradeCutoff*game.TickRate {
	//				pS[e.Killer.SteamID64].Trades += 1
	//				traded = true
	//				break
	//			}
	//		}
	//
	//		killerTeam := e.Killer.Team
	//		if game.Flags.PrePlant {
	//			//normal base value
	//			if killerTeam == 2 {
	//				//taking site by T
	//				killValue = 1.2
	//			} else if killerTeam == 3 {
	//				//site Defense by CT
	//				killValue = 1
	//			}
	//		} else if game.Flags.PostPlant {
	//			//site D or retake
	//			if killerTeam == 2 {
	//				//site Defense by T
	//				killValue = 1
	//			} else if killerTeam == 3 {
	//				//retake
	//				killValue = 1.2
	//			}
	//		} else if game.Flags.PostWinCon {
	//			//exit or chase
	//			if game.PotentialRound.WinnerENUM == 2 { //Ts win
	//				if killerTeam == 2 { //chase
	//					killValue = 0.8
	//				}
	//				if killerTeam == 3 { //exit
	//					killValue = 0.6
	//				}
	//			} else if game.PotentialRound.WinnerENUM == 3 { //CTs win
	//				if killerTeam == 2 { //T kill in lost round
	//					killValue = 0.5
	//				}
	//				if killerTeam == 3 { //CT kill in won round
	//					if game.Flags.TMoney {
	//						killValue = 0.6
	//					} else {
	//						killValue = 0.8
	//					}
	//				}
	//			}
	//		}
	//
	//		if game.Flags.OpeningKill {
	//			game.Flags.OpeningKill = false
	//
	//			pS[e.Killer.SteamID64].Ok = 1
	//			pS[e.Victim.SteamID64].Ol = 1
	//
	//			if killerTeam == 2 { //T entry/opener {
	//				if game.Flags.PrePlant {
	//					multiplier += 0.8
	//					pS[e.Killer.SteamID64].Entries = 1
	//				} else {
	//					multiplier += 0.3
	//				}
	//			} else if killerTeam == 3 { //CT opener
	//				multiplier += 0.5
	//			}
	//
	//		} else if traded {
	//			multiplier += 0.3
	//		}
	//
	//		if flashAssisted { //flash assisted kill
	//			multiplier += 0.2
	//		}
	//		if assisted { //assisted kill
	//			killValue -= 0.15
	//			pS[e.Assister.SteamID64].ImpactPoints += 0.15
	//		}
	//
	//		killValue *= multiplier
	//
	//		ecoRatio := float64(e.Victim.EquipmentValueCurrent()) / float64(e.Killer.EquipmentValueCurrent())
	//		ecoMod := 1.0
	//		if ecoRatio > 4 {
	//			ecoMod += 0.25
	//		} else if ecoRatio > 2 {
	//			ecoMod += 0.14
	//		} else if ecoRatio < 0.25 {
	//			ecoMod -= 0.25
	//		} else if ecoRatio < 0.5 {
	//			ecoMod -= 0.14
	//		}
	//		killValue *= ecoMod
	//
	//		pS[e.Killer.SteamID64].KillPoints += killValue
	//	}
	//}
}
