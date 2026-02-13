package types

type ClutchAttempt struct {
	RoundNum     int
	IsInClutch   bool
	ClutchValue  int
	IsSuccessful bool
}

func NewClutchAttempt(roundNum, clutchValue int) *ClutchAttempt {
	return &ClutchAttempt{
		RoundNum:     roundNum,
		IsInClutch:   true,
		ClutchValue:  clutchValue,
		IsSuccessful: false,
	}
}
