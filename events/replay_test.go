package events_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/burrbd/dip/events"
	"github.com/cheekybits/is"
)

// mockEngine records SubmitOrder calls for test assertions.
type mockEngine struct {
	submitted []submittedOrder
	submitErr error
}

type submittedOrder struct {
	nation string
	order  string
}

func (m *mockEngine) SubmitOrder(nation, orderText string) error {
	if m.submitErr != nil {
		return m.submitErr
	}
	m.submitted = append(m.submitted, submittedOrder{nation, orderText})
	return nil
}

// mockLoader returns a configured mockEngine.
type mockLoader struct {
	eng     *mockEngine
	loadErr error
}

func (ml *mockLoader) Load(_ []byte) (events.EngineState, error) {
	if ml.loadErr != nil {
		return nil, ml.loadErr
	}
	return ml.eng, nil
}

// TestRebuild_NoSnapshot returns an error when the channel has no GameStarted
// or PhaseResolved event.
func TestRebuild_NoSnapshot(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	_ = events.Write(ch, "c", events.TypePlayerJoined, events.PlayerJoined{UserID: "u1", Nation: "England"})

	loader := &mockLoader{eng: &mockEngine{}}
	_, err := events.Rebuild(ch, "c", loader.Load)
	is.Err(err)
}

// TestRebuild_FromGameStarted restores engine from a GameStarted snapshot.
func TestRebuild_FromGameStarted(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	snap := json.RawMessage(`{"phase":"Spring 1901 Movement"}`)
	_ = events.Write(ch, "c", events.TypeGameStarted, events.GameStarted{InitialState: snap})

	eng := &mockEngine{}
	loader := &mockLoader{eng: eng}
	got, err := events.Rebuild(ch, "c", loader.Load)
	is.NoErr(err)
	is.NotNil(got)
	// No orders were submitted after the snapshot.
	is.Equal(len(eng.submitted), 0)
}

// TestRebuild_FromPhaseResolved restores engine from a PhaseResolved snapshot.
func TestRebuild_FromPhaseResolved(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	snap := json.RawMessage(`{"phase":"Fall 1901 Movement"}`)
	_ = events.Write(ch, "c", events.TypePhaseResolved, events.PhaseResolved{
		Phase:         "Fall 1901 Movement",
		StateSnapshot: snap,
	})

	eng := &mockEngine{}
	loader := &mockLoader{eng: eng}
	got, err := events.Rebuild(ch, "c", loader.Load)
	is.NoErr(err)
	is.NotNil(got)
}

// TestRebuild_ReplaysStagedOrders replays OrderSubmitted events posted after
// the snapshot, staging each order on the restored engine.
func TestRebuild_ReplaysStagedOrders(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}

	snap := json.RawMessage(`{}`)
	_ = events.Write(ch, "c", events.TypeGameStarted, events.GameStarted{InitialState: snap})
	_ = events.Write(ch, "c", events.TypeOrderSubmitted, events.OrderSubmitted{
		Nation: "England",
		Orders: []string{"A Lon-Nth", "F Eng-Nth"},
	})
	_ = events.Write(ch, "c", events.TypeOrderSubmitted, events.OrderSubmitted{
		Nation: "France",
		Orders: []string{"A Par-Bur"},
	})

	eng := &mockEngine{}
	loader := &mockLoader{eng: eng}
	got, err := events.Rebuild(ch, "c", loader.Load)
	is.NoErr(err)
	is.NotNil(got)

	// Three orders total (2 England + 1 France).
	is.Equal(len(eng.submitted), 3)
	is.Equal(eng.submitted[0], submittedOrder{"England", "A Lon-Nth"})
	is.Equal(eng.submitted[1], submittedOrder{"England", "F Eng-Nth"})
	is.Equal(eng.submitted[2], submittedOrder{"France", "A Par-Bur"})
}

