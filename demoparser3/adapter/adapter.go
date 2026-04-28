package adapter

import "C"
import (
	"errors"
	"io"
	"strings"

	"github.com/csconfederation/demoparser3/domain"
	adapterevents "github.com/csconfederation/demoparser3/domain/events"
	"github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs"
	demoinfocscommon "github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs/common"
	demoinfocsevents "github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs/events"
	demoinfocsmsg "github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs/msg"
)

type Adapter struct {
	parser demoParser
	bus    *domain.EventBus
	err    error
}

func NewAdapterWithParser(parser demoParser, bus *domain.EventBus) *Adapter {
	return &Adapter{
		parser: parser,
		bus:    bus,
	}
}

func NewAdapter(demo io.ReadCloser, bus *domain.EventBus) *Adapter {
	return NewAdapterWithParser(demoinfocs.NewParser(demo), bus)
}

func getSteamID(player *demoinfocscommon.Player) uint64 {
	if player == nil {
		return 0
	}

	return player.SteamID64
}

func isUtility(weapon *demoinfocscommon.Equipment) bool {
	return (weapon.Type == demoinfocscommon.EqHE) || (weapon.Type == demoinfocscommon.EqMolotov) ||
		(weapon.Type == demoinfocscommon.EqIncendiary) || (weapon.Type == demoinfocscommon.EqSmoke) ||
		(weapon.Type == demoinfocscommon.EqDecoy)
}

func isGameLive(gameState demoinfocs.GameState) bool {
	return gameState.IsMatchStarted()
}

func filterAlive(players []*demoinfocscommon.Player) []*demoinfocscommon.Player {
	alive := make([]*demoinfocscommon.Player, 0)
	for _, member := range players {
		if member.IsAlive() {
			alive = append(alive, member)
		}
	}
	return alive
}

func calculatePairwiseDistances(players []*demoinfocscommon.Player) map[uint64]map[uint64]float64 {
	distances := make(map[uint64]map[uint64]float64)
	for i := 0; i < len(players); i++ {
		for j := i + 1; j < len(players); j++ {
			p1 := players[i]
			p2 := players[j]

			dist := p1.Position().Distance(p2.Position())

			if distances[p1.SteamID64] == nil {
				distances[p1.SteamID64] = make(map[uint64]float64)
			}
			distances[p1.SteamID64][p2.SteamID64] = dist

			if distances[p2.SteamID64] == nil {
				distances[p2.SteamID64] = make(map[uint64]float64)
			}
			distances[p2.SteamID64][p1.SteamID64] = dist
		}
	}

	return distances
}

