package demoscrape2

import (
	"fmt"
	"strings"

	dem "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
)

type Dictionary map[string]interface{}

func contains(players []*common.Player, player *common.Player) bool {
	for _, p := range players {
		if p.SteamID64 == player.SteamID64 {
			return true
		}
	}
	return false
}

func getTeamMembers(team *common.TeamState, game *Game, p dem.Parser) []*common.Player {
	players := team.Members()
	allPlayers := p.GameState().Participants().All()
	// Filter players by the Team from the team state
	teamPlayers := make([]*common.Player, 0)

	// Helper function to find player index in teamPlayers
	findPlayerIndex := func(slice []*common.Player, steamId uint64) int {
		for i, player := range slice {
			if player.SteamID64 == steamId {
				return i
			}
		}
		return -1
	}

	for _, player := range players {
		if player.Team == team.Team() {
			if game.ConnectedAfterRoundStart[player.SteamID64] {
				continue
			}
			teamPlayers = append(teamPlayers, player)
		}
	}

	// Grab reconnected players and check for duplicates
	for steamId, connected := range game.ReconnectedPlayers {
		if !connected {
			continue
		}
		for _, player := range allPlayers {
			// If the player is in connectedAfterRoundStart, do not return them
			if game.ConnectedAfterRoundStart[player.SteamID64] {
				continue
			}

			if player.SteamID64 == steamId && player.Team == team.Team() {
				// Check if player is already in teamPlayers
				idx := findPlayerIndex(teamPlayers, steamId)
				if idx != -1 {
					// Remove the existing record
					teamPlayers = append(teamPlayers[:idx], teamPlayers[idx+1:]...)
				}
				// Append the new record
				teamPlayers = append(teamPlayers, player)
			}
		}
	}

	return teamPlayers
}

func isDuringExpectedRound(game *Game, p dem.Parser) bool {
	isPreWinCon := int(game.PotentialRound.RoundNum) == p.GameState().TotalRoundsPlayed()+1
	isAfterWinCon := int(game.PotentialRound.RoundNum) == p.GameState().TotalRoundsPlayed() && game.Flags.PostWinCon
	return isPreWinCon || isAfterWinCon
}

func isRoundFinalInHalf(round int8) bool {
	return round%MR == 0 || (round > (MR*2) && round%3 == 0)
}

func validateTeamName(game *Game, teamName string, teamNum common.Team) string {
	if teamName != "" && game.Teams["StartedT"] == nil {
		name := ""
		if strings.HasPrefix(teamName, "[") {
			if len(teamName) == 31 {
				//name here will be truncated
				name = strings.Split(teamName, "] ")[1]
				for _, team := range game.Teams {
					if strings.Contains(team.Name, name) {
						return team.Name
					}
				}
				fmt.Print("OH NOEY")
				return name
			} else {
				name = strings.Split(teamName, "] ")[1]
				return name
			}
		} else {
			return teamName
		}
	} else {
		//this demo no have team names, so we are big fucked
		//we are hardcoding during what rounds each team will have what side
		round := game.PotentialRound.RoundNum
		swap := false
		if round >= MR+1 && round <= (MR*2)+3 {
			swap = true
		} else if round >= (MR*2)+4 {
			//we are now in OT hell :)
			if (round-((MR*2)+4))/6%2 != 0 {
				swap = true
			}
		}
		if !swap {
			if teamNum == 2 {
				return "StartedT"
			} else if teamNum == 3 {
				return "StartedCT"
			}
		} else {
			if teamNum == 2 {
				return "StartedCT"
			} else if teamNum == 3 {
				return "StartedT"
			}
			return "SPECs"
		}
		return "SPECs"
	}
}

func calculateTeamEquipmentValue(game *Game, team *common.TeamState, p dem.Parser) int {
	equipment := 0
	for _, teamMember := range getTeamMembers(team, game, p) {
		if teamMember.IsAlive() && game.PotentialRound.PlayerStats[teamMember.SteamID64].Health > 0 {
			equipment += teamMember.EquipmentValueCurrent()
		}
	}
	return equipment
}

// works for grenades, needs to be modified for other types
func calculateTeamEquipmentNum(game *Game, team *common.TeamState, equipmentENUM int, p dem.Parser) int {
	equipment := 0
	for _, teamMember := range getTeamMembers(team, game, p) {
		if teamMember.IsAlive() && game.PotentialRound.PlayerStats[teamMember.SteamID64].Health > 0 {
			//log.Debug(teamMember.Inventory)
			//log.Debug(teamMember.Weapons())
			//log.Debug(teamMember.AmmoLeft)
			//gren := teamMember.Inventory[equipmentENUM]
			equipment += teamMember.AmmoLeft[equipmentENUM]
		}
	}
	return equipment
}

