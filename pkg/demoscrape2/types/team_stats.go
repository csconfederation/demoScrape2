package types

type TeamStats struct {
	WinPoints      float64 `json:"winPoints"`
	ImpactPoints   float64 `json:"impactPoints"`
	TWinPoints     float64 `json:"TWinPoints"`
	CtWinPoints    float64 `json:"ctWinPoints"`
	TImpactPoints  float64 `json:"TImpactPoints"`
	CtImpactPoints float64 `json:"ctImpactPoints"`
	FourVFiveW     int     `json:"fourVFiveW"`
	FourVFiveS     int     `json:"fourVFiveS"`
	FiveVFourW     int     `json:"fiveVFourW"`
	FiveVFourS     int     `json:"fiveVFourS"`
	Pistols        int     `json:"pistols"`
	PistolsW       int     `json:"pistolsW"`
	Saves          int     `json:"saves"`
	Clutches       int     `json:"clutches"`
	Traded         int     `json:"traded"`
	Fass           int     `json:"fass"`
	Ef             int     `json:"ef"`
	Ud             int     `json:"ud"`
	Util           int     `json:"util"`
	CtR            int     `json:"ctR"`
	CtRW           int     `json:"ctRW"`
	TR             int     `json:"TR"`
	TRW            int     `json:"TRW"`
	Deaths         int     `json:"deaths"`

	//kinda garbo
	Normalizer int `json:"normalizer"`
}

func NewTR() *TeamStats {
	return &TeamStats{
		TR: 1,
	}
}

func NewCtR() *TeamStats {
	return &TeamStats{
		CtR: 1,
	}
}
