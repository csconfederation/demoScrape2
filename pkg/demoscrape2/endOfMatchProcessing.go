package demoscrape2

import (
	"math"
	"strconv"

	log "github.com/sirupsen/logrus"
)

func removeInvalidRounds(game *Game) {
	//we want to remove bad rounds (knife/veto rounds, incomplete rounds, redo rounds)
	validRoundsMap := make(map[int8]bool)
	validRounds := make([]*round, 0)
	for i := len(game.Rounds) - 1; i >= 0; i-- {
		_, validRoundExists := validRoundsMap[game.Rounds[i].RoundNum]
		if game.Rounds[i].IntegrityCheck && !game.Rounds[i].KnifeRound && !validRoundExists {
			//this i-th round is good to add
			validRoundsMap[game.Rounds[i].RoundNum] = true
			validRounds = append(validRounds, game.Rounds[i])
		} else {
			//this i-th round is bad and we will remove it
		}
	}
	for i, j := 0, len(validRounds)-1; i < j; i, j = i+1, j-1 {
		validRounds[i], validRounds[j] = validRounds[j], validRounds[i]
	}
	game.Rounds = validRounds
}

func endOfMatchProcessing(game *Game) {
	removeInvalidRounds(game)

	game.TotalPlayerStats = make(map[uint64]*playerStats)
	game.TotalTeamStats = make(map[string]*teamStats)
	game.TotalWPAlog = make([]*wpalog, 0)

	for i := len(game.Rounds) - 1; i >= 0; i-- {
		game.TotalWPAlog = append(game.TotalWPAlog, game.Rounds[i].WPAlog...)
		game.Rounds[i].ServerNormalizer += game.Rounds[i].InitTerroristCount + game.Rounds[i].InitCTerroristCount

		for teamName, team := range game.Rounds[i].TeamStats {
			if game.TotalTeamStats[teamName] == nil && teamName != "" {
				game.TotalTeamStats[teamName] = &teamStats{}
			}
			game.TotalTeamStats[teamName].Pistols += team.Pistols
			game.TotalTeamStats[teamName].PistolsW += team.PistolsW
			game.TotalTeamStats[teamName].FourVFiveS += team.FourVFiveS
			game.TotalTeamStats[teamName].FourVFiveW += team.FourVFiveW
			game.TotalTeamStats[teamName].FiveVFourS += team.FiveVFourS
			game.TotalTeamStats[teamName].FiveVFourW += team.FiveVFourW
			game.TotalTeamStats[teamName].Saves += team.Saves
			game.TotalTeamStats[teamName].Clutches += team.Clutches
			game.TotalTeamStats[teamName].CtR += team.CtR
			game.TotalTeamStats[teamName].CtRW += team.CtRW
			game.TotalTeamStats[teamName].TR += team.TR
			game.TotalTeamStats[teamName].TRW += team.TRW
		}

		//add to round master stats
		log.Debug(game.Rounds[i].RoundNum)
		for steam, player := range (*game.Rounds[i]).PlayerStats {
			if game.TotalPlayerStats[steam] == nil {
				game.TotalPlayerStats[steam] = &playerStats{Name: player.Name, SteamID: player.SteamID, TeamClanName: player.TeamClanName}
			}
			game.TotalPlayerStats[steam].Rounds += 1
			game.TotalPlayerStats[steam].Kills += player.Kills
			game.TotalPlayerStats[steam].Assists += player.Assists
			game.TotalPlayerStats[steam].Deaths += player.Deaths
			game.TotalPlayerStats[steam].Damage += player.Damage
			game.TotalPlayerStats[steam].TicksAlive += player.TicksAlive
			game.TotalPlayerStats[steam].DeathPlacement += player.DeathPlacement
			game.TotalPlayerStats[steam].Trades += player.Trades
			game.TotalPlayerStats[steam].Traded += player.Traded
			game.TotalPlayerStats[steam].Ok += player.Ok
			game.TotalPlayerStats[steam].Ol += player.Ol
			game.TotalPlayerStats[steam].KillPoints += player.KillPoints
			game.TotalPlayerStats[steam].Cl_1 += player.Cl_1
			game.TotalPlayerStats[steam].Cl_2 += player.Cl_2
			game.TotalPlayerStats[steam].Cl_3 += player.Cl_3
			game.TotalPlayerStats[steam].Cl_4 += player.Cl_4
			game.TotalPlayerStats[steam].Cl_5 += player.Cl_5
			game.TotalPlayerStats[steam].TwoK += player.TwoK
			game.TotalPlayerStats[steam].ThreeK += player.ThreeK
			game.TotalPlayerStats[steam].FourK += player.FourK
			game.TotalPlayerStats[steam].FiveK += player.FiveK
			game.TotalPlayerStats[steam].NadeDmg += player.NadeDmg
			game.TotalPlayerStats[steam].InfernoDmg += player.InfernoDmg
			game.TotalPlayerStats[steam].UtilDmg += player.UtilDmg
			game.TotalPlayerStats[steam].Ef += player.Ef
			game.TotalPlayerStats[steam].FAss += player.FAss
			game.TotalPlayerStats[steam].EnemyFlashTime += player.EnemyFlashTime
			game.TotalPlayerStats[steam].Hs += player.Hs
			game.TotalPlayerStats[steam].KastRounds += player.KastRounds
			game.TotalPlayerStats[steam].Saves += player.Saves
			game.TotalPlayerStats[steam].Entries += player.Entries
			game.TotalPlayerStats[steam].ImpactPoints += player.ImpactPoints
			game.TotalPlayerStats[steam].WinPoints += player.WinPoints
			game.TotalPlayerStats[steam].AwpKills += player.AwpKills
			game.TotalPlayerStats[steam].RF += player.RF
			game.TotalPlayerStats[steam].RA += player.RA
			game.TotalPlayerStats[steam].NadesThrown += player.NadesThrown
			game.TotalPlayerStats[steam].SmokeThrown += player.SmokeThrown
			game.TotalPlayerStats[steam].FlashThrown += player.FlashThrown
			game.TotalPlayerStats[steam].FiresThrown += player.FiresThrown
			game.TotalPlayerStats[steam].DamageTaken += player.DamageTaken
			game.TotalPlayerStats[steam].SuppDamage += player.SuppDamage
			game.TotalPlayerStats[steam].SuppRounds += player.SuppRounds
			game.TotalPlayerStats[steam].Rwk += player.Rwk
			game.TotalPlayerStats[steam].Mip += player.Mip
			game.TotalPlayerStats[steam].Eac += player.Eac
			game.TotalPlayerStats[steam].Side = 4

			if player.IsBot {
				game.TotalPlayerStats[steam].IsBot = true
			}

			if player.RF == 1 {
				game.Rounds[i].WinTeamDmg += player.Damage
			}

			if player.Side == 2 {
				game.TotalPlayerStats[steam].WinPointsNormalizer += game.Rounds[i].InitTerroristCount
				game.TotalPlayerStats[steam].TImpactPoints += player.ImpactPoints
				game.TotalPlayerStats[steam].TWinPoints += player.WinPoints
				game.TotalPlayerStats[steam].TOK += player.Ok
				game.TotalPlayerStats[steam].TOL += player.Ol
				game.TotalPlayerStats[steam].TKills += player.Kills
				game.TotalPlayerStats[steam].TDeaths += player.Deaths
				game.TotalPlayerStats[steam].TKASTRounds += player.KastRounds
				game.TotalPlayerStats[steam].TDamage += player.Damage
				game.TotalPlayerStats[steam].TADP += player.DeathPlacement
				//Game.TotalPlayerStats[steam].tTeamsWinPoints +=
				game.TotalPlayerStats[steam].TWinPointsNormalizer += game.Rounds[i].InitTerroristCount
				game.TotalPlayerStats[steam].TRounds += 1
				game.TotalPlayerStats[steam].TRF += player.RF
				game.TotalPlayerStats[steam].LurkRounds += player.LurkRounds
				if player.LurkRounds != 0 {
					game.TotalPlayerStats[steam].Wlp += player.WinPoints
				}

				game.Rounds[i].TeamStats[player.TeamClanName].TWinPoints += player.WinPoints
				game.Rounds[i].TeamStats[player.TeamClanName].TImpactPoints += player.ImpactPoints
			} else if player.Side == 3 {
				game.TotalPlayerStats[steam].WinPointsNormalizer += game.Rounds[i].InitCTerroristCount
				game.TotalPlayerStats[steam].CtImpactPoints += player.ImpactPoints
				game.TotalPlayerStats[steam].CtWinPoints += player.WinPoints
				game.TotalPlayerStats[steam].CtOK += player.Ok
				game.TotalPlayerStats[steam].CtOL += player.Ol
				game.TotalPlayerStats[steam].CtKills += player.Kills
				game.TotalPlayerStats[steam].CtDeaths += player.Deaths
				game.TotalPlayerStats[steam].CtKASTRounds += player.KastRounds
				game.TotalPlayerStats[steam].CtDamage += player.Damage
				game.TotalPlayerStats[steam].CtADP += player.DeathPlacement
				//Game.TotalPlayerStats[steam].tTeamsWinPoints +=
				game.TotalPlayerStats[steam].CtWinPointsNormalizer += game.Rounds[i].InitCTerroristCount
				game.TotalPlayerStats[steam].CtRounds += 1
				game.TotalPlayerStats[steam].CtAWP += player.CtAWP

				game.Rounds[i].TeamStats[player.TeamClanName].CtWinPoints += player.WinPoints
				game.Rounds[i].TeamStats[player.TeamClanName].CtImpactPoints += player.ImpactPoints
			}

			game.Rounds[i].TeamStats[player.TeamClanName].WinPoints += player.WinPoints
			game.Rounds[i].TeamStats[player.TeamClanName].ImpactPoints += player.ImpactPoints

		}
		for steam, player := range (*game.Rounds[i]).PlayerStats {
			game.TotalPlayerStats[steam].TeamsWinPoints += game.Rounds[i].TeamStats[player.TeamClanName].WinPoints
			game.TotalPlayerStats[steam].TTeamsWinPoints += game.Rounds[i].TeamStats[player.TeamClanName].TWinPoints
			game.TotalPlayerStats[steam].CtTeamsWinPoints += game.Rounds[i].TeamStats[player.TeamClanName].CtWinPoints

			//give players rws
			if player.RF != 0 {
				if game.Rounds[i].EndDueToBombEvent {
					player.Rws = 70 * (float64(player.Damage) / float64(game.Rounds[i].WinTeamDmg))
					steamId64, _ := strconv.ParseUint(player.SteamID, 10, 64)
					if player.Side == 2 && game.Rounds[i].Planter == steamId64 {
						player.Rws += 30
					} else if player.Side == 3 && game.Rounds[i].Defuser == steamId64 {
						player.Rws += 30
					}
				} else { //round ended due to damage/time
					player.Rws = 100 * (float64(player.Damage) / float64(game.Rounds[i].WinTeamDmg))
				}
				if math.IsNaN(player.Rws) {
					player.Rws = 0.0
				}
				game.TotalPlayerStats[steam].Rws += player.Rws
			}
		}
	}

	for _, player := range game.TotalPlayerStats {
		game.TotalTeamStats[player.TeamClanName].Util += player.SmokeThrown + player.FlashThrown + player.NadesThrown + player.FiresThrown
		game.TotalTeamStats[player.TeamClanName].Ud += player.UtilDmg
		game.TotalTeamStats[player.TeamClanName].Ef += player.Ef
		game.TotalTeamStats[player.TeamClanName].Fass += player.FAss
		game.TotalTeamStats[player.TeamClanName].Traded += player.Traded
		game.TotalTeamStats[player.TeamClanName].Deaths += int(player.Deaths)
	}

	calculateDerivedFields(game)
	return
}

