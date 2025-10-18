package events

import (
	"github.com/csconfederation/demoScrape2/pkg/demoscrape2/types"
	dem "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
	log "github.com/sirupsen/logrus"
)

func RegisterBombEvents(parser dem.Parser, game *types.Game) {
	parser.RegisterEventHandler(func(e events.BombPlanted) {
		log.Debug("Bomb Planted\n")
		if game.Flags.IsGameLive && !game.Flags.PostWinCon {
			game.Flags.PrePlant = false
			game.Flags.PostPlant = true
			game.Flags.TMoney = true
			game.PotentialRound.Planter = e.BombEvent.Player.SteamID64
			game.PotentialRound.BombStartTick = parser.GameState().IngameTick()
		}
	})

	parser.RegisterEventHandler(func(e events.BombDefused) {
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

	parser.RegisterEventHandler(func(e events.BombExplode) {
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
}
