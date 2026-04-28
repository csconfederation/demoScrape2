package adapter

import (
	"testing"

	"github.com/csconfederation/demoparser3/domain"
	"github.com/csconfederation/demoparser3/domain/events"
	"github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs"
	demoinfocscommon "github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs/common"
	demoinfocsmsg "github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs/msg"
	dp "github.com/markus-wa/godispatch"
	"github.com/stretchr/testify/assert"
)

type fakeParser struct {
	netHandler func(*demoinfocsmsg.CSVCMsg_ServerInfo)
}

func (f *fakeParser) ParseToEnd() error { return nil }
func (f *fakeParser) RegisterNetMessageHandler(handler any) dp.HandlerIdentifier {
	f.netHandler = handler.(func(*demoinfocsmsg.CSVCMsg_ServerInfo))
	return nil
}
func (f *fakeParser) RegisterEventHandler(_ any) dp.HandlerIdentifier { return nil }
func (f *fakeParser) GameState() demoinfocs.GameState                 { return nil }

type fakeHandler struct {
	called bool
	got    events.Event
	err    error
}

func (f *fakeHandler) Handle(event events.Event) error {
	f.called = true
	f.got = event
	return f.err
}

func TestNewAdapterWithParser(t *testing.T) {
	bus := domain.NewEventBus()
	parser := fakeParser{}
	adapter := NewAdapterWithParser(&parser, bus)

	assert.NotNil(t, adapter)
	assert.Same(t, &parser, adapter.parser)
	assert.Same(t, bus, adapter.bus)
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

func TestRegisterHandlers_GameStart(t *testing.T) {
	bus := domain.NewEventBus()
	parser := fakeParser{}
	adapter := NewAdapterWithParser(&parser, bus)
	handler := &fakeHandler{}

	bus.Subscribe("GameStart", handler)
	adapter.registerGameStart()

	msg := demoinfocsmsg.CSVCMsg_ServerInfo{MapName: new("test")}

	parser.netHandler(&msg)

	assert.True(t, handler.called)
	got := handler.got.(events.GameStart)
	assert.Equal(t, "test", got.MapName)
}
