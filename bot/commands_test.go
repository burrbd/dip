package bot

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/burrbd/dip/engine"
	"github.com/burrbd/dip/events"
	"github.com/cheekybits/is"
)

// ---- mock channel -----------------------------------------------------------

type mockChannel struct {
	msgs    []string
	postErr error
	histErr error
}

func (m *mockChannel) Post(_, text string) error {
	if m.postErr != nil {
		return m.postErr
	}
	m.msgs = append(m.msgs, text)
	return nil
}

func (m *mockChannel) History(_ string) ([]string, error) {
	if m.histErr != nil {
		return nil, m.histErr
	}
	return m.msgs, nil
}

func (m *mockChannel) lastEventType() events.EventType {
	if len(m.msgs) == 0 {
		return ""
	}
	var env events.Envelope
	if err := json.Unmarshal([]byte(m.msgs[len(m.msgs)-1]), &env); err != nil {
		return ""
	}
	return env.Type
}

// ---- mock notifier ----------------------------------------------------------

type mockNotifier struct{}

func (n *mockNotifier) Notify(_, _ string) error { return nil }

// ---- mock engine ------------------------------------------------------------

type mockEngine struct {
	phase    string
	dump     []byte
	dumpErr  error
	startErr error
}

func (e *mockEngine) SubmitOrder(_, _ string) error             { return nil }
func (e *mockEngine) Resolve() (engine.ResolutionResult, error) { return engine.ResolutionResult{}, nil }
func (e *mockEngine) Advance() error                            { return nil }
func (e *mockEngine) SoloWinner() string                        { return "" }
func (e *mockEngine) Dump() ([]byte, error)                     { return e.dump, e.dumpErr }
func (e *mockEngine) Phase() string                             { return e.phase }

// ---- helpers ----------------------------------------------------------------

func goodEngine() *mockEngine {
	return &mockEngine{phase: "Spring 1901 Movement", dump: []byte(`{}`)}
}

func goodFactory() EngineFactory {
	return func(_ string) (engine.Engine, error) { return goodEngine(), nil }
}

func newTestDispatcher(ch *mockChannel) *Dispatcher {
	return New(ch, &mockNotifier{}, nil, goodFactory())
}

// joinPlayers posts PlayerJoined events for n players with distinct nations.
func joinPlayers(ch *mockChannel, channelID string, n int) {
	nations := []string{"England", "France", "Germany", "Italy", "Austria", "Russia", "Turkey"}
	for i := 0; i < n; i++ {
		_ = events.Write(ch, channelID, events.TypePlayerJoined, events.PlayerJoined{
			UserID: "u" + string(rune('1'+i)),
			Nation: nations[i],
		})
	}
}

// seedGameCreated posts a GameCreated event to ch.
func seedGameCreated(ch *mockChannel, gmID string) {
	_ = events.Write(ch, "chan1", events.TypeGameCreated, events.GameCreated{
		Variant:       "classical",
		DeadlineHours: 24,
		GMUserID:      gmID,
	})
}

// ---- Dispatch unknown command ------------------------------------------------

func TestDispatch_UnknownCommand(t *testing.T) {
	is := is.New(t)
	d := newTestDispatcher(&mockChannel{})
	_, err := d.Dispatch(Command{Name: "unknown", ChannelID: "chan1", UserID: "u1"})
	is.Err(err)
}

// ---- /newgame ----------------------------------------------------------------

func TestDispatchNewGame_PostsGameCreatedEvent(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)

	_, err := d.Dispatch(Command{Name: "newgame", ChannelID: "chan1", UserID: "gm1"})
	is.NoErr(err)
	is.Equal(ch.lastEventType(), events.TypeGameCreated)
}

func TestDispatchNewGame_SetsGMToCallerUserID(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)

	_, err := d.Dispatch(Command{Name: "newgame", ChannelID: "chan1", UserID: "gm42"})
	is.NoErr(err)

	var env events.Envelope
	is.NoErr(json.Unmarshal([]byte(ch.msgs[0]), &env))
	var gc events.GameCreated
	is.NoErr(json.Unmarshal(env.Payload, &gc))
	is.Equal(gc.GMUserID, "gm42")
}

