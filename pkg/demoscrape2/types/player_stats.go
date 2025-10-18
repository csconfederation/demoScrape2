package types

import (
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
)

type PlayerStats struct {
	Name    string `json:"name"`
	SteamID uint64 `json:"steamID"`
	IsBot   bool   `json:"isBot"`
	//teamID  int
	TeamENUM     int    `json:"teamENUM"`
	TeamClanName string `json:"teamClanName"`
	Side         int    `json:"side"`
	Rounds       int    `json:"rounds"`
	//playerPoints float32
	//teamPoints float32
	Damage              int     `json:"damage"`
	Kills               uint8   `json:"kills"`
	Assists             uint8   `json:"assists"`
	Deaths              uint8   `json:"deaths"`
	DeathTick           int     `json:"deathTick"`
	DeathPlacement      float64 `json:"deathPlacement"`
	TicksAlive          int     `json:"ticksAlive"`
	Trades              int     `json:"trades"`
	Traded              int     `json:"traded"`
	Ok                  int     `json:"ok"`
	Ol                  int     `json:"ol"`
	Cl_1                int     `json:"cl_1"`
	Cl_2                int     `json:"cl_2"`
	Cl_3                int     `json:"cl_3"`
	Cl_4                int     `json:"cl_4"`
	Cl_5                int     `json:"cl_5"`
	TwoK                int     `json:"twoK"`
	ThreeK              int     `json:"threeK"`
	FourK               int     `json:"fourK"`
	FiveK               int     `json:"fiveK"`
	NadeDmg             int     `json:"nadeDmg"`
	InfernoDmg          int     `json:"infernoDmg"`
	UtilDmg             int     `json:"utilDmg"`
	Ef                  int     `json:"ef"`
	FAss                int     `json:"FAss"`
	EnemyFlashTime      float64 `json:"enemyFlashTime"`
	Hs                  int     `json:"hs"`
	KastRounds          float64 `json:"kastRounds"`
	Saves               int     `json:"saves"`
	Entries             int     `json:"entries"`
	KillPoints          float64 `json:"killPoints"`
	ImpactPoints        float64 `json:"impactPoints"`
	WinPoints           float64 `json:"winPoints"`
	AwpKills            int     `json:"awpKills"`
	RF                  int     `json:"RF"`
	RA                  int     `json:"RA"`
	NadesThrown         int     `json:"nadesThrown"`
	FiresThrown         int     `json:"firesThrown"`
	FlashThrown         int     `json:"flashThrown"`
	SmokeThrown         int     `json:"smokeThrown"`
	DamageTaken         int     `json:"damageTaken"`
	SuppRounds          int     `json:"suppRounds"`
	SuppDamage          int     `json:"suppDamage"`
	LurkerBlips         int     `json:"lurkerBlips"`
	DistanceToTeammates int     `json:"distanceToTeammates"`
	LurkRounds          int     `json:"lurkRounds"`
	Wlp                 float64 `json:"wlp"`
	Mip                 float64 `json:"mip"`
	Rws                 float64 `json:"rws"` //round win shares
	Eac                 int     `json:"eac"` //effective assist contributions

	Rwk int `json:"rwk"` //rounds with Kills

	//derived
	UtilThrown   int     `json:"utilThrown"`
	Atd          int     `json:"atd"`
	Kast         float64 `json:"kast"`
	KillPointAvg float64 `json:"killPointAvg"`
	Iiwr         float64 `json:"iiwr"`
	Adr          float64 `json:"adr"`
	DrDiff       float64 `json:"drDiff"`
	KR           float64 `json:"KR"`
	Tr           float64 `json:"tr"` //trade ratio
	ImpactRating float64 `json:"impactRating"`
	Rating       float64 `json:"rating"`

	//side specific
	TDamage               int     `json:"TDamage"`
	CtDamage              int     `json:"ctDamage"`
	TImpactPoints         float64 `json:"TImpactPoints"`
	TWinPoints            float64 `json:"TWinPoints"`
	TOK                   int     `json:"TOK"`
	TOL                   int     `json:"TOL"`
	CtImpactPoints        float64 `json:"ctImpactPoints"`
	CtWinPoints           float64 `json:"ctWinPoints"`
	CtOK                  int     `json:"ctOK"`
	CtOL                  int     `json:"ctOL"`
	TKills                uint8   `json:"TKills"`
	TDeaths               uint8   `json:"TDeaths"`
	TKAST                 float64 `json:"TKAST"`
	TKASTRounds           float64 `json:"TKASTRounds"`
	TADR                  float64 `json:"TADR"`
	CtKills               uint8   `json:"ctKills"`
	CtDeaths              uint8   `json:"ctDeaths"`
	CtKAST                float64 `json:"ctKAST"`
	CtKASTRounds          float64 `json:"ctKASTRounds"`
	CtADR                 float64 `json:"ctADR"`
	TTeamsWinPoints       float64 `json:"TTeamsWinPoints"`
	CtTeamsWinPoints      float64 `json:"ctTeamsWinPoints"`
	TWinPointsNormalizer  int     `json:"TWinPointsNormalizer"`
	CtWinPointsNormalizer int     `json:"ctWinPointsNormalizer"`
	TRounds               int     `json:"TRounds"`
	CtRounds              int     `json:"ctRounds"`
	CtRating              float64 `json:"ctRating"`
	CtImpactRating        float64 `json:"ctImpactRating"`
	TRating               float64 `json:"TRating"`
	TImpactRating         float64 `json:"TImpactRating"`
	TADP                  float64 `json:"TADP"`
	CtADP                 float64 `json:"ctADP"`

	TRF   int `json:"TRF"`
	CtAWP int `json:"ctAWP"`

	//kinda garbo
	TeamsWinPoints      float64 `json:"teamsWinPoints"`
	WinPointsNormalizer int     `json:"winPointsNormalizer"`

	//"flags"
	Health             int            `json:"health"`
	TradeList          map[uint64]int `json:"tradeList"`
	MostRecentFlasher  uint64         `json:"mostRecentFlasher"`
	MostRecentFlashVal float64        `json:"mostRecentFlashVal"`
	DamageList         map[uint64]int `json:"damageList"`
}

func NewPlayerStatsFromPlayer(player *common.Player, team *common.TeamState, game *Game) *PlayerStats {
	if player == nil || player.TeamState == nil {
		return nil
	}

	return &PlayerStats{
		SteamID:      player.SteamID64,
		Name:         player.Name,
		IsBot:        player.IsBot,
		Side:         int(team.Team()),
		TeamENUM:     team.ID(),
		TeamClanName: ValidateTeamName(team, game.Teams, game.PotentialRound),
		Health:       100,
		TradeList:    make(map[uint64]int),
		DamageList:   make(map[uint64]int),
	}
}
