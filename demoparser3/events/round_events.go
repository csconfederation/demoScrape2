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
		game.CurrentRound = types.NewRound()
		game.CurrentRound.RoundNum = parser.GameState().TotalRoundsPlayed() + 1

		counterTerrorists := parser.GameState().TeamCounterTerrorists()
		terrorists := parser.GameState().TeamTerrorists()

		if game.CurrentRound.RoundNum == 1 {
			game.Teams[counterTerrorists.ClanName()] = types.NewTeam(counterTerrorists.ClanName())
			game.Teams[terrorists.ClanName()] = types.NewTeam(terrorists.ClanName())
		}

		for _, player := range counterTerrorists.Members() {
			game.CurrentRound.AllPlayersStats[player.SteamID64] = types.NewPlayerStats()
		}
		for _, player := range terrorists.Members() {
			game.CurrentRound.AllPlayersStats[player.SteamID64] = types.NewPlayerStats()
		}
	})

	parser.RegisterEventHandler(func(event events.RoundEnd) {

		if !parser.GameState().IsMatchStarted() {
			return
		}

		// TODO: Check if required
		//game.Flags.DidRoundEndFire = true

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

		game.CurrentRound.RoundEnd(game.Teams[event.WinnerState.ClanName()],
			game.Teams[event.LoserState.ClanName()], event.Reason)
	})

	parser.RegisterEventHandler(func(event events.RoundFreezetimeEnd) {
		game.CurrentRound.RoundFreezetimeEnd(parser.GameState().TotalRoundsPlayed())
	})

	//round end official doesn't fire on the last round
	parser.RegisterEventHandler(func(event events.RoundEndOfficial) {

		//TODO: Do we need this
		//if !game.Flags.DidRoundEndFire {
		//	game.Flags.RoundIntegrityEnd -= 1
		//}
		//

		//totalRoundsPlayed := parser.GameState().TotalRoundsPlayed()
		//if game.Flags.IsGameLive && game.Flags.RoundIntegrityEnd == totalRoundsPlayed {
		//	game.ProcessRoundFinal(false, parser.GameState().IngameTick(), totalRoundsPlayed)
		//}

		game.Rounds = append(game.Rounds, game.CurrentRound)
	})
}
