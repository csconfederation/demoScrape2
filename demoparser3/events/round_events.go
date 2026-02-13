package events

import (
	"github.com/csconfederation/demoparser3/logger"
	"github.com/csconfederation/demoparser3/types"
	"github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs"
	"github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs/common"
	"github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs/events"
)

func RegisterRoundEvents(parser demoinfocs.Parser, game *types.Game) {

	parser.RegisterEventHandler(func(e events.RoundStart) {
		game.CurrentRound = types.NewRound(parser.GameState().IngameTick())
		game.CurrentRound.RoundNum = parser.GameState().TotalRoundsPlayed() + 1

		counterTerrorists := parser.GameState().TeamCounterTerrorists()
		terrorists := parser.GameState().TeamTerrorists()

		if game.CurrentRound.RoundNum == 1 {
			game.Start(counterTerrorists, terrorists)
		}

		connectedTeamPlayers := 0
		for _, player := range counterTerrorists.Members() {
			game.CurrentRound.AllPlayersStats[player.SteamID64] = types.NewPlayerStats()
			if player.IsConnected {
				connectedTeamPlayers += 1
			}
		}
		game.CurrentRound.TeamStats[counterTerrorists.ClanName()] = types.NewTeamStats(connectedTeamPlayers)

		connectedTeamPlayers = 0
		for _, player := range terrorists.Members() {
			game.CurrentRound.AllPlayersStats[player.SteamID64] = types.NewPlayerStats()
			if player.IsConnected {
				connectedTeamPlayers += 1
			}
		}
		game.CurrentRound.TeamStats[terrorists.ClanName()] = types.NewTeamStats(connectedTeamPlayers)
	})

	parser.RegisterEventHandler(func(event events.RoundEnd) {

		if !parser.GameState().IsMatchStarted() {
			return
		}

		if event.Winner == common.TeamUnassigned || event.Winner == common.TeamSpectators {
			logger.Warn("An invalid team won round %d or it's a draw", game.CurrentRound.RoundNum)
			return
		}

		if game.CurrentRound == nil {
			logger.Error("Game has no current round - RoundEnd")
			return
		}

		if event.Winner == common.TeamTerrorists {
			//TODO: Check if needed
			//game.CurrentRound.TMoney = true
		}

		if game.Teams[event.WinnerState.ClanName()] == nil {
			logger.Error("Team not found at round end",
				"clanName", event.WinnerState.ClanName(),
				"winnerState", event.WinnerState,
			)
			return
		}

		game.Teams[event.WinnerState.ClanName()].Score += 1

		game.CurrentRound.RoundEnd(game.Teams[event.WinnerState.ClanName()],
			game.Teams[event.LoserState.ClanName()], event.Reason)

		if game.CurrentRound.IsFinalRound {
			game.Rounds = append(game.Rounds, game.CurrentRound)
			game.CurrentRound.IsValid = game.CurrentRound.RoundNum == parser.GameState().TotalRoundsPlayed()
		}
	})

	// TODO: Test when pistol round is reset
	parser.RegisterEventHandler(func(event events.RoundFreezetimeEnd) {

		actualRoundsPlayed := 0
		for _, team := range game.Teams {
			actualRoundsPlayed += team.Score
		}

		game.CurrentRound.RoundFreezetimeEnd(actualRoundsPlayed)
	})

	//round end official doesn't fire on the last round
	parser.RegisterEventHandler(func(event events.RoundEndOfficial) {
		game.Rounds = append(game.Rounds, game.CurrentRound)

		// round integrity check
		game.CurrentRound.IsValid = game.CurrentRound.RoundNum == parser.GameState().TotalRoundsPlayed()
	})
}