func calculateDerivedFields(game *Game) {

	impactRoundAvg := 0.0
	killRoundAvg := 0.0
	deathRoundAvg := 0.0
	kastRoundAvg := 0.0
	adrAvg := 0.0
	roundNormalizer := 0

	tImpactRoundAvg := 0.0
	tKillRoundAvg := 0.0
	tDeathRoundAvg := 0.0
	tKastRoundAvg := 0.0
	tAdrAvg := 0.0
	tRoundNormalizer := 0

	ctImpactRoundAvg := 0.0
	ctKillRoundAvg := 0.0
	ctDeathRoundAvg := 0.0
	ctKastRoundAvg := 0.0
	ctAdrAvg := 0.0
	ctRoundNormalizer := 0

	//check our shit
	for _, player := range game.TotalPlayerStats {

		player.Atd = player.TicksAlive / player.Rounds / game.TickRate
		player.DeathPlacement = player.DeathPlacement / float64(player.Deaths)
		player.Kast = player.KastRounds / float64(player.Rounds)
		player.KillPointAvg = player.KillPoints / float64(player.Kills)
		if player.Kills == 0 {
			player.KillPointAvg = 0
		}
		player.Iiwr = player.WinPoints / player.ImpactPoints
		player.Adr = float64(player.Damage) / float64(player.Rounds)
		player.DrDiff = player.Adr - (float64(player.DamageTaken) / float64(player.Rounds))
		player.Tr = float64(player.Traded) / float64(player.Deaths)
		player.KR = float64(player.Kills) / float64(player.Rounds)
		player.UtilThrown = player.SmokeThrown + player.FlashThrown + player.NadesThrown + player.FiresThrown
		player.CtADR = float64(player.CtDamage) / float64(player.CtRounds)
		player.TADR = float64(player.TDamage) / float64(player.TRounds)
		player.TKAST = player.TKASTRounds / float64(player.TRounds)
		player.CtKAST = player.CtKASTRounds / float64(player.CtRounds)
		player.TADP = player.TADP / float64(player.TDeaths)
		player.CtADP = player.CtADP / float64(player.CtDeaths)
		player.Rws = player.Rws / float64(player.Rounds)

		if math.IsNaN(player.Rws) {
			player.Rws = 0.0
		}
		if player.ImpactPoints == 0 {
			player.Iiwr = 0
		}
		if player.Deaths == 0 {
			player.DeathPlacement = 0
			player.Tr = .50
		}
		if player.TDeaths == 0 {
			player.TADP = 0
		}
		if player.CtDeaths == 0 {
			player.CtADP = 0
		}

		roundNormalizer += player.Rounds
		impactRoundAvg += player.ImpactPoints
		killRoundAvg += float64(player.Kills)
		deathRoundAvg += float64(player.Deaths)
		kastRoundAvg += player.KastRounds
		adrAvg += float64(player.Damage)

		tImpactRoundAvg += player.TImpactPoints
		tKillRoundAvg += float64(player.TKills)
		tDeathRoundAvg += float64(player.CtDeaths)
		tKastRoundAvg += player.TKASTRounds
		tAdrAvg += float64(player.TDamage)
		tRoundNormalizer += player.TRounds

		ctImpactRoundAvg += player.CtImpactPoints
		ctKillRoundAvg += float64(player.CtKills)
		ctDeathRoundAvg += float64(player.TDeaths)
		ctKastRoundAvg += player.CtKASTRounds
		ctAdrAvg += float64(player.CtDamage)
		ctRoundNormalizer += player.CtRounds

	}

	impactRoundAvg /= float64(roundNormalizer)
	killRoundAvg /= float64(roundNormalizer)
	deathRoundAvg /= float64(roundNormalizer)
	kastRoundAvg /= float64(roundNormalizer)
	adrAvg /= float64(roundNormalizer)

	tImpactRoundAvg /= float64(tRoundNormalizer)
	tKillRoundAvg /= float64(tRoundNormalizer)
	tDeathRoundAvg /= float64(tRoundNormalizer)
	tKastRoundAvg /= float64(tRoundNormalizer)
	tAdrAvg /= float64(tRoundNormalizer)

	ctImpactRoundAvg /= float64(ctRoundNormalizer)
	ctKillRoundAvg /= float64(ctRoundNormalizer)
	ctDeathRoundAvg /= float64(ctRoundNormalizer)
	ctKastRoundAvg /= float64(ctRoundNormalizer)
	ctAdrAvg /= float64(ctRoundNormalizer)

	for _, player := range game.TotalPlayerStats {
		openingFactor := (float64(player.Ok-player.Ol) / 13.0) + 1 //move from 13 to (rounds / 5)
		playerIPR := player.ImpactPoints / float64(player.Rounds)
		playerWPR := player.WinPoints / float64(player.Rounds)

		if player.TeamsWinPoints != 0 {
			player.ImpactRating = (0.1 * float64(openingFactor)) + (0.6 * (playerIPR / impactRoundAvg)) + (0.3 * (playerWPR / (player.TeamsWinPoints / float64(player.WinPointsNormalizer))))
		} else {
			log.Debug("UH 16-0?")
			player.ImpactRating = (0.1 * float64(openingFactor)) + (0.6 * (playerIPR / impactRoundAvg))
		}
		playerDR := float64(player.Deaths) / float64(player.Rounds)
		playerRatingDeathComponent := 0.07 * (deathRoundAvg / playerDR)
		if player.Deaths == 0 || playerRatingDeathComponent > 0.21 {
			playerRatingDeathComponent = 0.21
		}
		player.Rating = (0.3 * player.ImpactRating) + (0.35 * (player.KR / killRoundAvg)) + playerRatingDeathComponent + (0.08 * (player.Kast / kastRoundAvg)) + (0.2 * (player.Adr / adrAvg))
		game.TotalTeamStats[player.TeamClanName].RatingAvg += player.Rating
		game.TotalTeamStats[player.TeamClanName].Normalizer += 1

		//ctRating
		openingFactor = (float64(player.CtOK-player.CtOL) / 13.0) + 1
		playerIPR = player.CtImpactPoints / float64(player.CtRounds)
		playerWPR = player.CtWinPoints / float64(player.CtRounds)

		if player.CtTeamsWinPoints != 0 {
			player.CtImpactRating = (0.1 * float64(openingFactor)) + (0.6 * (playerIPR / ctImpactRoundAvg)) + (0.3 * (playerWPR / (player.CtTeamsWinPoints / float64(player.CtWinPointsNormalizer))))
		} else {
			log.Debug("UH 16-0?")
			player.CtImpactRating = (0.1 * float64(openingFactor)) + (0.6 * (playerIPR / ctImpactRoundAvg))
		}
		playerDR = float64(player.CtDeaths) / float64(player.CtRounds)
		playerRatingDeathComponent = 0.07 * (ctDeathRoundAvg / playerDR)
		if player.CtDeaths == 0 || playerRatingDeathComponent > 0.21 {
			playerRatingDeathComponent = 0.21
		}
		player.CtRating = (0.3 * player.CtImpactRating) + (0.35 * ((float64(player.CtKills) / float64(player.CtRounds)) / ctKillRoundAvg)) + playerRatingDeathComponent + (0.08 * (player.CtKAST / ctKastRoundAvg)) + (0.2 * (player.CtADR / ctAdrAvg))

		//tRating
		openingFactor = (float64(player.TOK-player.TOL) / 13.0) + 1
		playerIPR = player.TImpactPoints / float64(player.TRounds)
		playerWPR = player.TWinPoints / float64(player.TRounds)

		if player.TTeamsWinPoints != 0 {
			player.TImpactRating = (0.1 * float64(openingFactor)) + (0.6 * (playerIPR / tImpactRoundAvg)) + (0.3 * (playerWPR / (player.TTeamsWinPoints / float64(player.TWinPointsNormalizer))))
		} else {
			log.Debug("UH 16-0?")
			player.TImpactRating = (0.1 * float64(openingFactor)) + (0.6 * (playerIPR / tImpactRoundAvg))
		}
		playerDR = float64(player.TDeaths) / float64(player.TRounds)
		playerRatingDeathComponent = 0.07 * (tDeathRoundAvg / playerDR)
		if player.TDeaths == 0 || playerRatingDeathComponent > 0.21 {
			playerRatingDeathComponent = 0.21
		}
		player.TRating = (0.3 * player.TImpactRating) + (0.35 * ((float64(player.TKills) / float64(player.TRounds)) / tKillRoundAvg)) + playerRatingDeathComponent + (0.08 * (player.TKAST / tKastRoundAvg)) + (0.2 * (player.TADR / tAdrAvg))

		log.Debug("openingFactor", 0.1*float64(openingFactor))
		log.Debug("playerIPR", 0.6*(playerIPR/impactRoundAvg))
		log.Debug("playerWPR", 0.3*(playerWPR/(player.TeamsWinPoints/float64(player.WinPointsNormalizer))))
		log.Debug("player.teamsWinPoints", player.TeamsWinPoints)
		log.Debug("player.winPointsNormalizer", player.WinPointsNormalizer)

		log.Debug("%+v\n\n", player)
	}
	log.Debug("impactRoundAvg", impactRoundAvg)
	log.Debug("killRoundAvg", killRoundAvg)
	log.Debug("deathRoundAvg", deathRoundAvg)
	log.Debug("kastRoundAvg", kastRoundAvg)
	log.Debug("adrAvg", adrAvg)

	for _, player := range game.TotalPlayerStats {
		players := float64(game.TotalTeamStats[player.TeamClanName].Normalizer)
		teamAvg := (game.TotalTeamStats[player.TeamClanName].RatingAvg / players)
		adjustedAvg := (players*teamAvg - player.Rating) / (players - 1)
		player.DeltaRating = player.Rating - adjustedAvg
	}

	for _, team := range game.TotalTeamStats {
		team.RatingAvg = team.RatingAvg / float64(team.Normalizer)
	}

	calculateSidedStats(game)
	return
}

