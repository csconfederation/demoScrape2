package demoscrape2

import (
	"errors"
	"strings"
	"testing"
)

// TestProcessDemoErrorContract pins the error semantics ProcessDemo relies on
// when it does `errors.Join(err, procErr)` around ParseToEnd's error and
// ErrNoValidRounds (see parser.go). ProcessDemo itself needs a real demo file
// to exercise end-to-end, which this package has no fixture for, so this pins
// the join behavior directly using errors.Join the same way ProcessDemo does.
//
// Two things downstream depends on and that would silently break if
// errors.Join's semantics ever changed out from under this code:
//
//  1. csgo-demo-worker classifies the joined error with errors.Is(err,
//     ErrNoValidRounds) to map it to a 422.
//  2. csgo-demo-worker (and other classification, e.g. ErrUnexpectedEndOfDemo
//     / ErrInvalidFileType from demoinfocs) also does substring matching via
//     strings.Contains(err.Error(), "...") against the *upstream* parse
//     error's message, so the joined error's text must still contain it.
func TestProcessDemoErrorContract(t *testing.T) {
	parseErr := errors.New("some upstream ParseToEnd failure: unexpected end of demo")

	joined := errors.Join(parseErr, ErrNoValidRounds)

	if joined == nil {
		t.Fatal("errors.Join(parseErr, ErrNoValidRounds) = nil, want non-nil")
	}
	if !errors.Is(joined, ErrNoValidRounds) {
		t.Fatalf("errors.Is(joined, ErrNoValidRounds) = false, want true (joined = %v)", joined)
	}
	if !strings.Contains(joined.Error(), parseErr.Error()) {
		t.Fatalf("joined.Error() = %q, want it to contain upstream message %q", joined.Error(), parseErr.Error())
	}
	if !strings.Contains(joined.Error(), ErrNoValidRounds.Error()) {
		t.Fatalf("joined.Error() = %q, want it to contain ErrNoValidRounds message %q", joined.Error(), ErrNoValidRounds.Error())
	}
}

// TestProcessDemoErrorContractNilParseErr pins the same contract for the case
// where ParseToEnd itself succeeds (err == nil) but endOfMatchProcessing still
// finds no valid rounds. ProcessDemo joins a nil with ErrNoValidRounds in that
// case (parser.go: errors.Join(err, procErr) where err is nil) — this must
// still produce a non-nil, errors.Is-matchable error, not silently vanish the
// way a naive "if err != nil" wrapper might.
func TestProcessDemoErrorContractNilParseErr(t *testing.T) {
	joined := errors.Join(nil, ErrNoValidRounds)

	if joined == nil {
		t.Fatal("errors.Join(nil, ErrNoValidRounds) = nil, want non-nil")
	}
	if !errors.Is(joined, ErrNoValidRounds) {
		t.Fatalf("errors.Is(joined, ErrNoValidRounds) = false, want true (joined = %v)", joined)
	}
}