func TestDispatchNewGame_ReturnsNonEmptyText(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)

	resp, err := d.Dispatch(Command{Name: "newgame", ChannelID: "chan1", UserID: "gm1"})
	is.NoErr(err)
	is.NotNil(resp)
	if resp == "" {
		t.Error("expected non-empty response text")
	}
}

func TestDispatchNewGame_RejectsIfGameAlreadyExists(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	seedGameCreated(ch, "gm1")
	d := newTestDispatcher(ch)

	_, err := d.Dispatch(Command{Name: "newgame", ChannelID: "chan1", UserID: "gm2"})
	is.Err(err)
}

func TestDispatchNewGame_RejectsWhenChannelFails(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{histErr: errors.New("channel down")}
	d := newTestDispatcher(ch)

	_, err := d.Dispatch(Command{Name: "newgame", ChannelID: "chan1", UserID: "gm1"})
	is.Err(err)
}

func TestDispatchNewGame_RejectsWhenPostFails(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{postErr: errors.New("post failed")}
	d := newTestDispatcher(ch)

	_, err := d.Dispatch(Command{Name: "newgame", ChannelID: "chan1", UserID: "gm1"})
	is.Err(err)
}

// ---- /join ------------------------------------------------------------------

func TestDispatchJoin_PostsPlayerJoinedEvent(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	seedGameCreated(ch, "gm1")
	d := newTestDispatcher(ch)

	_, err := d.Dispatch(Command{Name: "join", Args: []string{"England"}, ChannelID: "chan1", UserID: "u1"})
	is.NoErr(err)
	is.Equal(ch.lastEventType(), events.TypePlayerJoined)
}

func TestDispatchJoin_RecordsUserAndNation(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	seedGameCreated(ch, "gm1")
	d := newTestDispatcher(ch)

	_, err := d.Dispatch(Command{Name: "join", Args: []string{"France"}, ChannelID: "chan1", UserID: "u99"})
	is.NoErr(err)

	var env events.Envelope
	is.NoErr(json.Unmarshal([]byte(ch.msgs[len(ch.msgs)-1]), &env))
	var pj events.PlayerJoined
	is.NoErr(json.Unmarshal(env.Payload, &pj))
	is.Equal(pj.UserID, "u99")
	is.Equal(pj.Nation, "France")
}

func TestDispatchJoin_RejectsIfNoGame(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)

	_, err := d.Dispatch(Command{Name: "join", Args: []string{"England"}, ChannelID: "chan1", UserID: "u1"})
	is.Err(err)
}

func TestDispatchJoin_RejectsIfGameAlreadyStarted(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	seedGameCreated(ch, "gm1")
	_ = events.Write(ch, "chan1", events.TypeGameStarted, events.GameStarted{
		InitialState: json.RawMessage(`{}`),
	})
	d := newTestDispatcher(ch)

	_, err := d.Dispatch(Command{Name: "join", Args: []string{"England"}, ChannelID: "chan1", UserID: "u1"})
	is.Err(err)
}

func TestDispatchJoin_RejectsNationAlreadyTaken(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	seedGameCreated(ch, "gm1")
	_ = events.Write(ch, "chan1", events.TypePlayerJoined, events.PlayerJoined{UserID: "u1", Nation: "England"})
	d := newTestDispatcher(ch)

	_, err := d.Dispatch(Command{Name: "join", Args: []string{"England"}, ChannelID: "chan1", UserID: "u2"})
	is.Err(err)
}

func TestDispatchJoin_RejectsUnknownNation(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	seedGameCreated(ch, "gm1")
	d := newTestDispatcher(ch)

	_, err := d.Dispatch(Command{Name: "join", Args: []string{"Atlantis"}, ChannelID: "chan1", UserID: "u1"})
	is.Err(err)
}