func (a *Adapter) RegisterHandlers() error {
	a.parser.RegisterNetMessageHandler(func(msg *demoinfocsmsg.CSVCMsg_ServerInfo) {
		err := a.bus.Publish(adapterevents.GameStart{
			MapName: msg.GetMapName(),
		})
		if err != nil {
			a.err = err
		}
	})

	if a.err != nil {
		return a.err
	}

	a.parser.RegisterEventHandler(func(e demoinfocsevents.MatchStart) {
		CTSide := a.parser.GameState().TeamCounterTerrorists()
		TSide := a.parser.GameState().TeamTerrorists()
		CTMembers := map[uint64]string{}
		TMembers := map[uint64]string{}

		for _, player := range CTSide.Members() {
			CTMembers[player.SteamID64] = player.Name
		}

		for _, player := range TSide.Members() {
			TMembers[player.SteamID64] = player.Name
		}

		err := a.bus.Publish(adapterevents.MatchStart{
			CTSide:    CTSide.ClanName(),
			TSide:     TSide.ClanName(),
			CTMembers: CTMembers,
			TMembers:  TMembers,
		})
		if err != nil {
			a.err = err
		}
	})

	if a.err != nil {
		return a.err
	}

	a.parser.RegisterEventHandler(func(e demoinfocsevents.RoundFreezetimeEnd) {
		if !isGameLive(a.parser.GameState()) {
			return
		}
		// TODO: Verify this is accurate
		CTMembers := make([]uint64, 0)
		TMembers := make([]uint64, 0)
		for _, member := range a.parser.GameState().TeamCounterTerrorists().Members() {
			CTMembers = append(CTMembers, member.SteamID64)
		}
		for _, member := range a.parser.GameState().TeamTerrorists().Members() {
			TMembers = append(TMembers, member.SteamID64)
		}
		err := a.bus.Publish(adapterevents.RoundStart{
			TotalRoundsPlayed:  a.parser.GameState().TotalRoundsPlayed(),
			ConnectedCTPlayers: CTMembers,
			ConnectedTPlayers:  TMembers,
		})
		if err != nil {
			a.err = err
		}
	})

	a.parser.RegisterEventHandler(func(e demoinfocsevents.RoundEnd) {
		err := a.bus.Publish(adapterevents.RoundEnd{
			WinningTeamName: e.WinnerState.ClanName(),
		})
		if err != nil {
			a.err = err
		}
	})

	a.parser.RegisterEventHandler(func(e demoinfocsevents.RoundEndOfficial) {
		err := a.bus.Publish(adapterevents.RoundEndOfficial{})
		if err != nil {
			a.err = err
		}
	})

	if a.err != nil {
		return a.err
	}

	a.parser.RegisterEventHandler(func(e demoinfocsevents.GameHalfEnded) {
		err := a.bus.Publish(adapterevents.GameHalfEnded{})
		if err != nil {
			a.err = err
		}
	})

	if a.err != nil {
		return a.err
	}

	a.parser.RegisterEventHandler(func(e demoinfocsevents.FrameDone) {
		if !isGameLive(a.parser.GameState()) {
			return
		}
		tMembers := a.parser.GameState().TeamTerrorists().Members()
		aliveTMembers := filterAlive(tMembers)
		distances := calculatePairwiseDistances(aliveTMembers)
		err := a.bus.Publish(adapterevents.FrameDone{
			Tick:      a.parser.GameState().IngameTick(),
			Distances: distances,
		})
		if err != nil {
			a.err = err
		}
	})

	if a.err != nil {
		return a.err
	}

	// Called before roundend
	a.parser.RegisterEventHandler(func(e demoinfocsevents.ScoreUpdated) {
		if !isGameLive(a.parser.GameState()) {
			return
		}
		err := a.bus.Publish(adapterevents.ScoreUpdated{
			TeamName: e.TeamState.ClanName(),
			OldScore: e.OldScore,
			NewScore: e.NewScore,
		})
		if err != nil {
			a.err = err
		}
	})

	if a.err != nil {
		return a.err
	}

	a.parser.RegisterEventHandler(func(e demoinfocsevents.GrenadeProjectileThrow) {
		if !isGameLive(a.parser.GameState()) {
			return
		}
		utility := e.Projectile.WeaponInstance.Type.String()
		if utility == demoinfocscommon.EqMolotov.String() || utility == demoinfocscommon.EqIncendiary.String() {
			utility = "fire"
		}
		err := a.bus.Publish(adapterevents.UtilityThrown{
			ThrowerID:   getSteamID(e.Projectile.Thrower),
			UtilityType: strings.ToLower(utility),
		})
		if err != nil {
			a.err = err
		}
	})

	if a.err != nil {
		return a.err
	}

	a.parser.RegisterEventHandler(func(e demoinfocsevents.Kill) {
		if !isGameLive(a.parser.GameState()) {
			return
		}

		if e.Victim != nil {
			err := a.bus.Publish(adapterevents.Death{
				VictimID:       getSteamID(e.Victim),
				VictimTeamName: e.Victim.TeamState.ClanName(),
			})
			if err != nil {
				a.err = err
				return
			}

			if e.Assister != nil && e.Assister.Team != e.Victim.Team {
				err := a.bus.Publish(adapterevents.Assist{
					VictimID:        getSteamID(e.Victim),
					AssisterID:      getSteamID(e.Assister),
					IsAssistedFlash: e.AssistedFlash,
					Tick:            a.parser.GameState().IngameTick(),
				})
				if err != nil {
					a.err = err
					return
				}
			}

			if e.Killer != nil {
				weapon := e.Weapon.Type
				if weapon == demoinfocscommon.EqWorld {
					return
				}
				isAWPKill := false
				if weapon == demoinfocscommon.EqAWP {
					isAWPKill = true
				}
				err = a.bus.Publish(adapterevents.Kill{
					KillerID:       getSteamID(e.Killer),
					VictimID:       getSteamID(e.Victim),
					IsAWPKill:      isAWPKill,
					IsHeadshot:     e.IsHeadshot,
					Tick:           a.parser.GameState().IngameTick(),
					KillerTeamName: e.Killer.TeamState.ClanName(),
				})
				if err != nil {
					a.err = err
					return
				}
			}
		}
	})

	if a.err != nil {
		return a.err
	}

	a.parser.RegisterEventHandler(func(e demoinfocsevents.BombExplode) {
		err := a.bus.Publish(adapterevents.BombExplode{
			PlanterID: getSteamID(e.Player),
		})
		if err != nil {
			a.err = err
		}
	})

	if a.err != nil {
		return a.err
	}

	a.parser.RegisterEventHandler(func(e demoinfocsevents.BombDefused) {
		err := a.bus.Publish(adapterevents.BombDefused{
			DefuserID: getSteamID(e.Player),
		})
		if err != nil {
			a.err = err
		}
	})

	if a.err != nil {
		return a.err
	}

	a.parser.RegisterEventHandler(func(e demoinfocsevents.PlayerHurt) {
		if e.Weapon.Type == demoinfocscommon.EqBomb {
			return
		}
		isFireDamage := false
		isHEDamage := false
		if e.Weapon.Type == demoinfocscommon.EqMolotov || e.Weapon.Type == demoinfocscommon.EqIncendiary {
			isFireDamage = true
		}

		if e.Weapon.Type == demoinfocscommon.EqHE {
			isHEDamage = true
		}

		err := a.bus.Publish(adapterevents.PlayerHurt{
			AttackerID:   getSteamID(e.Player),
			VictimID:     getSteamID(e.Player),
			Weapon:       e.WeaponString,
			HealthDamage: e.HealthDamage,
			IsUtility:    isUtility(e.Weapon),
			IsFireDamage: isFireDamage,
			IsHEDamage:   isHEDamage,
		})
		if err != nil {
			a.err = err
		}
	})

	if a.err != nil {
		return a.err
	}

	a.parser.RegisterEventHandler(func(e demoinfocsevents.PlayerFlashed) {
		if e.Attacker.Team == e.Player.Team {
			return
		}

		if e.FlashDuration().Seconds() < 1.0 {
			return
		}

		err := a.bus.Publish(adapterevents.PlayerFlashed{
			AttackerID:    getSteamID(e.Attacker),
			VictimID:      getSteamID(e.Player),
			FlashDuration: e.FlashDuration().Seconds(),
			Tick:          a.parser.GameState().IngameTick(),
		})
		if err != nil {
			a.err = err
		}
	})

	if a.err != nil {
		return a.err
	}

	return nil
}

func (a *Adapter) Parse() error {
	err := a.RegisterHandlers()
	if err != nil {
		return err
	}
	err = a.parser.ParseToEnd()
	if err != nil {
		if errors.Is(err, demoinfocs.ErrUnexpectedEndOfDemo) {
			return nil
		}
		return err
	}
	return nil
	//stats := a.bus.Publish()
}
