package demoscrape2

import (
	"errors"
	"testing"
)

// A demo that recorded no rounds at all — the shape produced by an aborted
// force-start, which used to panic here with "index out of range [-1]".
func TestRemoveInvalidRoundsNoRounds(t *testing.T) {
	game := InitGameObject()
	game.Rounds = nil

	err := removeInvalidRounds(game)

	if !errors.Is(err, ErrNoValidRounds) {
		t.Fatalf("removeInvalidRounds() = %v, want ErrNoValidRounds", err)
	}
}

// Every recorded round is filtered out (knife round only, no integrity). The
// result is just as unimportable as having no rounds, and must not pass as an
// empty-but-successful parse.
func TestRemoveInvalidRoundsAllFiltered(t *testing.T) {
	game := InitGameObject()
	game.Rounds = []*round{
		{RoundNum: 1, KnifeRound: true, IntegrityCheck: true},
	}

	err := removeInvalidRounds(game)

	if !errors.Is(err, ErrNoValidRounds) {
		t.Fatalf("removeInvalidRounds() = %v, want ErrNoValidRounds", err)
	}
	if len(game.Rounds) != 0 {
		t.Fatalf("len(game.Rounds) = %d, want 0", len(game.Rounds))
	}
}

func TestRemoveInvalidRoundsKeepsValidRounds(t *testing.T) {
	game := InitGameObject()
	game.Rounds = []*round{
		{RoundNum: 1, IntegrityCheck: true},
		{RoundNum: 2, IntegrityCheck: true},
	}

	if err := removeInvalidRounds(game); err != nil {
		t.Fatalf("removeInvalidRounds() = %v, want nil", err)
	}
	if len(game.Rounds) != 2 {
		t.Fatalf("len(game.Rounds) = %d, want 2", len(game.Rounds))
	}
	// rounds must come back in ascending order, not the reversed collection order
	if game.Rounds[0].RoundNum != 1 || game.Rounds[1].RoundNum != 2 {
		t.Fatalf("rounds out of order: got %d, %d", game.Rounds[0].RoundNum, game.Rounds[1].RoundNum)
	}
}

// endOfMatchProcessing must surface the error rather than returning a
// well-formed, totally empty Game — that would silently import a 0-round match.
func TestEndOfMatchProcessingNoRounds(t *testing.T) {
	game := InitGameObject()
	game.Rounds = nil

	err := endOfMatchProcessing(game)

	if !errors.Is(err, ErrNoValidRounds) {
		t.Fatalf("endOfMatchProcessing() = %v, want ErrNoValidRounds", err)
	}
}
