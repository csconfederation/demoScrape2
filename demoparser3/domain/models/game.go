package models

import (
	"github.com/csconfederation/demoparser3/domain/models/stats"
	_ "github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs"
)

const MR = 12
const TickRate = 64

type Game struct {
	MapName      string
	IsGameLive   bool
	PlayersStats map[uint64]*stats.PlayerStats
	Rounds       []*Round
}

func NewGame() *Game {
	return &Game{}
}
