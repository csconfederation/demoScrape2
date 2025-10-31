package events

import (
	"github.com/csconfederation/demoparser3/logger"
	"github.com/csconfederation/demoparser3/types"
	"github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs"
	"github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs/common"
	"github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs/events"
)

func RegisterPlayerEvents(parser demoinfocs.Parser, game *types.Game) {
	//parser.RegisterEventHandler(func(event events.PlayerInfo) {
	//	//TODO: Handle connected and reconnected players (if required)
	//	//player := parser.GameState().Participants().AllByUserID()[event.Index]
	//	//if player != nil {
	//	//	game.ReconnectedPlayers[player.SteamID64] = true
	//	//	if game.Flags.InRound && game.Flags.IsGameLive {
	//	//		game.ConnectedAfterRoundStart[player.SteamID64] = true
	//	//	}
	//	//}
	//})

	parser.RegisterEventHandler(func(event events.PlayerHurt) {

		if !parser.GameState().IsMatchStarted() {
			return
		}

		if event.Weapon.Type == common.EqBomb || event.Weapon.Type == common.EqUnknown {
			return
		}

		if event.Player == nil {
			logger.Error("Player is nil")
			return
		}

		if event.Attacker == nil {
			logger.Error("Attacker is nil")
			return
		}

		attackerStats := game.CurrentRound.AllPlayersStats[event.Attacker.SteamID64]
		game.CurrentRound.AllPlayersStats[event.Player.SteamID64].PlayerHurt(attackerStats, event.HealthDamageTaken, event.Weapon.Type)

	})

	parser.RegisterEventHandler(func(event events.PlayerFlashed) {

		if !parser.GameState().IsMatchStarted() {
			return
		}

		if event.Player == nil {
			logger.Error("Player is nil")
			return
		}

		if event.Attacker == nil {
			logger.Error("Attacker is nil")
			return
		}

		if event.Player.Team == event.Attacker.Team {
			// team flash
			return
		}

		attackerStats := game.CurrentRound.AllPlayersStats[event.Attacker.SteamID64]
		game.CurrentRound.AllPlayersStats[event.Player.SteamID64].PlayerFlashed(attackerStats, event.FlashDuration())
	})
}

//TODO: Check if this event even affects anything significant

//	parser.RegisterEventHandler(func(event events.PlayerTeamChange) {
//		log.Debug("Player Changed Team:", event.Player, event.OldTeam, event.NewTeam)
//
//		if !game.Flags.IsGameLive || !game.Flags.InRound {
//			return
//		}
//
//		// joins same team
//		if event.NewTeam <= 1 {
//			return
//		}
//
//		//we are joining an actual team
//		if game.PotentialRound.PlayerStats[event.Player.SteamID64] == nil && event.Player.IsBot && event.Player.IsAlive() {
//			player := types.NewPlayerStatsFromPlayer(event.Player, event.NewTeamState, game)
//			game.PotentialRound.PlayerStats[player.SteamID] = player
//		}
//	})
//
//TODO: Check if this event even affects anything significant

//	parser.RegisterEventHandler(func(event events.PlayerDisconnected) {
//		log.Debug("Player DC", event.Player)
//
//		if game.ReconnectedPlayers[event.Player.SteamID64] {
//			game.ReconnectedPlayers[event.Player.SteamID64] = false
//		}
//
//		if !game.Flags.IsGameLive {
//			return
//		}
//
//		//update alive players
//
//		game.Flags.TAlive = 0
//		game.Flags.CtAlive = 0
//
//		membersT := types.GetTeamMembers(parser.GameState().TeamTerrorists(), game, parser)
//		for _, terrorist := range membersT {
//			if terrorist.IsAlive() {
//				game.Flags.TAlive += 1
//			}
//		}
//		membersCT := types.GetTeamMembers(parser.GameState().TeamCounterTerrorists(), game, parser)
//		for _, counterTerrorist := range membersCT {
//			if counterTerrorist.IsAlive() {
//				game.Flags.CtAlive += 1
//			}
//		}
//
//	})
//}
