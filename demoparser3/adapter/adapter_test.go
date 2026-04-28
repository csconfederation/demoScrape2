package adapter

import (
	"testing"

	"github.com/csconfederation/demoparser3/domain"
	"github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs"
	demoinfocscommon "github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs/common"
	dp "github.com/markus-wa/godispatch"
	"github.com/stretchr/testify/assert"
)

type fakeParser struct{}

func (fakeParser) ParseToEnd() error                                    { return nil }
func (fakeParser) RegisterNetMessageHandler(_ any) dp.HandlerIdentifier { return nil }
func (fakeParser) RegisterEventHandler(_ any) dp.HandlerIdentifier      { return nil }
func (fakeParser) GameState() demoinfocs.GameState                      { return nil }

func TestNewAdapterWithParser(t *testing.T) {
	bus := domain.NewEventBus()

	adapter := NewAdapterWithParser(fakeParser{}, bus)

	assert.NotNil(t, adapter)
	assert.Equal(t, adapter.parser, fakeParser{})
	assert.Equal(t, adapter.bus, bus)
}

func TestGetSteamID(t *testing.T) {
	player := &demoinfocscommon.Player{
		SteamID64: 123,
	}

	got := getSteamID(player)

	assert.Equal(t, uint64(123), got)
}

func TestIsUtility(t *testing.T) {
	he := &demoinfocscommon.Equipment{Type: demoinfocscommon.EqHE}
	molotov := &demoinfocscommon.Equipment{Type: demoinfocscommon.EqMolotov}
	incendiary := &demoinfocscommon.Equipment{Type: demoinfocscommon.EqIncendiary}
	smoke := &demoinfocscommon.Equipment{Type: demoinfocscommon.EqSmoke}
	decoy := &demoinfocscommon.Equipment{Type: demoinfocscommon.EqDecoy}
	ak := &demoinfocscommon.Equipment{Type: demoinfocscommon.EqAK47}
	bomb := &demoinfocscommon.Equipment{Type: demoinfocscommon.EqBomb}

	gotHE := isUtility(he)
	gotMolotov := isUtility(molotov)
	gotIncendiary := isUtility(incendiary)
	gotSmoke := isUtility(smoke)
	gotDecoy := isUtility(decoy)
	gotNil := isUtility(nil)
	gotAK := isUtility(ak)
	gotBomb := isUtility(bomb)

	assert.True(t, gotHE)
	assert.True(t, gotMolotov)
	assert.True(t, gotIncendiary)
	assert.True(t, gotSmoke)
	assert.True(t, gotDecoy)
	assert.False(t, gotNil)
	assert.False(t, gotAK)
	assert.False(t, gotBomb)
}
