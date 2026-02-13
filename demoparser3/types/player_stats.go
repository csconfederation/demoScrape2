package types

import (
	"time"

	"github.com/csconfederation/demoparser3/logger"
	"github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs/common"
)

type PlayerStats struct {
	ImpactPoints float64 `json:"impactPoints" end_of_match_sum:"true"`
	DamageDone   int     `json:"damageDone" end_of_match_sum:"true"`
	DamageTaken  int     `json:"damageTaken" end_of_match_sum:"true"`
	//TODO: Maybe separate utility? I don't see much of a benefit other than readability
	UtilityDamage      int     `json:"utilityDamage" end_of_match_sum:"true"`
	HEDamage           int     `json:"heDamage" end_of_match_sum:"true"`
	FireDamage         int     `json:"fireDamage" end_of_match_sum:"true"`
	EnemiesFlashed     int     `json:"enemiesFlashed" end_of_match_sum:"true"`
	EnemyFlashDuration float64 `json:"enemyFlashTime" end_of_match_sum:"true"`
	//TODO: check
	//MostRecentFlasher uint64         `json:"mostRecentFlasher"`
	UtilityThrown map[string]int `json:"utilityThrown" end_of_match_sum:"true"`
	Kills         int            `json:"kills" end_of_match_sum:"true"`
	Deaths        int            `json:"deaths" end_of_match_sum:"true"`
	DeathOrder    int            `json:"deathOrder" end_of_match_sum:"true"` // Order in which a player died. Ex.: Entry fraggers will have the lowest order (possibly 0)
	Assists       int            `json:"assists" end_of_match_sum:"true"`
	EAC           int            `json:"eac" end_of_match_sum:"true"` // Effective Assists Contribution
	FlashAssists  int            `json:"flashAssists" end_of_match_sum:"true"`
	Clutch        *ClutchAttempt `json:"clutch" end_of_match_sum:"true"`
	DamageList    map[uint64]int `json:"damageList"`
	SupportDamage int            `json:"supportDamage" end_of_match_sum:"true"`
	LurkerBlips   int            `json:"lurkerBlips" end_of_match_sum:"true"`
}

func NewPlayerStats() *PlayerStats {
	return &PlayerStats{
		ImpactPoints:  0,
		UtilityThrown: make(map[string]int),
		DamageList:    make(map[uint64]int),
	}
}

func (player *PlayerStats) PlayerHurt(attacker *PlayerStats, damage int, equipmentType common.EquipmentType) {
	player.DamageTaken += damage
	attacker.DamageDone += damage
	//TODO: Check if we need this
	//player.DamageList[attacker.] += damage
	if isUtility(equipmentType) {
		attacker.UtilityDamage += damage
		switch equipmentType {
		case common.EqHE:
			attacker.HEDamage += damage
			break
		case common.EqMolotov, common.EqIncendiary:
			attacker.FireDamage += damage
			break
		}
	}
}

func (player *PlayerStats) PlayerFlashed(attacker *PlayerStats, flashDuration time.Duration) {
	//TODO: Is this required
	//player.MostRecentFlasher = attacker
	attacker.EnemiesFlashed += 1
	attacker.EnemyFlashDuration += flashDuration.Seconds()
}

func (player *PlayerStats) GrenadeThrown(grenade common.EquipmentType) {
	if !isUtility(grenade) {
		logger.Error("Not a grenade was thrown????",
			"player", player,
			"grenade", grenade)
	}

	player.UtilityThrown[grenade.String()] += 1
}

func (player *PlayerStats) Kill() {
	player.Kills += 1
}

func (player *PlayerStats) Death(deathOrder int) {
	player.Deaths += 1
	player.DeathOrder = deathOrder
}

func (player *PlayerStats) Assist(isFlashAssist bool) {
	player.Assists += 1
	if isFlashAssist {
		player.FlashAssists += 1
	}

	//TODO: How to handle this
	//pS[e.Assister.SteamID64].Eac += 1
	//player.KASTRound = 1
	//player.SupportRound = 1

	//else if float64(parser.GameState().IngameTick()) < pS[e.Victim.SteamID64].MostRecentFlashVal {
	//	//this will trigger if there is both a flash assist and a damage assist
	//	pS[pS[e.Victim.SteamID64].MostRecentFlasher].FAss += 1
	//	pS[pS[e.Victim.SteamID64].MostRecentFlasher].Eac += 1
	//	pS[pS[e.Victim.SteamID64].MostRecentFlasher].SuppRounds = 1
	//	flashAssisted = true
	//	flashAssister = pS[pS[e.Victim.SteamID64].MostRecentFlasher].Name
	//}

}

