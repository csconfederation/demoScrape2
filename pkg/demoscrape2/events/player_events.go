package events

import (
	"github.com/csconfederation/demoScrape2/pkg/demoscrape2/types"
	dem "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
	log "github.com/sirupsen/logrus"
)

func RegisterPlayerEvents(parser dem.Parser, game *types.Game) {
	parser.RegisterEventHandler(func(event events.PlayerInfo) {
		player := parser.GameState().Participants().AllByUserID()[event.Index]
		if player != nil {
			game.ReconnectedPlayers[player.SteamID64] = true
			if game.Flags.InRound && game.Flags.IsGameLive {
				game.ConnectedAfterRoundStart[player.SteamID64] = true
			}
		}
	})

	parser.RegisterEventHandler(func(event events.PlayerHurt) {
		if !game.Flags.IsGameLive || event.Player == nil {
			return
		}

		equipment := getEquipmentType(event.Weapon)
		playerStats := game.PotentialRound.PlayerStats[event.Player.SteamID64]

		if shouldUpdateDamageTaken(playerStats, event, equipment, game.PotentialRound,
			parser.GameState().IsFreezetimePeriod()) {
			game.PotentialRound.PlayerStats[event.Player.SteamID64].DamageTaken += event.HealthDamageTaken
		}

		if shouldUpdateAttackerDamage(playerStats, event) {
			attackerStats := game.PotentialRound.PlayerStats[event.Attacker.SteamID64]
			victimStats := game.PotentialRound.PlayerStats[event.Player.SteamID64]
			updateAttackerDamageStats(attackerStats, victimStats, event.HealthDamageTaken, equipment)
		}
	})

	parser.RegisterEventHandler(func(event events.PlayerFlashed) {
		if !game.Flags.IsGameLive || event.Player == nil || event.Attacker == nil {
			return
		}

		tick := float64(parser.GameState().IngameTick())
		blindTicks := event.FlashDuration().Seconds() * float64(game.TickRate)
		victim := event.Player
		flasher := event.Attacker

		if flasher.Team == victim.Team {
			return
		} // team flash
		if blindTicks <= float64(game.TickRate) {
			return
		} // not blind
		if !victim.IsAlive() {
			return
		}
		if float64(victim.FlashDuration) >= (blindTicks/float64(game.TickRate) + 1) {
			return
		} // flash duration

		game.PotentialRound.PlayerStats[flasher.SteamID64].Ef += 1
		game.PotentialRound.PlayerStats[flasher.SteamID64].EnemyFlashTime += blindTicks / float64(game.TickRate)
		if tick+blindTicks > game.PotentialRound.PlayerStats[victim.SteamID64].MostRecentFlashVal {
			game.PotentialRound.PlayerStats[victim.SteamID64].MostRecentFlashVal = tick + blindTicks
			game.PotentialRound.PlayerStats[victim.SteamID64].MostRecentFlasher = flasher.SteamID64
		}
	})

	parser.RegisterEventHandler(func(event events.PlayerTeamChange) {
		log.Debug("Player Changed Team:", event.Player, event.OldTeam, event.NewTeam)

		if !game.Flags.IsGameLive || !game.Flags.InRound {
			return
		}

		// joins same team
		if event.NewTeam <= 1 {
			return
		}

		//we are joining an actual team
		if game.PotentialRound.PlayerStats[event.Player.SteamID64] == nil && event.Player.IsBot && event.Player.IsAlive() {
			player := types.NewPlayerStatsFromPlayer(event.Player, event.NewTeamState, game)
			game.PotentialRound.PlayerStats[player.SteamID] = player
		}
	})

	parser.RegisterEventHandler(func(event events.PlayerDisconnected) {
		log.Debug("Player DC", event.Player)

		if game.ReconnectedPlayers[event.Player.SteamID64] {
			game.ReconnectedPlayers[event.Player.SteamID64] = false
		}

		if !game.Flags.IsGameLive {
			return
		}

		//update alive players

		game.Flags.TAlive = 0
		game.Flags.CtAlive = 0

		membersT := types.GetTeamMembers(parser.GameState().TeamTerrorists(), game, parser)
		for _, terrorist := range membersT {
			if terrorist.IsAlive() {
				game.Flags.TAlive += 1
			}
		}
		membersCT := types.GetTeamMembers(parser.GameState().TeamCounterTerrorists(), game, parser)
		for _, counterTerrorist := range membersCT {
			if counterTerrorist.IsAlive() {
				game.Flags.CtAlive += 1
			}
		}

	})
}

func getEquipmentType(weapon *common.Equipment) common.EquipmentType {
	if weapon == nil {
		return -999
	}
	return weapon.Type
}

func shouldUpdateDamageTaken(playerStats *types.PlayerStats, event events.PlayerHurt,
	equipment common.EquipmentType, round *types.Round, isFreezeTime bool) bool {
	if event.Player == nil || !event.Player.IsConnected {
		return false
	}

	if equipment == common.EqBomb && round.IsRoundFinalInHalf() {
		return false
	}

	if !isFreezeTime {
		if playerStats == nil {
			panic("Connected player exists but has no player stats (not in freeze time)")
		}
		return false
	}

	return true
}

func shouldUpdateAttackerDamage(playerStats *types.PlayerStats, event events.PlayerHurt) bool {
	if event.Player == nil || event.Attacker == nil {
		return false
	}

	if playerStats == nil {
		return false
	}

	return event.Player.Team != event.Attacker.Team
}

func updateAttackerDamageStats(attackerStats, victimStats *types.PlayerStats, damage int,
	equipment common.EquipmentType) {

	attackerStats.Damage += damage
	victimStats.DamageList[attackerStats.SteamID] += damage
	updateUtilityDamage(attackerStats, equipment, damage)
}

func updateUtilityDamage(stats *types.PlayerStats, equipment common.EquipmentType, damage int) {

	if !isUtility(equipment) {
		return
	}

	stats.UtilDmg += damage

	switch equipment {
	case common.EqHE:
		stats.NadeDmg += damage
		break
	case common.EqMolotov, common.EqIncendiary:
		stats.InfernoDmg += damage
		break
	}
}

func isUtility(equipmentType common.EquipmentType) bool {
	return equipmentType == common.EqHE || equipmentType == common.EqIncendiary || equipmentType == common.EqMolotov
}
