package types

type Flag struct {
	//all our sentinals and shit
	HasGameStarted            bool `json:"hasGameStarted"`
	IsGameLive                bool `json:"isGameLive"`
	IsGameOver                bool `json:"isGameOver"`
	InRound                   bool `json:"inRound"`
	PrePlant                  bool `json:"prePlant"`
	PostPlant                 bool `json:"postPlant"`
	PostWinCon                bool `json:"postWinCon"`
	RoundIntegrityStart       int  `json:"roundIntegrityStart"`
	RoundIntegrityEnd         int  `json:"roundIntegrityEnd"`
	RoundIntegrityEndOfficial int  `json:"roundIntegrityEndOfficial"`

	//for the round (gets reset on a new round) maybe should be in a new struct
	TAlive            int    `json:"TAlive"`
	CtAlive           int    `json:"ctAlive"`
	TMoney            bool   `json:"TMoney"`
	TClutchVal        int    `json:"TClutchVal"`
	CtClutchVal       int    `json:"ctClutchVal"`
	TClutchSteam      uint64 `json:"TClutchSteam"`
	CtClutchSteam     uint64 `json:"ctClutchSteam"`
	OpeningKill       bool   `json:"openingKill"`
	LastTickProcessed int    `json:"lastTickProcessed"`
	TicksProcessed    int    `json:"ticksProcessed"`
	DidRoundEndFire   bool   `json:"didRoundEndFire"`
	RoundStartedAt    int    `json:"roundStartedAt"`
	HaveInitRound     bool   `json:"haveInitRound"`
}