//type PlayerStats struct {
//	Name    string `json:"name"`
//	SteamID uint64 `json:"steamID"`
//	IsBot   bool   `json:"isBot"`
//	//teamID  int
//	TeamENUM     int    `json:"teamENUM"`
//	TeamClanName string `json:"teamClanName"`
//	Side         int    `json:"side"`
//	Rounds       int    `json:"rounds"`
//	//playerPoints float32
//	//teamPoints float32
//	Damage              int     `json:"damage" end_of_match_sum:"true"`
//	Kills               uint8   `json:"kills" end_of_match_sum:"true"`
//	Assists             uint8   `json:"assists" end_of_match_sum:"true"`
//	Deaths              uint8   `json:"deaths" end_of_match_sum:"true"`
//	DeathTick           int     `json:"deathTick"`
//	DeathPlacement      float64 `json:"deathPlacement" end_of_match_sum:"true"`
//	TicksAlive          int     `json:"ticksAlive" end_of_match_sum:"true"`
//	Trades              int     `json:"trades" end_of_match_sum:"true"`
//	Traded              int     `json:"traded" end_of_match_sum:"true"`
//	Ok                  int     `json:"ok" end_of_match_sum:"true"`
//	Ol                  int     `json:"ol" end_of_match_sum:"true"`
//	Cl_1                int     `json:"cl_1" end_of_match_sum:"true"`
//	Cl_2                int     `json:"cl_2" end_of_match_sum:"true"`
//	Cl_3                int     `json:"cl_3" end_of_match_sum:"true"`
//	Cl_4                int     `json:"cl_4" end_of_match_sum:"true"`
//	Cl_5                int     `json:"cl_5" end_of_match_sum:"true"`
//	TwoK                int     `json:"twoK" end_of_match_sum:"true"`
//	ThreeK              int     `json:"threeK" end_of_match_sum:"true"`
//	FourK               int     `json:"fourK" end_of_match_sum:"true"`
//	FiveK               int     `json:"fiveK" end_of_match_sum:"true"`
//	NadeDmg             int     `json:"nadeDmg" end_of_match_sum:"true"`
//	InfernoDmg          int     `json:"infernoDmg" end_of_match_sum:"true"`
//	UtilDmg             int     `json:"utilDmg" end_of_match_sum:"true"`
//	Ef                  int     `json:"ef" end_of_match_sum:"true"`
//	FAss                int     `json:"FAss" end_of_match_sum:"true"`
//	EnemyFlashTime      float64 `json:"enemyFlashTime" end_of_match_sum:"true"`
//	Hs                  int     `json:"hs" end_of_match_sum:"true"`
//	KastRounds          float64 `json:"kastRounds" end_of_match_sum:"true"`
//	Saves               int     `json:"saves" end_of_match_sum:"true"`
//	Entries             int     `json:"entries" end_of_match_sum:"true"`
//	KillPoints          float64 `json:"killPoints"`
//	ImpactPoints        float64 `json:"impactPoints" end_of_match_sum:"true"`
//	WinPoints           float64 `json:"winPoints" end_of_match_sum:"true"`
//	AwpKills            int     `json:"awpKills" end_of_match_sum:"true"`
//	RF                  int     `json:"RF" end_of_match_sum:"true"`
//	RA                  int     `json:"RA" end_of_match_sum:"true"`
//	NadesThrown         int     `json:"nadesThrown" end_of_match_sum:"true"`
//	FiresThrown         int     `json:"firesThrown" end_of_match_sum:"true"`
//	FlashThrown         int     `json:"flashThrown" end_of_match_sum:"true"`
//	SmokeThrown         int     `json:"smokeThrown" end_of_match_sum:"true"`
//	DamageTaken         int     `json:"damageTaken" end_of_match_sum:"true"`
//	SuppRounds          int     `json:"suppRounds" end_of_match_sum:"true"`
//	SuppDamage          int     `json:"suppDamage" end_of_match_sum:"true"`
//	LurkerBlips         int     `json:"lurkerBlips"`
//	DistanceToTeammates int     `json:"distanceToTeammates"`
//	LurkRounds          int     `json:"lurkRounds"`
//	Wlp                 float64 `json:"wlp"`
//	Mip                 float64 `json:"mip" end_of_match_sum:"true"`
//	Rws                 float64 `json:"rws"`                         //round win shares
//	Eac                 int     `json:"eac" end_of_match_sum:"true"` //effective assist contributions
//
//	Rwk int `json:"rwk" end_of_match_sum:"true"` //rounds with Kills
//
//	//derived
//	UtilThrown   int     `json:"utilThrown"`
//	Atd          int     `json:"atd"`
//	Kast         float64 `json:"kast"`
//	KillPointAvg float64 `json:"killPointAvg"`
//	Iiwr         float64 `json:"iiwr"`
//	Adr          float64 `json:"adr"`
//	DrDiff       float64 `json:"drDiff"`
//	KR           float64 `json:"KR"`
//	Tr           float64 `json:"tr"` //trade ratio
//	ImpactRating float64 `json:"impactRating"`
//	Rating       float64 `json:"rating"`
//
//	//side specific
//	TDamage               int     `json:"TDamage"`
//	CtDamage              int     `json:"ctDamage"`
//	TImpactPoints         float64 `json:"TImpactPoints"`
//	TWinPoints            float64 `json:"TWinPoints"`
//	TOK                   int     `json:"TOK"`
//	TOL                   int     `json:"TOL"`
//	CtImpactPoints        float64 `json:"ctImpactPoints"`
//	CtWinPoints           float64 `json:"ctWinPoints"`
//	CtOK                  int     `json:"ctOK"`
//	CtOL                  int     `json:"ctOL"`
//	TKills                uint8   `json:"TKills"`
//	TDeaths               uint8   `json:"TDeaths"`
//	TKAST                 float64 `json:"TKAST"`
//	TKASTRounds           float64 `json:"TKASTRounds"`
//	TADR                  float64 `json:"TADR"`
//	CtKills               uint8   `json:"ctKills"`
//	CtDeaths              uint8   `json:"ctDeaths"`
//	CtKAST                float64 `json:"ctKAST"`
//	CtKASTRounds          float64 `json:"ctKASTRounds"`
//	CtADR                 float64 `json:"ctADR"`
//	TTeamsWinPoints       float64 `json:"TTeamsWinPoints"`
//	CtTeamsWinPoints      float64 `json:"ctTeamsWinPoints"`
//	TWinPointsNormalizer  int     `json:"TWinPointsNormalizer"`
//	CtWinPointsNormalizer int     `json:"ctWinPointsNormalizer"`
//	TRounds               int     `json:"TRounds"`
//	CtRounds              int     `json:"ctRounds"`
//	CtRating              float64 `json:"ctRating"`
//	CtImpactRating        float64 `json:"ctImpactRating"`
//	TRating               float64 `json:"TRating"`
//	TImpactRating         float64 `json:"TImpactRating"`
//	TADP                  float64 `json:"TADP"`
//	CtADP                 float64 `json:"ctADP"`
//
//	TRF   int `json:"TRF"`
//	CtAWP int `json:"ctAWP"`
//
//	//kinda garbo
//	TeamsWinPoints      float64 `json:"teamsWinPoints"`
//	WinPointsNormalizer int     `json:"winPointsNormalizer"`
//
//	//"flags"
//	Health             int            `json:"health"`
//	TradeList          map[uint64]int `json:"tradeList"`
//	MostRecentFlasher  uint64         `json:"mostRecentFlasher"`
//	MostRecentFlashVal float64        `json:"mostRecentFlashVal"`
//	DamageList         map[uint64]int `json:"damageList"`
//}
//
//func NewPlayerStatsFromPlayer(player *common.Player, team *common.TeamState, game *Game) *PlayerStats {
//	if player == nil || player.TeamState == nil {
//		return nil
//	}
//
//	return &PlayerStats{
//		SteamID:      player.SteamID64,
//		Name:         player.Name,
//		IsBot:        player.IsBot,
//		Side:         int(team.Team()),
//		TeamENUM:     team.ID(),
//		TeamClanName: ValidateTeamName(team, game.Teams, game.PotentialRound),
//		Health:       100,
//		TradeList:    make(map[uint64]int),
//		DamageList:   make(map[uint64]int),
//	}
//}

func isUtility(equipmentType common.EquipmentType) bool {
	return equipmentType == common.EqHE || equipmentType == common.EqIncendiary ||
		equipmentType == common.EqMolotov || equipmentType == common.EqFlash ||
		equipmentType == common.EqSmoke || equipmentType == common.EqDecoy
}
