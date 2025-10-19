package demoscrape2

import (
	"io"

	"github.com/csconfederation/demoScrape2/pkg/demoscrape2/events"
	"github.com/csconfederation/demoScrape2/pkg/demoscrape2/types"
	dem "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs"
	log "github.com/sirupsen/logrus"
)

//TODO
//"Catch up on the score" - dont remember what this is lol

//BUG fix
//MAKE ROUNDENDOFFICIAL Redundant (may have done, make sure it passes validation tho)
//add verification for missed event triggers if someone DCs/Cs after the redundant event that is stale
//csgo bots all have same steamID, need to use something else just in case for bots
//MM bug?

//add support for esea games? need to change validation and how we determine what round it is (gamestats rounds doesnt work)

//FUNCTIONAL CHANGES
//add verification for if a round event has triggered so far in the round (avoid double roundEnds)
//check for Game start without pistol (if we have bad demo)
//Add backend support
//Add anchor stuff
//Add team economy round stats (ecos, forces, etc)
//Add various nil checking

//CLEAN CODE
//TODO: create a outputPlayer function to clean up output.go
//TODO: convert rating calculations to a function
//TODO: actually use killValues lmao

const DEBUG = false

//const suppressNormalOutput = false

// globals
const printChatLog = true
const printDebugLog = true
const FORCE_NEW_STATS_UPLOAD = false
const BACKEND_PUSHING = true

var killValues = map[string]float64{
	"attacking":     1.2, //base values
	"defending":     1.0,
	"bombDefense":   1.0,
	"retake":        1.2,
	"chase":         0.8,
	"exit":          0.6,
	"t_consolation": 0.5,
	"gravy":         0.6,
	"punish":        0.8,
	"entry":         0.8, //multipliers
	"t_opener":      0.3,
	"ct_opener":     0.5,
	"trade":         0.3,
	"flashAssist":   0.2,
	"assist":        0.15,
}

func ProcessDemo(demo io.ReadCloser) (*types.Game, error) {

	game := types.NewGame()
	parser := dem.NewParser(demo)
	defer func() {
		if err := parser.Close(); err != nil {
			log.Printf("Failed to close parser: %v", err)
		}
	}()

	//must parse header to get header info - method deprecated in v5
	header, err := parser.ParseHeader()
	if err != nil {
		return nil, err
	}

	game.SetTickLength(header.PlaybackTicks)

	events.RegisterEvents(parser, game)

	// Parse to end
	err = parser.ParseToEnd()

	game.EndOfMatchProcessing()

	return game, err

}
