package types

type Accolades struct {
	Awp        int `json:"awp"`
	Deagle     int `json:"deagle"`
	Knife      int `json:"knife"`
	Dinks      int `json:"dinks"`
	BlindKills int `json:"blindKills"`
	BombPlants int `json:"bombPlants"`
	Jumps      int `json:"jumps"`
	TeamDMG    int `json:"teamDMG"`
	SelfDMG    int `json:"selfDMG"`
	Ping       int `json:"ping"`
	PingPoints int `json:"pingPoints"`
	//footsteps         int //unnecessary processing?
	BombTaps          int `json:"bombTaps"`
	KillsThroughSmoke int `json:"killsThroughSmoke"`
	Penetrations      int `json:"penetrations"`
	NoScopes          int `json:"noScopes"`
	MidairKills       int `json:"midairKills"`
	CrouchedKills     int `json:"crouchedKills"`
	BombzoneKills     int `json:"bombzoneKills"`
	KillsWhileMoving  int `json:"killsWhileMoving"`
	MostMoneySpent    int `json:"mostMoneySpent"`
	MostShotsOnLegs   int `json:"mostShotsOnLegs"`
	ShotsFired        int `json:"shotsFired"`
	Ak                int `json:"ak"`
	M4                int `json:"m4"`
	Pistol            int `json:"pistol"`
	Scout             int `json:"scout"`
}
