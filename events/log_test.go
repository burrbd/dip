package events_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/burrbd/dip/events"
	"github.com/cheekybits/is"
)

// mockChannel is an in-memory Channel for tests.
type mockChannel struct {
	messages  []string
	postErr   error
	historyErr error
}

func (m *mockChannel) Post(_, text string) error {
	if m.postErr != nil {
		return m.postErr
	}
	m.messages = append(m.messages, text)
	return nil
}

func (m *mockChannel) History(_ string) ([]string, error) {
	if m.historyErr != nil {
		return nil, m.historyErr
	}
	return m.messages, nil
}

// TestWrite_PostsJSONEnvelope verifies that Write encodes the event as a JSON
// Envelope and posts it to the channel.
func TestWrite_PostsJSONEnvelope(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}

	payload := events.GameCreated{Variant: "classical", DeadlineHours: 24, GMUserID: "u1"}
	err := events.Write(ch, "chan1", events.TypeGameCreated, payload)
	is.NoErr(err)
	is.Equal(len(ch.messages), 1)

	var env events.Envelope
	is.NoErr(json.Unmarshal([]byte(ch.messages[0]), &env))
	is.Equal(env.Type, events.TypeGameCreated)

	var got events.GameCreated
	is.NoErr(json.Unmarshal(env.Payload, &got))
	is.Equal(got.Variant, "classical")
	is.Equal(got.DeadlineHours, 24)
	is.Equal(got.GMUserID, "u1")
}

// TestWrite_PropagatesPostError verifies that Write returns the channel error.
func TestWrite_PropagatesPostError(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{postErr: errors.New("network error")}

	err := events.Write(ch, "chan1", events.TypeGameCreated, events.GameCreated{})
	is.Err(err)
}

// TestScan_ParsesEvents verifies that Scan returns all valid Envelope events
// in channel-history order.
func TestScan_ParsesEvents(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}

	_ = events.Write(ch, "chan1", events.TypeGameCreated, events.GameCreated{Variant: "classical"})
	_ = events.Write(ch, "chan1", events.TypePlayerJoined, events.PlayerJoined{UserID: "u1", Nation: "England"})

	envs, err := events.Scan(ch, "chan1")
	is.NoErr(err)
	is.Equal(len(envs), 2)
	is.Equal(envs[0].Type, events.TypeGameCreated)
	is.Equal(envs[1].Type, events.TypePlayerJoined)
}

// TestScan_SkipsNonEventMessages verifies that Scan silently ignores messages
// that are not valid JSON Envelopes.
func TestScan_SkipsNonEventMessages(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{messages: []string{
		"Hello everyone, let's play Diplomacy!",
		`{"not":"an event"}`,
	}}
	_ = events.Write(ch, "chan1", events.TypeGameCreated, events.GameCreated{Variant: "classical"})

	envs, err := events.Scan(ch, "chan1")
	is.NoErr(err)
	is.Equal(len(envs), 1)
	is.Equal(envs[0].Type, events.TypeGameCreated)
}

// TestScan_PropagatesHistoryError verifies that Scan returns the channel error.
func TestScan_PropagatesHistoryError(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{historyErr: errors.New("history unavailable")}

	_, err := events.Scan(ch, "chan1")
	is.Err(err)
}

// TestScan_EmptyChannel returns no events for an empty channel.
func TestScan_EmptyChannel(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	envs, err := events.Scan(ch, "chan1")
	is.NoErr(err)
	is.Equal(len(envs), 0)
}

// TestWrite_PropagatesMarshalError verifies that Write returns an error when
// the payload cannot be marshalled to JSON (e.g. a channel value).
func TestWrite_PropagatesMarshalError(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	// chan int is not serialisable to JSON.
	err := events.Write(ch, "chan1", events.TypeGameCreated, make(chan int))
	is.Err(err)
	is.Equal(len(ch.messages), 0)
}

// TestRoundTrip_AllEventTypes verifies serialisation round-trips for every
// event type defined in the package.
func TestRoundTrip_AllEventTypes(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}

	payloads := []struct {
		typ     events.EventType
		payload any
	}{
		{events.TypeGameCreated, events.GameCreated{Variant: "classical", DeadlineHours: 48, GMUserID: "gm1"}},
		{events.TypePlayerJoined, events.PlayerJoined{UserID: "u2", Nation: "France"}},
		{events.TypeGameStarted, events.GameStarted{InitialState: json.RawMessage(`{}`)}},
		{events.TypeOrderSubmitted, events.OrderSubmitted{UserID: "u1", Nation: "England", Orders: []string{"A Lon-Nth"}, Phase: "Movement"}},
		{events.TypePhaseResolved, events.PhaseResolved{Phase: "Spring 1901 Movement", StateSnapshot: json.RawMessage(`{}`)}},
		{events.TypePhaseSkipped, events.PhaseSkipped{Phase: "Spring 1901 Retreat", Reason: "no_dislodgements"}},
		{events.TypeNMRRecorded, events.NMRRecorded{Nation: "Russia", Phase: "Movement", AutoOrders: []string{"A Mos H"}}},
		{events.TypeDrawProposed, events.DrawProposed{ProposerNation: "Turkey"}},
		{events.TypeDrawVoted, events.DrawVoted{Nation: "France", Accept: true}},
		{events.TypeGameEnded, events.GameEnded{Result: "solo", Winner: "England", FinalState: json.RawMessage(`{}`)}},
	}

	for _, p := range payloads {
		is.NoErr(events.Write(ch, "chan1", p.typ, p.payload))
	}

	envs, err := events.Scan(ch, "chan1")
	is.NoErr(err)
	is.Equal(len(envs), len(payloads))
	for i, p := range payloads {
		is.Equal(envs[i].Type, p.typ)
	}
}
