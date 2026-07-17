package demoscrape2

import "errors"

// ErrNoValidRounds reports a demo that yielded nothing importable: either it
// recorded no rounds at all, or every round it did record was filtered out by
// removeInvalidRounds (knife/veto/redo rounds).
//
// This is a property of the demo, not of the parser, so callers should treat it
// as a permanent rejection of that file rather than a retryable failure.
var ErrNoValidRounds = errors.New("demoscrape2: demo contains no valid rounds (ErrNoValidRounds)")
