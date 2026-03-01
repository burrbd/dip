package events

import (
	"encoding/json"
	"fmt"
)

// EngineState is the subset of engine.Engine that Rebuild requires: staging
// orders onto a restored game state. engine.Engine satisfies this interface.
type EngineState interface {
	SubmitOrder(nation, orderText string) error
}

// Loader restores an EngineState from a JSON snapshot produced by Dump.
// engine.Load satisfies this function signature.
type Loader func(snapshot []byte) (EngineState, error)

// Rebuild reconstructs the current game state from the channel's event log.
// It finds the most recent GameStarted or PhaseResolved event, calls load to
// restore the engine from its snapshot, then replays any OrderSubmitted events
// that were posted after that snapshot.
//
// Returns an error if no snapshot event is found, if load fails, or if a
// replayed order cannot be staged.
func Rebuild(ch Channel, channelID string, load Loader) (EngineState, error) {
	envs, err := Scan(ch, channelID)
	if err != nil {
		return nil, err
	}

	// Find the index of the last snapshot event (GameStarted or PhaseResolved).
	snapshotIdx := -1
	var snapshotBytes []byte
	for i, env := range envs {
		switch env.Type {
		case TypeGameStarted:
			var gs GameStarted
			if err := json.Unmarshal(env.Payload, &gs); err != nil {
				continue
			}
			snapshotIdx = i
			snapshotBytes = gs.InitialState
		case TypePhaseResolved:
			var pr PhaseResolved
			if err := json.Unmarshal(env.Payload, &pr); err != nil {
				continue
			}
			snapshotIdx = i
			snapshotBytes = pr.StateSnapshot
		}
	}

	if snapshotIdx < 0 {
		return nil, fmt.Errorf("events: no snapshot found in channel %q", channelID)
	}

	eng, err := load(snapshotBytes)
	if err != nil {
		return nil, fmt.Errorf("events: load snapshot: %w", err)
	}

	// Replay OrderSubmitted events posted after the snapshot.
	for _, env := range envs[snapshotIdx+1:] {
		if env.Type != TypeOrderSubmitted {
			continue
		}
		var os OrderSubmitted
		if err := json.Unmarshal(env.Payload, &os); err != nil {
			continue
		}
		for _, order := range os.Orders {
			if err := eng.SubmitOrder(os.Nation, order); err != nil {
				return nil, fmt.Errorf("events: replay order %q for %s: %w", order, os.Nation, err)
			}
		}
	}

	return eng, nil
}