func closestCTDisttoBomb(game *Game, team *common.TeamState, bomb *common.Bomb, p dem.Parser) int {
	var distance = 999999
	for _, teamMember := range getTeamMembers(team, game, p) {
		if teamMember.IsAlive() && game.PotentialRound.PlayerStats[teamMember.SteamID64].Health > 0 {
			if bomb.Position().Distance(teamMember.Position()) < float64(distance) {
				distance = int(bomb.Position().Distance(teamMember.Position()))
			}
		}
	}
	return distance
}

func numOfKits(game *Game, team *common.TeamState, p dem.Parser) int {
	kits := 0
	for _, teamMember := range getTeamMembers(team, game, p) {
		if teamMember.IsAlive() && game.PotentialRound.PlayerStats[teamMember.SteamID64].Health > 0 {
			if teamMember.HasDefuseKit() {
				kits += 1
			}
		}
	}
	return kits
}

func playersWithArmor(game *Game, team *common.TeamState, p dem.Parser) int {
	armor := 0
	for _, teamMember := range getTeamMembers(team, game, p) {
		if teamMember.IsAlive() && game.PotentialRound.PlayerStats[teamMember.SteamID64].Health > 0 {
			if teamMember.Armor() > 0 {
				armor += 1
			}
		}
	}
	return armor
}

var roundEndReasons = map[int]string{
	0:  "StillInProgress", //base values
	1:  "TargetBombed",
	2:  "VIPEscaped",
	3:  "VIPKilled",
	4:  "TerroristsEscaped",
	5:  "CTStoppedEscape",
	6:  "TerroristsStopped",
	7:  "BombDefused",
	8:  "CTWin",
	9:  "TerroristsWin",
	10: "Draw",
	11: "HostagesRescued",
	12: "TargetSaved",
	13: "HostagesNotRescued",
	14: "TerroristsNotEscaped",
	15: "VIPNotEscaped",
	16: "GameStart",
	17: "TerroristsSurrender",
	18: "CTSurrender",
	19: "TerroristsPlanted",
	20: "CTsReachedHostage",
}

func getPlayerAPIDict(side string, player *playerStats) Dictionary {

	return Dictionary{
		"playerSteamId": player.SteamID,
		"side":          side,
		"teamName":      player.TeamClanName,
		"adp":           player.DeathPlacement,
		"adr":           player.Adr,
		"assists":       player.Assists,
		"atd":           player.Atd,
		"awpK":          player.AwpKills,
		"damageDealt":   player.Damage,
		"damageTaken":   player.DamageTaken,
		"deaths":        player.Deaths,
		"eac":           player.Eac,
		"ef":            player.Ef,
		"eft":           player.EnemyFlashTime,
		"fAss":          player.FAss,
		"fDeaths":       player.Ol,
		"fireDamage":    player.InfernoDmg,
		"fires":         player.FiresThrown,
		"fiveK":         player.FiveK,
		"fourK":         player.FourK,
		"threeK":        player.ThreeK,
		"twoK":          player.TwoK,
		"fKills":        player.Ok,
		"flashes":       player.FlashThrown,
		"hs":            player.Hs,
		"impact":        player.ImpactRating,
		"iwr":           player.Iiwr,
		"jumps":         0,
		"kast":          player.Kast,
		"kills":         player.Kills,
		"kpa":           player.KillPointAvg,
		"lurks":         player.LurkRounds,
		"mip":           player.Mip,
		"nadeDamage":    player.NadeDmg,
		"nades":         player.NadesThrown,
		"oneVFive":      player.Cl_5,
		"oneVFour":      player.Cl_4,
		"oneVThree":     player.Cl_3,
		"oneVTwo":       player.Cl_2,
		"oneVOne":       player.Cl_1,
		"ra":            player.RA,
		"rating":        player.Rating,
		"rf":            player.RF,
		"rounds":        player.Rounds,
		"rwk":           player.Rwk,
		"rws":           player.Rws,
		"saves":         player.Saves,
		"smokes":        player.SmokeThrown,
		"suppR":         player.SuppRounds,
		"suppX":         player.SuppDamage,
		"traded":        player.Traded,
		"trades":        player.Trades,
		"ud":            player.UtilDmg,
		"util":          player.UtilThrown,
		"wlp":           player.Wlp,
	}
}