func calculateSidedStats(game *Game) {

	game.CtPlayerStats = make(map[uint64]*playerStats)
	game.TPlayerStats = make(map[uint64]*playerStats)

	for i := len(game.Rounds) - 1; i >= 0; i-- {
		game.Rounds[i].ServerNormalizer += game.Rounds[i].InitTerroristCount + game.Rounds[i].InitCTerroristCount

		//add to round master stats
		for steam, player := range (*game.Rounds[i]).PlayerStats {
			//sidedStats := make(map[uint64]*playerStats)
			sidedStats := game.CtPlayerStats
			if player.Side == 2 {
				sidedStats = game.TPlayerStats
			}
			if sidedStats[steam] == nil {
				sidedStats[steam] = &playerStats{Name: player.Name, SteamID: player.SteamID, TeamClanName: player.TeamClanName}
			}
			sidedStats[steam].Rounds += 1
			sidedStats[steam].Kills += player.Kills
			sidedStats[steam].Assists += player.Assists
			sidedStats[steam].Deaths += player.Deaths
			sidedStats[steam].Damage += player.Damage
			sidedStats[steam].TicksAlive += player.TicksAlive
			sidedStats[steam].DeathPlacement += player.DeathPlacement
			sidedStats[steam].Trades += player.Trades
			sidedStats[steam].Traded += player.Traded
			sidedStats[steam].Ok += player.Ok
			sidedStats[steam].Ol += player.Ol
			sidedStats[steam].KillPoints += player.KillPoints
			sidedStats[steam].Cl_1 += player.Cl_1
			sidedStats[steam].Cl_2 += player.Cl_2
			sidedStats[steam].Cl_3 += player.Cl_3
			sidedStats[steam].Cl_4 += player.Cl_4
			sidedStats[steam].Cl_5 += player.Cl_5
			sidedStats[steam].TwoK += player.TwoK
			sidedStats[steam].ThreeK += player.ThreeK
			sidedStats[steam].FourK += player.FourK
			sidedStats[steam].FiveK += player.FiveK
			sidedStats[steam].NadeDmg += player.NadeDmg
			sidedStats[steam].InfernoDmg += player.InfernoDmg
			sidedStats[steam].UtilDmg += player.UtilDmg
			sidedStats[steam].Ef += player.Ef
			sidedStats[steam].FAss += player.FAss
			sidedStats[steam].EnemyFlashTime += player.EnemyFlashTime
			sidedStats[steam].Hs += player.Hs
			sidedStats[steam].KastRounds += player.KastRounds
			sidedStats[steam].Saves += player.Saves
			sidedStats[steam].Entries += player.Entries
			sidedStats[steam].ImpactPoints += player.ImpactPoints
			sidedStats[steam].WinPoints += player.WinPoints
			sidedStats[steam].AwpKills += player.AwpKills
			sidedStats[steam].RF += player.RF
			sidedStats[steam].RA += player.RA
			sidedStats[steam].NadesThrown += player.NadesThrown
			sidedStats[steam].SmokeThrown += player.SmokeThrown
			sidedStats[steam].FlashThrown += player.FlashThrown
			sidedStats[steam].FiresThrown += player.FiresThrown
			sidedStats[steam].DamageTaken += player.DamageTaken
			sidedStats[steam].SuppDamage += player.SuppDamage
			sidedStats[steam].SuppRounds += player.SuppRounds
			sidedStats[steam].Rwk += player.Rwk
			sidedStats[steam].Mip += player.Mip
			sidedStats[steam].Eac += player.Eac
			sidedStats[steam].Side = player.Side

			if player.IsBot {
				sidedStats[steam].IsBot = true
			}

			sidedStats[steam].LurkRounds += player.LurkRounds
			if player.LurkRounds != 0 {
				sidedStats[steam].Wlp += player.WinPoints
			}

			if math.IsNaN(player.Rws) {
				player.Rws = 0.0
			}
			sidedStats[steam].Rws += player.Rws

			if player.Side == 2 {
				sidedStats[steam].Rating = game.TotalPlayerStats[steam].TRating
				sidedStats[steam].ImpactRating = game.TotalPlayerStats[steam].TImpactRating
			} else {
				sidedStats[steam].Rating = game.TotalPlayerStats[steam].CtRating
				sidedStats[steam].ImpactRating = game.TotalPlayerStats[steam].CtImpactRating
			}

		}
	}

	for _, player := range game.CtPlayerStats {
		player.Atd = player.TicksAlive / player.Rounds / game.TickRate
		player.DeathPlacement = player.DeathPlacement / float64(player.Deaths)
		player.Kast = player.KastRounds / float64(player.Rounds)
		player.KillPointAvg = player.KillPoints / float64(player.Kills)
		if player.Kills == 0 {
			player.KillPointAvg = 0
		}
		player.Iiwr = player.WinPoints / player.ImpactPoints
		player.Adr = float64(player.Damage) / float64(player.Rounds)
		player.DrDiff = player.Adr - (float64(player.DamageTaken) / float64(player.Rounds))
		player.Tr = float64(player.Traded) / float64(player.Deaths)
		player.KR = float64(player.Kills) / float64(player.Rounds)
		player.UtilThrown = player.SmokeThrown + player.FlashThrown + player.NadesThrown + player.FiresThrown
		player.Rws = player.Rws / float64(player.Rounds)
		if math.IsNaN(player.Rws) {
			player.Rws = 0.0
		}
		if player.ImpactPoints == 0 {
			player.Iiwr = 0
		}
		if player.Deaths == 0 {
			player.DeathPlacement = 0
			player.Tr = .50
		}
		if player.TDeaths == 0 {
			player.TADP = 0
		}
		if player.CtDeaths == 0 {
			player.CtADP = 0
		}
	}
	for _, player := range game.TPlayerStats {
		player.Atd = player.TicksAlive / player.Rounds / game.TickRate
		player.DeathPlacement = player.DeathPlacement / float64(player.Deaths)
		player.Kast = player.KastRounds / float64(player.Rounds)
		player.KillPointAvg = player.KillPoints / float64(player.Kills)
		if player.Kills == 0 {
			player.KillPointAvg = 0
		}
		player.Iiwr = player.WinPoints / player.ImpactPoints
		player.Adr = float64(player.Damage) / float64(player.Rounds)
		player.DrDiff = player.Adr - (float64(player.DamageTaken) / float64(player.Rounds))
		player.Tr = float64(player.Traded) / float64(player.Deaths)
		player.KR = float64(player.Kills) / float64(player.Rounds)
		player.UtilThrown = player.SmokeThrown + player.FlashThrown + player.NadesThrown + player.FiresThrown
		player.Rws = player.Rws / float64(player.Rounds)
		if math.IsNaN(player.Rws) {
			player.Rws = 0.0
		}
		if player.ImpactPoints == 0 {
			player.Iiwr = 0
		}
		if player.Deaths == 0 {
			player.DeathPlacement = 0
			player.Tr = .50
		}
		if player.TDeaths == 0 {
			player.TADP = 0
		}
		if player.CtDeaths == 0 {
			player.CtADP = 0
		}
	}

	return
}
