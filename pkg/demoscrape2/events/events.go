package events

import (
	"github.com/csconfederation/demoScrape2/pkg/demoscrape2/types"
	dem "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs"
)

func RegisterEvents(parser dem.Parser, game *types.Game) {
	RegisterGameEvents(parser, game)
	RegisterPlayerEvents(parser, game)
	RegisterRoundEvents(parser, game)
	RegisterBombEvents(parser, game)
}
