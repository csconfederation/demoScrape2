package demoscrape2

import (
	"io"
	"strings"
	"testing"
)

// TestProcessDemoReturnsNonNilGameOnError pins the invariant documented on
// ProcessDemo (parser.go): it always returns a non-nil *Game alongside a
// non-nil error, even when the input can't be parsed at all.
// csgo-demo-worker/main.go reads game.Result on the err != nil path and would
// nil-deref if this ever regressed.
//
// This doesn't need a real demo fixture: garbage input fails fast inside
// demoinfocs (invalid/missing header), which is enough to exercise the
// non-nil-Game contract on ProcessDemo's error path without a large fixture
// or a slow parse.
func TestProcessDemoReturnsNonNilGameOnError(t *testing.T) {
	garbage := io.NopCloser(strings.NewReader("this is not a valid demo file"))

	game, err := ProcessDemo(garbage)

	if game == nil {
		t.Fatal("ProcessDemo() game = nil, want non-nil *Game even on error")
	}
	if err == nil {
		t.Fatal("ProcessDemo() err = nil, want non-nil error for garbage input")
	}
}