// TestRebuild_IgnoresOrdersBeforeSnapshot verifies that OrderSubmitted events
// posted before the snapshot are not replayed.
func TestRebuild_IgnoresOrdersBeforeSnapshot(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}

	// An order before the snapshot.
	_ = events.Write(ch, "c", events.TypeOrderSubmitted, events.OrderSubmitted{
		Nation: "England",
		Orders: []string{"A Lon-Nth"},
	})
	// The snapshot resets state.
	snap := json.RawMessage(`{}`)
	_ = events.Write(ch, "c", events.TypeGameStarted, events.GameStarted{InitialState: snap})
	// An order after the snapshot.
	_ = events.Write(ch, "c", events.TypeOrderSubmitted, events.OrderSubmitted{
		Nation: "France",
		Orders: []string{"A Par-Bur"},
	})

	eng := &mockEngine{}
	loader := &mockLoader{eng: eng}
	got, err := events.Rebuild(ch, "c", loader.Load)
	is.NoErr(err)
	is.NotNil(got)

	// Only the post-snapshot order must be staged.
	is.Equal(len(eng.submitted), 1)
	is.Equal(eng.submitted[0], submittedOrder{"France", "A Par-Bur"})
}

// TestRebuild_UsesLastSnapshot verifies that when multiple snapshots exist,
// Rebuild uses the most recent one.
func TestRebuild_UsesLastSnapshot(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}

	_ = events.Write(ch, "c", events.TypeGameStarted, events.GameStarted{InitialState: json.RawMessage(`{}`)})
	_ = events.Write(ch, "c", events.TypeOrderSubmitted, events.OrderSubmitted{
		Nation: "England", Orders: []string{"A Lon-Nth"},
	})
	_ = events.Write(ch, "c", events.TypePhaseResolved, events.PhaseResolved{
		Phase:         "Spring 1901 Movement",
		StateSnapshot: json.RawMessage(`{}`),
	})
	// Only this order is after the last snapshot.
	_ = events.Write(ch, "c", events.TypeOrderSubmitted, events.OrderSubmitted{
		Nation: "France", Orders: []string{"A Par-Bur"},
	})

	eng := &mockEngine{}
	loader := &mockLoader{eng: eng}
	got, err := events.Rebuild(ch, "c", loader.Load)
	is.NoErr(err)
	is.NotNil(got)

	is.Equal(len(eng.submitted), 1)
	is.Equal(eng.submitted[0], submittedOrder{"France", "A Par-Bur"})
}

// TestRebuild_PropagatesLoaderError verifies that Rebuild returns the loader's error.
func TestRebuild_PropagatesLoaderError(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	_ = events.Write(ch, "c", events.TypeGameStarted, events.GameStarted{InitialState: json.RawMessage(`{}`)})

	loader := &mockLoader{loadErr: errors.New("load failed")}
	_, err := events.Rebuild(ch, "c", loader.Load)
	is.Err(err)
}

// TestRebuild_SkipsNonOrderEventsAfterSnapshot verifies that non-OrderSubmitted
// events posted after the snapshot (e.g. PhaseSkipped) are ignored.
func TestRebuild_SkipsNonOrderEventsAfterSnapshot(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}

	_ = events.Write(ch, "c", events.TypeGameStarted, events.GameStarted{InitialState: json.RawMessage(`{}`)})
	_ = events.Write(ch, "c", events.TypePhaseSkipped, events.PhaseSkipped{Phase: "Spring 1901 Retreat", Reason: "no_dislodgements"})
	_ = events.Write(ch, "c", events.TypeOrderSubmitted, events.OrderSubmitted{
		Nation: "England", Orders: []string{"A Lon-Nth"},
	})

	eng := &mockEngine{}
	loader := &mockLoader{eng: eng}
	got, err := events.Rebuild(ch, "c", loader.Load)
	is.NoErr(err)
	is.NotNil(got)
	// Only the OrderSubmitted event stages an order; PhaseSkipped is ignored.
	is.Equal(len(eng.submitted), 1)
}

