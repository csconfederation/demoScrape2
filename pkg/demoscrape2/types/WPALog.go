package types

type WPALog struct {
	Round               int `json:"round"`
	Tick                int `json:"tick"`
	Clock               int `json:"clock"`
	Planted             int `json:"planted"`
	CtAlive             int `json:"ctAlive"`
	TAlive              int `json:"TAlive"`
	CtEquipVal          int `json:"ctEquipVal"`
	TEquipVal           int `json:"TEquipVal"`
	CtFlashes           int `json:"ctFlashes"`
	CtSmokes            int `json:"ctSmokes"`
	CtMolys             int `json:"ctMolys"`
	CtFrags             int `json:"ctFrags"`
	TFlashes            int `json:"TFlashes"`
	TSmokes             int `json:"TSmokes"`
	TMolys              int `json:"TMolys"`
	TFrags              int `json:"TFrags"`
	ClosestCTDisttoBomb int `json:"closestCTDisttoBomb"`
	Kits                int `json:"kits"`
	CtArmor             int `json:"ctArmor"`
	TArmor              int `json:"TArmor"`
	Winner              int `json:"winner"`
}
