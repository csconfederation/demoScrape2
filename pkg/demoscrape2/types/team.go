package types

import (
	"fmt"
	"strings"

	dem "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
)

type Team struct {
	//id    int //meaningless?
	Name          string `json:"name"`
	Score         int    `json:"score"`
	ScoreAdjusted int    `json:"scoreAdjusted"`
}

func ValidateTeamName(team *common.TeamState, allTeams map[string]*Team, round *Round) string {
	teamName := team.ClanName()
	teamNum := team.ID()
	if teamName != "" {
		name := ""
		if strings.HasPrefix(teamName, "[") {
			if len(teamName) == 31 {
				//name here will be truncated
				name = strings.Split(teamName, "] ")[1]
				for _, team := range allTeams {
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
		swap := false
		roundNum := round.RoundNum
		if roundNum >= MR+1 && roundNum <= (MR*2)+3 {
			swap = true
		} else if roundNum >= (MR*2)+4 {
			//we are now in OT hell :)
			if (roundNum-((MR*2)+4))/6%2 != 0 {
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

func NewTeamFromTeamState(team *common.TeamState, allTeams map[string]*Team, round *Round) *Team {
	return &Team{
		Name: ValidateTeamName(team, allTeams, round),
	}
}

func GetTeamMembers(team *common.TeamState, game *Game, p dem.Parser) []*common.Player {
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