// TestRebuild_SkipsMalformedGameStartedPayload verifies that Rebuild skips a
// GameStarted envelope whose payload cannot be decoded, then finds the next
// valid snapshot.
func TestRebuild_SkipsMalformedGameStartedPayload(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}

	// Inject a GameStarted envelope with a payload that cannot decode to GameStarted.
	bad := events.Envelope{Type: events.TypeGameStarted, Payload: json.RawMessage(`"not-an-object"`)}
	data, _ := json.Marshal(bad)
	ch.messages = append(ch.messages, string(data))

	// A valid GameStarted follows.
	_ = events.Write(ch, "c", events.TypeGameStarted, events.GameStarted{InitialState: json.RawMessage(`{}`)})

	eng := &mockEngine{}
	loader := &mockLoader{eng: eng}
	got, err := events.Rebuild(ch, "c", loader.Load)
	is.NoErr(err)
	is.NotNil(got)
}

// TestRebuild_SkipsMalformedPhaseResolvedPayload verifies that Rebuild skips a
// PhaseResolved envelope whose payload cannot be decoded.
func TestRebuild_SkipsMalformedPhaseResolvedPayload(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}

	// Inject a PhaseResolved envelope with an invalid payload.
	bad := events.Envelope{Type: events.TypePhaseResolved, Payload: json.RawMessage(`"not-an-object"`)}
	data, _ := json.Marshal(bad)
	ch.messages = append(ch.messages, string(data))

	// A valid snapshot follows so Rebuild can succeed.
	_ = events.Write(ch, "c", events.TypeGameStarted, events.GameStarted{InitialState: json.RawMessage(`{}`)})

	eng := &mockEngine{}
	loader := &mockLoader{eng: eng}
	got, err := events.Rebuild(ch, "c", loader.Load)
	is.NoErr(err)
	is.NotNil(got)
}

// TestRebuild_SkipsMalformedOrderSubmittedPayload verifies that Rebuild skips
// an OrderSubmitted envelope whose payload cannot be decoded and continues.
func TestRebuild_SkipsMalformedOrderSubmittedPayload(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}

	_ = events.Write(ch, "c", events.TypeGameStarted, events.GameStarted{InitialState: json.RawMessage(`{}`)})

	// Inject an OrderSubmitted envelope with an invalid payload after the snapshot.
	bad := events.Envelope{Type: events.TypeOrderSubmitted, Payload: json.RawMessage(`"not-an-object"`)}
	data, _ := json.Marshal(bad)
	ch.messages = append(ch.messages, string(data))

	eng := &mockEngine{}
	loader := &mockLoader{eng: eng}
	got, err := events.Rebuild(ch, "c", loader.Load)
	is.NoErr(err)
	is.NotNil(got)
	is.Equal(len(eng.submitted), 0)
}

// TestRebuild_PropagatesChannelError verifies that Rebuild returns a channel scan error.
func TestRebuild_PropagatesChannelError(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{historyErr: errors.New("channel down")}

	loader := &mockLoader{eng: &mockEngine{}}
	_, err := events.Rebuild(ch, "c", loader.Load)
	is.Err(err)
}

// TestRebuild_PropagatesSubmitError verifies that Rebuild returns an error when
// staging a replayed order fails.
func TestRebuild_PropagatesSubmitError(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}

	_ = events.Write(ch, "c", events.TypeGameStarted, events.GameStarted{InitialState: json.RawMessage(`{}`)})
	_ = events.Write(ch, "c", events.TypeOrderSubmitted, events.OrderSubmitted{
		Nation: "England", Orders: []string{"A Lon-Nth"},
	})

	eng := &mockEngine{submitErr: errors.New("parse error")}
	loader := &mockLoader{eng: eng}
	_, err := events.Rebuild(ch, "c", loader.Load)
	is.Err(err)
}