func TestDispatchJoin_RejectsPlayerAlreadyJoined(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	seedGameCreated(ch, "gm1")
	_ = events.Write(ch, "chan1", events.TypePlayerJoined, events.PlayerJoined{UserID: "u1", Nation: "England"})
	d := newTestDispatcher(ch)

	_, err := d.Dispatch(Command{Name: "join", Args: []string{"France"}, ChannelID: "chan1", UserID: "u1"})
	is.Err(err)
}

func TestDispatchJoin_RejectsMissingNationArg(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	seedGameCreated(ch, "gm1")
	d := newTestDispatcher(ch)

	_, err := d.Dispatch(Command{Name: "join", Args: []string{}, ChannelID: "chan1", UserID: "u1"})
	is.Err(err)
}

func TestDispatchJoin_RejectsWhenChannelFails(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{histErr: errors.New("channel down")}
	d := newTestDispatcher(ch)

	_, err := d.Dispatch(Command{Name: "join", Args: []string{"England"}, ChannelID: "chan1", UserID: "u1"})
	is.Err(err)
}

// ---- /start -----------------------------------------------------------------

func TestDispatchStart_PostsGameStartedEvent(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	seedGameCreated(ch, "gm1")
	joinPlayers(ch, "chan1", 2)
	d := newTestDispatcher(ch)
	d.sessions["chan1"] = nil // pre-set to allow overwrite

	_, err := d.Dispatch(Command{Name: "start", ChannelID: "chan1", UserID: "gm1"})
	is.NoErr(err)
	is.Equal(ch.lastEventType(), events.TypeGameStarted)
}

func TestDispatchStart_GameStartedContainsSnapshot(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	seedGameCreated(ch, "gm1")
	joinPlayers(ch, "chan1", 2)
	d := newTestDispatcher(ch)

	_, err := d.Dispatch(Command{Name: "start", ChannelID: "chan1", UserID: "gm1"})
	is.NoErr(err)

	var env events.Envelope
	is.NoErr(json.Unmarshal([]byte(ch.msgs[len(ch.msgs)-1]), &env))
	var gs events.GameStarted
	is.NoErr(json.Unmarshal(env.Payload, &gs))
	is.NotNil(gs.InitialState)
}

func TestDispatchStart_StoresSession(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	seedGameCreated(ch, "gm1")
	joinPlayers(ch, "chan1", 2)
	d := newTestDispatcher(ch)

	_, err := d.Dispatch(Command{Name: "start", ChannelID: "chan1", UserID: "gm1"})
	is.NoErr(err)
	is.NotNil(d.sessions["chan1"])
}

func TestDispatchStart_RequiresGM(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	seedGameCreated(ch, "gm1")
	joinPlayers(ch, "chan1", 2)
	d := newTestDispatcher(ch)

	_, err := d.Dispatch(Command{Name: "start", ChannelID: "chan1", UserID: "notgm"})
	is.Err(err)
}

func TestDispatchStart_RejectsFewerThan2Players(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	seedGameCreated(ch, "gm1")
	joinPlayers(ch, "chan1", 1)
	d := newTestDispatcher(ch)

	_, err := d.Dispatch(Command{Name: "start", ChannelID: "chan1", UserID: "gm1"})
	is.Err(err)
}

func TestDispatchStart_RejectsMoreThan7Players(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	seedGameCreated(ch, "gm1")
	// Manually add 8 distinct player entries (reuse nations to get 8 users but cheat on nations).
	// We can't use joinPlayers for 8 since classical only has 7 nations; instead seed directly.
	for i := 0; i < 8; i++ {
		_ = events.Write(ch, "chan1", events.TypePlayerJoined, events.PlayerJoined{
			UserID: "u" + string(rune('A'+i)),
			Nation: "England", // deliberately duplicated; readState just counts players
		})
	}
	d := newTestDispatcher(ch)

	_, err := d.Dispatch(Command{Name: "start", ChannelID: "chan1", UserID: "gm1"})
	is.Err(err)
}

func TestDispatchStart_RejectsIfNoGame(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)

	_, err := d.Dispatch(Command{Name: "start", ChannelID: "chan1", UserID: "gm1"})
	is.Err(err)
}

