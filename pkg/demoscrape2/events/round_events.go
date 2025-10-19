package events

import (
	"math"

	"github.com/csconfederation/demoScrape2/pkg/demoscrape2/helpers"
	types "github.com/csconfederation/demoScrape2/pkg/demoscrape2/types"
	dem "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
	log "github.com/sirupsen/logrus"
)

func RegisterRoundEvents(parser dem.Parser, game *types.Game) {
	parser.RegisterEventHandler(func(event events.RoundStart) {
		log.Debug("Round Start", parser.GameState().TotalRoundsPlayed())
		game.Flags.RoundStartedAt = parser.GameState().IngameTick()
	})

	parser.RegisterEventHandler(func(event events.RoundFreezetimeEnd) {
		log.Debug("Round Freeze Time End\n")
		pistol := false

		//we are going to check to see if the first pistol is actually starting
		membersT := types.GetTeamMembers(parser.GameState().TeamTerrorists(), game, parser)
		membersCT := types.GetTeamMembers(parser.GameState().TeamCounterTerrorists(), game, parser)
		if len(membersT) != 0 && len(membersCT) != 0 {
			if membersT[0].Money()+membersT[0].MoneySpentThisRound() == 800 && membersCT[0].Money()+
				membersCT[0].MoneySpentThisRound() == 800 {
				//start the Game
				if !game.Flags.HasGameStarted {
					tickRate := int(math.Round(parser.TickRate()))
					teamTerrorists := parser.GameState().TeamTerrorists()
					teamCounterTerrorists := parser.GameState().TeamCounterTerrorists()
					game.Start(tickRate, teamTerrorists, teamCounterTerrorists)
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

		if !game.Flags.IsGameLive {
			return
		}

		game.Flags.InRound = true
		game.PotentialRound.Start(game, parser)
		if pistol {
			for _, team := range game.PotentialRound.TeamStats {
				team.Pistols = 1
			}
		}
	})

	parser.RegisterEventHandler(func(event events.RoundEnd) {

		if !game.Flags.IsGameLive {
			return
		}

		game.Flags.DidRoundEndFire = true

		log.Debug("Round", parser.GameState().TotalRoundsPlayed(), "End", event.WinnerState.ClanName(), "won", "this determined from e.WinnerState.ClanName()")

		log.Debug("e.WinnerState.ID()", event.WinnerState.ID(), "and", "e.Winner", event.Winner, "and", "e.WinnerState.Team()", event.WinnerState.Team())

		validWinner := true
		if event.Winner != common.TeamTerrorists && event.Winner != common.TeamCounterTerrorists {
			validWinner = false
			//and set the integrity flag to false

		} else if event.Winner == common.TeamTerrorists {
			game.Flags.TMoney = true
		} else {
			//we need to check if the Game is over
		}

		//we want to actually process the round
		totalRoundsPlayed := parser.GameState().TotalRoundsPlayed()
		if validWinner && game.Flags.RoundIntegrityStart == totalRoundsPlayed {
			game.PotentialRound.WinnerENUM = int(event.Winner)
			game.PotentialRound.RoundEndReason = helpers.RoundEndReasons[int(event.Reason)]
			game.ProcessRoundOnWinCon(event.WinnerState, totalRoundsPlayed)

			//check last round
			roundWinnerScore := game.Teams[types.ValidateTeamName(event.WinnerState, game.Teams, game.PotentialRound)].Score
			roundLoserScore := game.Teams[types.ValidateTeamName(event.LoserState, game.Teams, game.PotentialRound)].Score
			log.Debug("winner Rounds", roundWinnerScore)
			log.Debug("loser Rounds", roundLoserScore)

			if game.RoundsToWin == 16 {
				//check for normal win
				if roundWinnerScore == 16 && roundLoserScore < 15 {
					//normal win
					game.WinnerClanName = game.PotentialRound.WinnerClanName
					game.ProcessRoundFinal(true, parser.GameState().IngameTick(), totalRoundsPlayed)
				} else if roundWinnerScore > 15 { //check for OT win
					overtime := ((roundWinnerScore+roundLoserScore)-30-1)/6 + 1
					//OT win
					if (roundWinnerScore-15-1)/3 == overtime {
						game.WinnerClanName = game.PotentialRound.WinnerClanName
						game.ProcessRoundFinal(true, parser.GameState().IngameTick(), totalRoundsPlayed)
					}
				}
			} else if game.RoundsToWin == 9 {
				//check for normal win
				if roundWinnerScore == 9 && roundLoserScore < 8 {
					//normal win
					game.WinnerClanName = game.PotentialRound.WinnerClanName
					game.ProcessRoundFinal(true, parser.GameState().IngameTick(), totalRoundsPlayed)
				} else if roundWinnerScore == 8 && roundLoserScore == 8 { //check for tie
					//tie
					game.WinnerClanName = game.PotentialRound.WinnerClanName
					game.ProcessRoundFinal(true, parser.GameState().IngameTick(), totalRoundsPlayed)
				}
			} else if game.RoundsToWin == 13 {
				//check for normal win
				if roundWinnerScore == 13 && roundLoserScore < 12 {
					//normal win
					game.WinnerClanName = game.PotentialRound.WinnerClanName
					game.ProcessRoundFinal(true, parser.GameState().IngameTick(), totalRoundsPlayed)
				} else if roundWinnerScore > 12 { //check for OT win
					overtime := ((roundWinnerScore+roundLoserScore)-24-1)/6 + 1
					//OT win
					if (roundWinnerScore-12-1)/3 == overtime {
						game.WinnerClanName = game.PotentialRound.WinnerClanName
						game.ProcessRoundFinal(true, parser.GameState().IngameTick(), totalRoundsPlayed)
					}
				}
			}
		}

		//check last round
		//or check overtime win
	})

	//round end official doesn't fire on the last round
	parser.RegisterEventHandler(func(e events.RoundEndOfficial) {

		log.Debug("Round End Official\n")

		if !game.Flags.DidRoundEndFire {
			game.Flags.RoundIntegrityEnd -= 1
		}

		log.Debug("isGameLive", game.Flags.IsGameLive, "roundIntegrityEnd", game.Flags.RoundIntegrityEnd,
			"pTotalRoundsPlayed", parser.GameState().TotalRoundsPlayed())

		totalRoundsPlayed := parser.GameState().TotalRoundsPlayed()
		if game.Flags.IsGameLive && game.Flags.RoundIntegrityEnd == totalRoundsPlayed {
			game.ProcessRoundFinal(false, parser.GameState().IngameTick(), totalRoundsPlayed)
		}
	})
}
