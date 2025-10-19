package types

type TeamStats struct {
	WinPoints      float64 `json:"winPoints"`
	ImpactPoints   float64 `json:"impactPoints"`
	TWinPoints     float64 `json:"TWinPoints"`
	CtWinPoints    float64 `json:"ctWinPoints"`
	TImpactPoints  float64 `json:"TImpactPoints"`
	CtImpactPoints float64 `json:"ctImpactPoints"`
	FourVFiveW     int     `json:"fourVFiveW" end_of_match_sum:"true"`
	FourVFiveS     int     `json:"fourVFiveS" end_of_match_sum:"true"`
	FiveVFourW     int     `json:"fiveVFourW" end_of_match_sum:"true"`
	FiveVFourS     int     `json:"fiveVFourS" end_of_match_sum:"true"`
	Pistols        int     `json:"pistols" end_of_match_sum:"true"`
	PistolsW       int     `json:"pistolsW" end_of_match_sum:"true"`
	Saves          int     `json:"saves" end_of_match_sum:"true"`
	Clutches       int     `json:"clutches" end_of_match_sum:"true"`
	Traded         int     `json:"traded"`
	Fass           int     `json:"fass"`
	Ef             int     `json:"ef"`
	Ud             int     `json:"ud"`
	Util           int     `json:"util"`
	CtR            int     `json:"ctR" end_of_match_sum:"true"`
	CtRW           int     `json:"ctRW" end_of_match_sum:"true"`
	TR             int     `json:"TR" end_of_match_sum:"true"`
	TRW            int     `json:"TRW" end_of_match_sum:"true"`
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
