package types

type Team struct {
	Name string `json:"name"`
	Score int `json:"score"`
}

func NewTeam(name string) *Team {
	return &Team{
		Name: name,
		Score: 0,
	}
}