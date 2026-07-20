package demoscrape2

import (
	"encoding/json"
	"testing"
)

const exactSteam64 uint64 = 76561198000000123

func TestSteam64FieldsMarshalAsExactDecimalStrings(t *testing.T) {
	game := Game{
		Flags: flag{
			TClutchSteam:  exactSteam64,
			CtClutchSteam: exactSteam64,
		},
		Rounds: []*round{{
			Planter: exactSteam64,
			Defuser: exactSteam64,
			PlayerStats: map[uint64]*playerStats{
				exactSteam64: {
					SteamID:           "76561198000000123",
					MostRecentFlasher: exactSteam64,
				},
			},
		}},
		PlayerOrder: []uint64{exactSteam64},
	}

	payload, err := json.Marshal(game)
	if err != nil {
		t.Fatalf("json.Marshal(Game) error = %v", err)
	}

	var decoded map[string]any
	if err := json.Unmarshal(payload, &decoded); err != nil {
		t.Fatalf("json.Unmarshal(payload) error = %v", err)
	}
	if _, ok := decoded["PlayerOrder"]; ok {
		t.Fatal("payload contains unexpected PlayerOrder key")
	}

	assertString := func(path string, value any) {
		t.Helper()
		if value != "76561198000000123" {
			t.Errorf("%s = %#v (%T), want exact decimal string", path, value, value)
		}
	}

	flags := decoded["flags"].(map[string]any)
	assertString("flags.TClutchSteam", flags["TClutchSteam"])
	assertString("flags.ctClutchSteam", flags["ctClutchSteam"])

	round := decoded["rounds"].([]any)[0].(map[string]any)
	assertString("round.planter", round["planter"])
	assertString("round.defuser", round["defuser"])
	player := round["playerStats"].(map[string]any)["76561198000000123"].(map[string]any)
	assertString("player.steamID", player["steamID"])
	assertString("player.mostRecentFlasher", player["mostRecentFlasher"])
	assertString("playerOrder[0]", decoded["playerOrder"].([]any)[0])
}

func TestNilPlayerOrderMarshalsAsNull(t *testing.T) {
	payload, err := json.Marshal(Game{})
	if err != nil {
		t.Fatalf("json.Marshal(Game{}) error = %v", err)
	}
	var decoded map[string]any
	if err := json.Unmarshal(payload, &decoded); err != nil {
		t.Fatalf("json.Unmarshal(payload) error = %v", err)
	}
	if decoded["playerOrder"] != nil {
		t.Fatalf("nil playerOrder = %#v, want null", decoded["playerOrder"])
	}
}

func TestSteam64ZeroSentinelsMarshalAsStrings(t *testing.T) {
	payload, err := json.Marshal(round{})
	if err != nil {
		t.Fatalf("json.Marshal(round{}) error = %v", err)
	}
	var decoded map[string]any
	if err := json.Unmarshal(payload, &decoded); err != nil {
		t.Fatalf("json.Unmarshal(payload) error = %v", err)
	}
	if decoded["planter"] != "0" || decoded["defuser"] != "0" {
		t.Fatalf("zero sentinels = planter:%#v defuser:%#v, want quoted zeroes", decoded["planter"], decoded["defuser"])
	}
}