func TestDispatchStart_RejectsIfAlreadyStarted(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	seedGameCreated(ch, "gm1")
	joinPlayers(ch, "chan1", 2)
	_ = events.Write(ch, "chan1", events.TypeGameStarted, events.GameStarted{
		InitialState: json.RawMessage(`{}`),
	})
	d := newTestDispatcher(ch)

	_, err := d.Dispatch(Command{Name: "start", ChannelID: "chan1", UserID: "gm1"})
	is.Err(err)
}

func TestDispatchStart_RejectsWhenEngineFactoryFails(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	seedGameCreated(ch, "gm1")
	joinPlayers(ch, "chan1", 2)
	d := New(ch, &mockNotifier{}, nil, func(_ string) (engine.Engine, error) {
		return nil, errors.New("engine create failed")
	})

	_, err := d.Dispatch(Command{Name: "start", ChannelID: "chan1", UserID: "gm1"})
	is.Err(err)
}

func TestDispatchStart_RejectsWhenDumpFails(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	seedGameCreated(ch, "gm1")
	joinPlayers(ch, "chan1", 2)
	d := New(ch, &mockNotifier{}, nil, func(_ string) (engine.Engine, error) {
		return &mockEngine{dump: nil, dumpErr: errors.New("dump failed")}, nil
	})

	_, err := d.Dispatch(Command{Name: "start", ChannelID: "chan1", UserID: "gm1"})
	is.Err(err)
}

func TestDispatchStart_RejectsWhenChannelFails(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{histErr: errors.New("channel down")}
	d := newTestDispatcher(ch)

	_, err := d.Dispatch(Command{Name: "start", ChannelID: "chan1", UserID: "gm1"})
	is.Err(err)
}

func TestDispatchStart_RejectsWhenPostFails(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	seedGameCreated(ch, "gm1")
	joinPlayers(ch, "chan1", 2)
	d := newTestDispatcher(ch)
	// Make the channel fail on the next post (GameStarted).
	ch.postErr = errors.New("post failed")

	_, err := d.Dispatch(Command{Name: "start", ChannelID: "chan1", UserID: "gm1"})
	is.Err(err)
}

// ---- readState malformed payload coverage -----------------------------------

// seedMalformed injects an envelope with a bad JSON payload for the given type.
func seedMalformed(ch *mockChannel, channelID string, evtType events.EventType) {
	env := events.Envelope{Type: evtType, Payload: json.RawMessage(`"bad"`)}
	data, _ := json.Marshal(env)
	ch.msgs = append(ch.msgs, string(data))
}

func TestReadState_SkipsMalformedGameCreated(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	seedMalformed(ch, "chan1", events.TypeGameCreated)
	// A valid GameCreated after the malformed one must still be picked up.
	seedGameCreated(ch, "gm1")
	d := newTestDispatcher(ch)

	state, err := d.readState("chan1")
	is.NoErr(err)
	is.Equal(state.created, true)
	is.Equal(state.gmID, "gm1")
}

func TestReadState_SkipsMalformedPlayerJoined(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	seedGameCreated(ch, "gm1")
	seedMalformed(ch, "chan1", events.TypePlayerJoined)
	// A valid join after the bad one must still be recorded.
	_ = events.Write(ch, "chan1", events.TypePlayerJoined, events.PlayerJoined{UserID: "u1", Nation: "England"})
	d := newTestDispatcher(ch)

	state, err := d.readState("chan1")
	is.NoErr(err)
	is.Equal(state.players["u1"], "England")
}

// ---- handleJoin post-error coverage -----------------------------------------

func TestDispatchJoin_RejectsWhenPostFails(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	seedGameCreated(ch, "gm1")
	d := newTestDispatcher(ch)
	// Make the channel fail on the next post (PlayerJoined).
	ch.postErr = errors.New("post failed")

	_, err := d.Dispatch(Command{Name: "join", Args: []string{"England"}, ChannelID: "chan1", UserID: "u1"})
	is.Err(err)
}
