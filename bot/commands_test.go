package bot

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/burrbd/dip/engine"
	"github.com/burrbd/dip/events"
	"github.com/burrbd/dip/session"
	"github.com/cheekybits/is"
)

// ---- mock channel -----------------------------------------------------------

type mockChannel struct {
	msgs      []string
	dms       map[string][]string
	postErr   error
	histErr   error
	dmPostErr error
	dmHistErr error
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

func (m *mockChannel) SendDM(userID, text string) error {
	if m.dmPostErr != nil {
		return m.dmPostErr
	}
	if m.dms == nil {
		m.dms = make(map[string][]string)
	}
	m.dms[userID] = append(m.dms[userID], text)
	return nil
}

func (m *mockChannel) DMHistory(userID string) ([]string, error) {
	if m.dmHistErr != nil {
		return nil, m.dmHistErr
	}
	return m.dms[userID], nil
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
	phase      string
	dump       []byte
	dumpErr    error
	orderErr   error
	resolveErr error
	advanceErr error
	soloWinner string
}

func (e *mockEngine) SubmitOrder(_, _ string) error {
	return e.orderErr
}
func (e *mockEngine) Resolve() (engine.ResolutionResult, error) {
	return engine.ResolutionResult{Phase: e.phase}, e.resolveErr
}
func (e *mockEngine) Advance() error    { return e.advanceErr }
func (e *mockEngine) SoloWinner() string { return e.soloWinner }
func (e *mockEngine) Dump() ([]byte, error) { return e.dump, e.dumpErr }
func (e *mockEngine) Phase() string         { return e.phase }

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

// ---- DM command helpers -----------------------------------------------------

// makeDMSession creates a started session with two players (u1→England, u2→France)
// keyed by gameChID in the dispatcher's sessions map. DeadlineHours is 0 to
// prevent background timers in tests.
func makeDMSession(d *Dispatcher, ch *mockChannel, gameChID string) *session.Session {
	players := map[string]string{"u1": "England", "u2": "France"}
	eng := goodEngine()
	sess := session.New(ch, gameChID, "gm1", "Spring 1901 Movement", players, 0, eng, &mockNotifier{})
	d.sessions[gameChID] = sess
	return sess
}

// makeSinglePlayerSession creates a started session with one player (u1→England).
func makeSinglePlayerSession(d *Dispatcher, ch *mockChannel, gameChID string) *session.Session {
	players := map[string]string{"u1": "England"}
	sess := session.New(ch, gameChID, "gm1", "Spring 1901 Movement", players, 0, goodEngine(), &mockNotifier{})
	d.sessions[gameChID] = sess
	return sess
}

// dmCmd builds a DM Command targeting the given game channel.
func dmCmd(name, gameChID, userID string, args ...string) Command {
	return Command{
		Name:          name,
		Args:          args,
		UserID:        userID,
		ChannelID:     "dm-" + userID,
		IsDM:          true,
		GameChannelID: gameChID,
	}
}

// ---- /order -----------------------------------------------------------------

func TestDispatchOrder_RejectsNonDM(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	makeDMSession(d, ch, "chan1")

	_, err := d.Dispatch(Command{Name: "order", Args: []string{"A Vie-Bud"}, ChannelID: "chan1", UserID: "u1"})
	is.Err(err)
}

func TestDispatchOrder_RejectsNoActiveGame(t *testing.T) {
	is := is.New(t)
	d := newTestDispatcher(&mockChannel{})

	_, err := d.Dispatch(dmCmd("order", "chan1", "u1", "A Vie-Bud"))
	is.Err(err)
}

func TestDispatchOrder_RejectsNonMovementPhase(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	sess := makeDMSession(d, ch, "chan1")
	sess.Phase = "Spring 1901 Retreat"

	_, err := d.Dispatch(dmCmd("order", "chan1", "u1", "A Vie-Bud"))
	is.Err(err)
}

func TestDispatchOrder_RejectsNonPlayer(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	makeDMSession(d, ch, "chan1")

	_, err := d.Dispatch(dmCmd("order", "chan1", "outsider"))
	is.Err(err)
}

func TestDispatchOrder_RejectsMissingArg(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	makeDMSession(d, ch, "chan1")

	_, err := d.Dispatch(dmCmd("order", "chan1", "u1"))
	is.Err(err)
}

func TestDispatchOrder_RejectsInvalidOrder(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	sess := makeDMSession(d, ch, "chan1")
	sess.Eng = &mockEngine{
		phase:    "Spring 1901 Movement",
		dump:     []byte(`{}`),
		orderErr: errors.New("bad order"),
	}

	_, err := d.Dispatch(dmCmd("order", "chan1", "u1", "A Xyz-Foo"))
	is.Err(err)
}

func TestDispatchOrder_StagesOrder(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	sess := makeDMSession(d, ch, "chan1")

	resp, err := d.Dispatch(dmCmd("order", "chan1", "u1", "A", "Vie-Bud"))
	is.NoErr(err)
	if resp == "" {
		t.Error("expected non-empty response")
	}
	is.Equal(len(sess.StagedOrders["England"]), 1)
	is.Equal(sess.StagedOrders["England"][0], "A Vie-Bud")
}

func TestDispatchOrder_MultiWordOrderJoined(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	sess := makeDMSession(d, ch, "chan1")

	_, err := d.Dispatch(dmCmd("order", "chan1", "u1", "A", "Vie", "S", "F", "Tri-Alb"))
	is.NoErr(err)
	is.Equal(sess.StagedOrders["England"][0], "A Vie S F Tri-Alb")
}

// ---- /orders ----------------------------------------------------------------

func TestDispatchOrders_RejectsNonDM(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	makeDMSession(d, ch, "chan1")

	_, err := d.Dispatch(Command{Name: "orders", ChannelID: "chan1", UserID: "u1"})
	is.Err(err)
}

func TestDispatchOrders_RejectsNoActiveGame(t *testing.T) {
	is := is.New(t)
	d := newTestDispatcher(&mockChannel{})

	_, err := d.Dispatch(dmCmd("orders", "chan1", "u1"))
	is.Err(err)
}

func TestDispatchOrders_RejectsNonPlayer(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	makeDMSession(d, ch, "chan1")

	_, err := d.Dispatch(dmCmd("orders", "chan1", "outsider"))
	is.Err(err)
}

func TestDispatchOrders_ReturnsNoOrders(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	makeDMSession(d, ch, "chan1")

	resp, err := d.Dispatch(dmCmd("orders", "chan1", "u1"))
	is.NoErr(err)
	is.Equal(resp, "No orders staged.")
}

func TestDispatchOrders_ListsStagedOrders(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	sess := makeDMSession(d, ch, "chan1")
	sess.StagedOrders["England"] = []string{"A Vie-Bud", "F Tri-Alb"}

	resp, err := d.Dispatch(dmCmd("orders", "chan1", "u1"))
	is.NoErr(err)
	if resp == "" {
		t.Error("expected non-empty response")
	}
	// Response must mention both orders.
	if !contains(resp, "A Vie-Bud") || !contains(resp, "F Tri-Alb") {
		t.Errorf("expected both orders in response, got: %q", resp)
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsStr(s, sub))
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

// ---- /clear -----------------------------------------------------------------

func TestDispatchClear_RejectsNonDM(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	makeDMSession(d, ch, "chan1")

	_, err := d.Dispatch(Command{Name: "clear", ChannelID: "chan1", UserID: "u1"})
	is.Err(err)
}

func TestDispatchClear_RejectsNoActiveGame(t *testing.T) {
	is := is.New(t)
	d := newTestDispatcher(&mockChannel{})

	_, err := d.Dispatch(dmCmd("clear", "chan1", "u1"))
	is.Err(err)
}

func TestDispatchClear_RejectsNonPlayer(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	makeDMSession(d, ch, "chan1")

	_, err := d.Dispatch(dmCmd("clear", "chan1", "outsider"))
	is.Err(err)
}

func TestDispatchClear_ClearsAllOrders(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	sess := makeDMSession(d, ch, "chan1")
	sess.StagedOrders["England"] = []string{"A Vie-Bud", "F Tri-Alb"}
	sess.Submitted["England"] = true

	_, err := d.Dispatch(dmCmd("clear", "chan1", "u1"))
	is.NoErr(err)
	is.Equal(len(sess.StagedOrders["England"]), 0)
	is.Equal(sess.Submitted["England"], false)
}

func TestDispatchClear_ClearsSpecificOrder(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	sess := makeDMSession(d, ch, "chan1")
	sess.StagedOrders["England"] = []string{"A Vie-Bud", "F Tri-Alb"}

	_, err := d.Dispatch(dmCmd("clear", "chan1", "u1", "A Vie-Bud"))
	is.NoErr(err)
	is.Equal(len(sess.StagedOrders["England"]), 1)
	is.Equal(sess.StagedOrders["England"][0], "F Tri-Alb")
}

func TestDispatchClear_RejectsOrderNotFound(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	sess := makeDMSession(d, ch, "chan1")
	sess.StagedOrders["England"] = []string{"A Vie-Bud"}

	_, err := d.Dispatch(dmCmd("clear", "chan1", "u1", "A Lon-Nth"))
	is.Err(err)
}

// ---- /submit ----------------------------------------------------------------

func TestDispatchSubmit_RejectsNonDM(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	makeDMSession(d, ch, "chan1")

	_, err := d.Dispatch(Command{Name: "submit", ChannelID: "chan1", UserID: "u1"})
	is.Err(err)
}

func TestDispatchSubmit_RejectsNoActiveGame(t *testing.T) {
	is := is.New(t)
	d := newTestDispatcher(&mockChannel{})

	_, err := d.Dispatch(dmCmd("submit", "chan1", "u1"))
	is.Err(err)
}

func TestDispatchSubmit_RejectsNonMovementPhase(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	sess := makeDMSession(d, ch, "chan1")
	sess.Phase = "Spring 1901 Retreat"

	_, err := d.Dispatch(dmCmd("submit", "chan1", "u1"))
	is.Err(err)
}

func TestDispatchSubmit_RejectsNonPlayer(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	makeDMSession(d, ch, "chan1")

	_, err := d.Dispatch(dmCmd("submit", "chan1", "outsider"))
	is.Err(err)
}

func TestDispatchSubmit_RejectsWriteDMError(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{dmPostErr: errors.New("dm unavailable")}
	d := newTestDispatcher(ch)
	makeDMSession(d, ch, "chan1")

	_, err := d.Dispatch(dmCmd("submit", "chan1", "u1"))
	is.Err(err)
}

func TestDispatchSubmit_WritesOrderSubmittedToDM(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	sess := makeDMSession(d, ch, "chan1")
	sess.StagedOrders["England"] = []string{"A Vie-Bud"}

	_, err := d.Dispatch(dmCmd("submit", "chan1", "u1"))
	is.NoErr(err)

	// DM thread for u1 must contain an OrderSubmitted event.
	envs, scanErr := events.ScanDM(ch, "u1")
	is.NoErr(scanErr)
	is.Equal(len(envs), 1)
	is.Equal(envs[0].Type, events.TypeOrderSubmitted)

	var os events.OrderSubmitted
	is.NoErr(json.Unmarshal(envs[0].Payload, &os))
	is.Equal(os.Nation, "England")
	is.Equal(os.Phase, "Spring 1901 Movement")
	is.Equal(len(os.Orders), 1)
	is.Equal(os.Orders[0], "A Vie-Bud")
}

func TestDispatchSubmit_DoesNotFireAdvanceTurnIfNotAllSubmitted(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	makeDMSession(d, ch, "chan1") // two players: u1 and u2

	// Only u1 submits; u2 has not.
	resp, err := d.Dispatch(dmCmd("submit", "chan1", "u1"))
	is.NoErr(err)
	is.Equal(resp, "Orders submitted.")

	// No PhaseResolved event should have been posted to the game channel.
	is.Equal(len(ch.msgs), 0)
}

func TestDispatchSubmit_FiresAdvanceTurnWhenAllSubmitted(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	makeSinglePlayerSession(d, ch, "chan1") // one player: u1 only

	resp, err := d.Dispatch(dmCmd("submit", "chan1", "u1"))
	is.NoErr(err)
	if resp == "Orders submitted." {
		t.Error("expected early-resolution response, got plain submit response")
	}

	// PhaseResolved event must have been posted to the game channel.
	is.Equal(len(ch.msgs), 1)
	var env events.Envelope
	is.NoErr(json.Unmarshal([]byte(ch.msgs[0]), &env))
	is.Equal(env.Type, events.TypePhaseResolved)
}

func TestDispatchSubmit_RejectsDMHistoryError(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	makeDMSession(d, ch, "chan1")
	// WriteDM succeeds but DMHistory fails on the check call.
	ch.dmHistErr = errors.New("dm history unavailable")

	_, err := d.Dispatch(dmCmd("submit", "chan1", "u1"))
	is.Err(err)
}

func TestDispatchSubmit_RejectsAdvanceTurnError(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	sess := makeSinglePlayerSession(d, ch, "chan1")
	// Make resolve fail so AdvanceTurn returns an error.
	sess.Eng = &mockEngine{
		phase:      "Spring 1901 Movement",
		dump:       []byte(`{}`),
		resolveErr: errors.New("resolve failed"),
	}

	_, err := d.Dispatch(dmCmd("submit", "chan1", "u1"))
	is.Err(err)
}

// ---- allNationsSubmitted edge cases -----------------------------------------

// TestAllNationsSubmitted_SkipsNonOrderSubmittedEvents verifies that
// non-OrderSubmitted envelopes in DM history are ignored.
func TestAllNationsSubmitted_SkipsNonOrderSubmittedEvents(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)

	// Pre-populate u1's DM with a GameCreated event (wrong type) and a
	// malformed JSON message before the valid OrderSubmitted.
	_ = ch.SendDM("u1", "not-json-at-all")
	_ = events.WriteDM(ch, "u1", events.TypeGameCreated, events.GameCreated{Variant: "classical"})
	// Malformed OrderSubmitted payload.
	badEnv := events.Envelope{Type: events.TypeOrderSubmitted, Payload: json.RawMessage(`"bad"`)}
	badData, _ := json.Marshal(badEnv)
	_ = ch.SendDM("u1", string(badData))
	// Wrong phase.
	_ = events.WriteDM(ch, "u1", events.TypeOrderSubmitted, events.OrderSubmitted{
		Nation: "England", Phase: "Fall 1901 Movement",
	})

	sess := &session.Session{
		Phase:        "Spring 1901 Movement",
		Players:      map[string]string{"u1": "England"},
		StagedOrders: make(map[string][]string),
		Submitted:    make(map[string]bool),
	}

	done, err := d.allNationsSubmitted(sess)
	is.NoErr(err)
	is.Equal(done, false) // no valid OrderSubmitted for the current phase
}
