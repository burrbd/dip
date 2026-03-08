package bot

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/burrbd/dip/dipmap"
	"github.com/burrbd/dip/engine"
	"github.com/burrbd/dip/events"
	"github.com/burrbd/dip/session"
	"github.com/cheekybits/is"
)

// ---- mock channel -----------------------------------------------------------

type mockChannel struct {
	msgs         []string
	dms          map[string][]string
	imgs         [][]byte
	postErr      error
	histErr      error
	dmPostErr    error
	dmHistErr    error
	postErrAfter int // if > 0, Post returns an error after this many successful posts
}

func (m *mockChannel) Post(_, text string) error {
	if m.postErr != nil {
		return m.postErr
	}
	if m.postErrAfter > 0 && len(m.msgs) >= m.postErrAfter {
		return errors.New("post failed")
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

func (m *mockChannel) PostImage(_ string, data []byte) error {
	if m.postErr != nil {
		return m.postErr
	}
	m.imgs = append(m.imgs, data)
	return nil
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
	dislodgeds map[string]string
	units      map[string]engine.UnitInfo
}

func (e *mockEngine) SubmitOrder(_, _ string) error {
	return e.orderErr
}
func (e *mockEngine) Resolve() (engine.ResolutionResult, error) {
	return engine.ResolutionResult{Phase: e.phase}, e.resolveErr
}
func (e *mockEngine) Advance() error     { return e.advanceErr }
func (e *mockEngine) SoloWinner() string { return e.soloWinner }
func (e *mockEngine) Dump() ([]byte, error) { return e.dump, e.dumpErr }
func (e *mockEngine) Phase() string         { return e.phase }
func (e *mockEngine) Dislodgeds() map[string]string {
	if e.dislodgeds == nil {
		return make(map[string]string)
	}
	return e.dislodgeds
}
func (e *mockEngine) SupplyCenters() map[string]int { return make(map[string]int) }
func (e *mockEngine) Units() map[string]engine.UnitInfo {
	if e.units != nil {
		return e.units
	}
	return make(map[string]engine.UnitInfo)
}

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

// ---- helpers for retreat/adjustment sessions --------------------------------

// makeRetreatSession creates a session in the Retreat phase with one player (u1→England).
// The engine has a dislodged unit at "vie" (lowercase, matching godip's province keys) belonging to England.
func makeRetreatSession(d *Dispatcher, ch *mockChannel, gameChID string) *session.Session {
	players := map[string]string{"u1": "England"}
	eng := &mockEngine{
		phase:      "Spring 1901 Retreat",
		dump:       []byte(`{}`),
		dislodgeds: map[string]string{"vie": "England"},
	}
	sess := session.New(ch, gameChID, "gm1", "Spring 1901 Retreat", players, 0, eng, &mockNotifier{})
	d.sessions[gameChID] = sess
	return sess
}

// makeAdjustmentSession creates a session in the Adjustment phase with one player (u1→England).
func makeAdjustmentSession(d *Dispatcher, ch *mockChannel, gameChID string) *session.Session {
	players := map[string]string{"u1": "England"}
	eng := &mockEngine{
		phase: "Fall 1901 Adjustment",
		dump:  []byte(`{}`),
	}
	sess := session.New(ch, gameChID, "gm1", "Fall 1901 Adjustment", players, 0, eng, &mockNotifier{})
	d.sessions[gameChID] = sess
	return sess
}

// ---- /retreat ---------------------------------------------------------------

func TestDispatchRetreat_RejectsNonDM(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	makeRetreatSession(d, ch, "chan1")

	_, err := d.Dispatch(Command{Name: "retreat", Args: []string{"A", "Vie", "Bud"}, ChannelID: "chan1", UserID: "u1"})
	is.Err(err)
}

func TestDispatchRetreat_RejectsNoActiveGame(t *testing.T) {
	is := is.New(t)
	d := newTestDispatcher(&mockChannel{})

	_, err := d.Dispatch(dmCmd("retreat", "chan1", "u1", "A", "Vie", "Bud"))
	is.Err(err)
}

func TestDispatchRetreat_RejectsNonRetreatPhase(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	makeDMSession(d, ch, "chan1") // Movement phase

	_, err := d.Dispatch(dmCmd("retreat", "chan1", "u1", "A", "Vie", "Bud"))
	is.Err(err)
}

func TestDispatchRetreat_RejectsNonPlayer(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	makeRetreatSession(d, ch, "chan1")

	_, err := d.Dispatch(dmCmd("retreat", "chan1", "outsider", "A", "Vie", "Bud"))
	is.Err(err)
}

func TestDispatchRetreat_RejectsMissingArgs(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	makeRetreatSession(d, ch, "chan1")

	_, err := d.Dispatch(dmCmd("retreat", "chan1", "u1", "A", "Vie")) // missing destination
	is.Err(err)
}

func TestDispatchRetreat_RejectsWrongNationDislodged(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	sess := makeRetreatSession(d, ch, "chan1")
	// Override dislodgeds so vie belongs to France, not England.
	sess.Eng = &mockEngine{
		phase:      "Spring 1901 Retreat",
		dump:       []byte(`{}`),
		dislodgeds: map[string]string{"vie": "France"},
	}

	_, err := d.Dispatch(dmCmd("retreat", "chan1", "u1", "A", "Vie", "Bud"))
	is.Err(err)
}

func TestDispatchRetreat_RejectsNoDislodgedAtProvince(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	sess := makeRetreatSession(d, ch, "chan1")
	// No dislodged unit at mun.
	sess.Eng = &mockEngine{
		phase:      "Spring 1901 Retreat",
		dump:       []byte(`{}`),
		dislodgeds: map[string]string{"vie": "England"},
	}

	_, err := d.Dispatch(dmCmd("retreat", "chan1", "u1", "A", "Mun", "Boh"))
	is.Err(err)
}

func TestDispatchRetreat_RejectsEngineError(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	sess := makeRetreatSession(d, ch, "chan1")
	sess.Eng = &mockEngine{
		phase:      "Spring 1901 Retreat",
		dump:       []byte(`{}`),
		dislodgeds: map[string]string{"vie": "England"},
		orderErr:   errors.New("invalid retreat destination"),
	}

	_, err := d.Dispatch(dmCmd("retreat", "chan1", "u1", "A", "Vie", "Bad"))
	is.Err(err)
}

func TestDispatchRetreat_StagesRetreatOrder(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	sess := makeRetreatSession(d, ch, "chan1")

	resp, err := d.Dispatch(dmCmd("retreat", "chan1", "u1", "A", "Vie", "Bud"))
	is.NoErr(err)
	if resp == "" {
		t.Error("expected non-empty response")
	}
	is.Equal(len(sess.StagedOrders["England"]), 1)
	is.Equal(sess.StagedOrders["England"][0], "A vie-Bud")
}

func TestDispatchRetreat_WriteDMError(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{dmPostErr: errors.New("dm write error")}
	d := newTestDispatcher(ch)
	makeRetreatSession(d, ch, "chan1")

	_, err := d.Dispatch(dmCmd("retreat", "chan1", "u1", "A", "Vie", "Bud"))
	is.Err(err)
}

// ---- /disband ---------------------------------------------------------------

func TestDispatchDisband_RejectsNonDM(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	makeRetreatSession(d, ch, "chan1")

	_, err := d.Dispatch(Command{Name: "disband", Args: []string{"A", "Vie"}, ChannelID: "chan1", UserID: "u1"})
	is.Err(err)
}

func TestDispatchDisband_RejectsNoActiveGame(t *testing.T) {
	is := is.New(t)
	d := newTestDispatcher(&mockChannel{})

	_, err := d.Dispatch(dmCmd("disband", "chan1", "u1", "A", "Vie"))
	is.Err(err)
}

func TestDispatchDisband_RejectsMovementPhase(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	makeDMSession(d, ch, "chan1") // Movement phase

	_, err := d.Dispatch(dmCmd("disband", "chan1", "u1", "A", "Vie"))
	is.Err(err)
}

func TestDispatchDisband_RejectsNonPlayer(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	makeRetreatSession(d, ch, "chan1")

	_, err := d.Dispatch(dmCmd("disband", "chan1", "outsider", "A", "Vie"))
	is.Err(err)
}

func TestDispatchDisband_RejectsMissingArgs(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	makeRetreatSession(d, ch, "chan1")

	_, err := d.Dispatch(dmCmd("disband", "chan1", "u1", "A")) // missing province
	is.Err(err)
}

func TestDispatchDisband_RejectsWrongNationDislodgedInRetreat(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	sess := makeRetreatSession(d, ch, "chan1")
	sess.Eng = &mockEngine{
		phase:      "Spring 1901 Retreat",
		dump:       []byte(`{}`),
		dislodgeds: map[string]string{"vie": "France"},
	}

	_, err := d.Dispatch(dmCmd("disband", "chan1", "u1", "A", "Vie"))
	is.Err(err)
}

func TestDispatchDisband_RejectsEngineError(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	sess := makeRetreatSession(d, ch, "chan1")
	sess.Eng = &mockEngine{
		phase:      "Spring 1901 Retreat",
		dump:       []byte(`{}`),
		dislodgeds: map[string]string{"vie": "England"},
		orderErr:   errors.New("bad disband"),
	}

	_, err := d.Dispatch(dmCmd("disband", "chan1", "u1", "A", "Vie"))
	is.Err(err)
}

func TestDispatchDisband_StagesDisbandInRetreatPhase(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	sess := makeRetreatSession(d, ch, "chan1")

	resp, err := d.Dispatch(dmCmd("disband", "chan1", "u1", "A", "Vie"))
	is.NoErr(err)
	if resp == "" {
		t.Error("expected non-empty response")
	}
	is.Equal(len(sess.StagedOrders["England"]), 1)
	is.Equal(sess.StagedOrders["England"][0], "A vie disband")
}

func TestDispatchDisband_StagesDisbandInAdjustmentPhase(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	sess := makeAdjustmentSession(d, ch, "chan1")

	resp, err := d.Dispatch(dmCmd("disband", "chan1", "u1", "A", "Lon"))
	is.NoErr(err)
	if resp == "" {
		t.Error("expected non-empty response")
	}
	is.Equal(len(sess.StagedOrders["England"]), 1)
	is.Equal(sess.StagedOrders["England"][0], "A lon disband")
}

func TestDispatchDisband_WriteDMErrorInRetreat(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{dmPostErr: errors.New("dm write error")}
	d := newTestDispatcher(ch)
	makeRetreatSession(d, ch, "chan1")

	_, err := d.Dispatch(dmCmd("disband", "chan1", "u1", "A", "Vie"))
	is.Err(err)
}

func TestDispatchDisband_WriteDMErrorInAdjustment(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{dmPostErr: errors.New("dm write error")}
	d := newTestDispatcher(ch)
	makeAdjustmentSession(d, ch, "chan1")

	_, err := d.Dispatch(dmCmd("disband", "chan1", "u1", "A", "Lon"))
	is.Err(err)
}

// ---- /build -----------------------------------------------------------------

func TestDispatchBuild_RejectsNonDM(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	makeAdjustmentSession(d, ch, "chan1")

	_, err := d.Dispatch(Command{Name: "build", Args: []string{"A", "Lon"}, ChannelID: "chan1", UserID: "u1"})
	is.Err(err)
}

func TestDispatchBuild_RejectsNoActiveGame(t *testing.T) {
	is := is.New(t)
	d := newTestDispatcher(&mockChannel{})

	_, err := d.Dispatch(dmCmd("build", "chan1", "u1", "A", "Lon"))
	is.Err(err)
}

func TestDispatchBuild_RejectsNonAdjustmentPhase(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	makeDMSession(d, ch, "chan1") // Movement phase

	_, err := d.Dispatch(dmCmd("build", "chan1", "u1", "A", "Lon"))
	is.Err(err)
}

func TestDispatchBuild_RejectsNonPlayer(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	makeAdjustmentSession(d, ch, "chan1")

	_, err := d.Dispatch(dmCmd("build", "chan1", "outsider", "A", "Lon"))
	is.Err(err)
}

func TestDispatchBuild_RejectsMissingArgs(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	makeAdjustmentSession(d, ch, "chan1")

	_, err := d.Dispatch(dmCmd("build", "chan1", "u1", "A")) // missing province
	is.Err(err)
}

func TestDispatchBuild_RejectsEngineError(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	sess := makeAdjustmentSession(d, ch, "chan1")
	sess.Eng = &mockEngine{
		phase:    "Fall 1901 Adjustment",
		dump:     []byte(`{}`),
		orderErr: errors.New("no build slot available"),
	}

	_, err := d.Dispatch(dmCmd("build", "chan1", "u1", "A", "Lon"))
	is.Err(err)
}

func TestDispatchBuild_StagesBuildOrder(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	sess := makeAdjustmentSession(d, ch, "chan1")

	resp, err := d.Dispatch(dmCmd("build", "chan1", "u1", "A", "Lon"))
	is.NoErr(err)
	if resp == "" {
		t.Error("expected non-empty response")
	}
	is.Equal(len(sess.StagedOrders["England"]), 1)
	is.Equal(sess.StagedOrders["England"][0], "build A Lon")
}

func TestDispatchBuild_WriteDMError(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{dmPostErr: errors.New("dm write error")}
	d := newTestDispatcher(ch)
	makeAdjustmentSession(d, ch, "chan1")

	_, err := d.Dispatch(dmCmd("build", "chan1", "u1", "A", "Lon"))
	is.Err(err)
}

// ---- /waive -----------------------------------------------------------------

func TestDispatchWaive_RejectsNonDM(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	makeAdjustmentSession(d, ch, "chan1")

	_, err := d.Dispatch(Command{Name: "waive", ChannelID: "chan1", UserID: "u1"})
	is.Err(err)
}

func TestDispatchWaive_RejectsNoActiveGame(t *testing.T) {
	is := is.New(t)
	d := newTestDispatcher(&mockChannel{})

	_, err := d.Dispatch(dmCmd("waive", "chan1", "u1"))
	is.Err(err)
}

func TestDispatchWaive_RejectsNonAdjustmentPhase(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	makeDMSession(d, ch, "chan1") // Movement phase

	_, err := d.Dispatch(dmCmd("waive", "chan1", "u1"))
	is.Err(err)
}

func TestDispatchWaive_RejectsNonPlayer(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	makeAdjustmentSession(d, ch, "chan1")

	_, err := d.Dispatch(dmCmd("waive", "chan1", "outsider"))
	is.Err(err)
}

func TestDispatchWaive_StagesWaiveOrder(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	sess := makeAdjustmentSession(d, ch, "chan1")

	resp, err := d.Dispatch(dmCmd("waive", "chan1", "u1"))
	is.NoErr(err)
	is.Equal(resp, "Waive order staged.")
	is.Equal(len(sess.StagedOrders["England"]), 1)
	is.Equal(sess.StagedOrders["England"][0], "Waive")
}

// ---- /retreat and /disband reject Adjustment / Retreat guard respectively ---

func TestDispatchRetreat_RejectsAdjustmentPhase(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	makeAdjustmentSession(d, ch, "chan1")

	_, err := d.Dispatch(dmCmd("retreat", "chan1", "u1", "A", "Vie", "Bud"))
	is.Err(err)
}

func TestDispatchBuild_RejectsRetreatPhase(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	makeRetreatSession(d, ch, "chan1")

	_, err := d.Dispatch(dmCmd("build", "chan1", "u1", "A", "Lon"))
	is.Err(err)
}

func TestDispatchWaive_RejectsRetreatPhase(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	makeRetreatSession(d, ch, "chan1")

	_, err := d.Dispatch(dmCmd("waive", "chan1", "u1"))
	is.Err(err)
}

// ---- /status ----------------------------------------------------------------

func TestDispatchStatus_RejectsNoSession(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)

	_, err := d.Dispatch(Command{Name: "status", ChannelID: "chan1", UserID: "u1"})
	is.Err(err)
}

func TestDispatchStatus_ReturnsPhaseAndNations(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	makeDMSession(d, ch, "chan1")

	resp, err := d.Dispatch(Command{Name: "status", ChannelID: "chan1", UserID: "u1"})
	is.NoErr(err)
	if !containsStr(resp, "Spring 1901 Movement") {
		t.Errorf("expected phase in status output, got: %q", resp)
	}
	if !containsStr(resp, "England") {
		t.Errorf("expected England in status output, got: %q", resp)
	}
}

// ---- /history ---------------------------------------------------------------

func TestDispatchHistory_RejectsMissingArg(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)

	_, err := d.Dispatch(Command{Name: "history", ChannelID: "chan1", UserID: "u1"})
	is.Err(err)
}

func TestDispatchHistory_RejectsChannelError(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{histErr: errors.New("channel down")}
	d := newTestDispatcher(ch)

	_, err := d.Dispatch(Command{Name: "history", Args: []string{"Spring 1901"}, ChannelID: "chan1", UserID: "u1"})
	is.Err(err)
}

func TestDispatchHistory_RejectsNotFound(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)

	_, err := d.Dispatch(Command{Name: "history", Args: []string{"Spring 1901"}, ChannelID: "chan1", UserID: "u1"})
	is.Err(err)
}

func TestDispatchHistory_ReturnsResultSummary(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	_ = events.Write(ch, "chan1", events.TypePhaseResolved, events.PhaseResolved{
		Phase:         "Spring 1901 Movement",
		StateSnapshot: json.RawMessage(`{}`),
		ResultSummary: json.RawMessage(`"Spring 1901 resolved"`),
	})
	d := newTestDispatcher(ch)

	resp, err := d.Dispatch(Command{Name: "history", Args: []string{"Spring 1901"}, ChannelID: "chan1", UserID: "u1"})
	is.NoErr(err)
	if !containsStr(resp, "Spring 1901 resolved") {
		t.Errorf("expected result summary in response, got: %q", resp)
	}
}

func TestDispatchHistory_NoResultSummaryFallback(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	_ = events.Write(ch, "chan1", events.TypePhaseResolved, events.PhaseResolved{
		Phase:         "Spring 1901 Movement",
		StateSnapshot: json.RawMessage(`{}`),
	})
	d := newTestDispatcher(ch)

	resp, err := d.Dispatch(Command{Name: "history", Args: []string{"Spring 1901"}, ChannelID: "chan1", UserID: "u1"})
	is.NoErr(err)
	if !containsStr(resp, "Spring 1901 Movement") {
		t.Errorf("expected phase in fallback response, got: %q", resp)
	}
}

func TestDispatchHistory_SkipsMalformedPhaseResolved(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	// Valid one first (earliest in history).
	_ = events.Write(ch, "chan1", events.TypePhaseResolved, events.PhaseResolved{
		Phase:         "Fall 1901 Movement",
		StateSnapshot: json.RawMessage(`{}`),
		ResultSummary: json.RawMessage(`"Fall 1901 resolved"`),
	})
	// Malformed PhaseResolved envelope after (latest in history — hit first by reverse scan).
	bad := events.Envelope{Type: events.TypePhaseResolved, Payload: json.RawMessage(`"bad"`)}
	data, _ := json.Marshal(bad)
	ch.msgs = append(ch.msgs, string(data))
	d := newTestDispatcher(ch)

	resp, err := d.Dispatch(Command{Name: "history", Args: []string{"Fall 1901"}, ChannelID: "chan1", UserID: "u1"})
	is.NoErr(err)
	if !containsStr(resp, "Fall 1901 resolved") {
		t.Errorf("expected result in response, got: %q", resp)
	}
}

func TestDispatchHistory_RejectsNotFoundWhenPhaseDoesNotMatch(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	// Seed a valid PhaseResolved whose phase does not match the search term.
	_ = events.Write(ch, "chan1", events.TypePhaseResolved, events.PhaseResolved{
		Phase:         "Fall 1901 Movement",
		StateSnapshot: json.RawMessage(`{}`),
		ResultSummary: json.RawMessage(`"Fall 1901 resolved"`),
	})
	d := newTestDispatcher(ch)

	_, err := d.Dispatch(Command{Name: "history", Args: []string{"Spring 1901"}, ChannelID: "chan1", UserID: "u1"})
	is.Err(err)
}

// ---- /map -------------------------------------------------------------------

func TestDispatchMap_RejectsNoSession(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)

	_, err := d.Dispatch(Command{Name: "map", ChannelID: "chan1", UserID: "u1"})
	is.Err(err)
}

func TestDispatchMap_PostsImageNoArgs(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	makeDMSession(d, ch, "chan1")

	resp, err := d.Dispatch(Command{Name: "map", ChannelID: "chan1", UserID: "u1"})
	is.NoErr(err)
	is.Equal(resp, "Map posted.")
	is.Equal(len(ch.imgs), 1)
}

func TestDispatchMap_PostsImageWithTerritory(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	makeDMSession(d, ch, "chan1")

	resp, err := d.Dispatch(Command{Name: "map", Args: []string{"Vienna"}, ChannelID: "chan1", UserID: "u1"})
	is.NoErr(err)
	is.Equal(resp, "Map posted.")
	is.Equal(len(ch.imgs), 1)
}

func TestDispatchMap_PostsImageWithTerritoryAndRadius(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	makeDMSession(d, ch, "chan1")

	resp, err := d.Dispatch(Command{Name: "map", Args: []string{"Vienna", "1"}, ChannelID: "chan1", UserID: "u1"})
	is.NoErr(err)
	is.Equal(resp, "Map posted.")
	is.Equal(len(ch.imgs), 1)
}

func TestDispatchMap_RejectsInvalidRadius(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	makeDMSession(d, ch, "chan1")

	_, err := d.Dispatch(Command{Name: "map", Args: []string{"Vienna", "notanumber"}, ChannelID: "chan1", UserID: "u1"})
	is.Err(err)
}

func TestDispatchMap_RejectsPostImageError(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	makeDMSession(d, ch, "chan1")
	ch.postErr = errors.New("image post failed")

	_, err := d.Dispatch(Command{Name: "map", ChannelID: "chan1", UserID: "u1"})
	is.Err(err)
}

func TestDispatchMap_RejectsSVGLoadError(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	makeDMSession(d, ch, "chan1")
	// svgFn is called first, before the full-board / zoom split.
	d.svgFn = func(_ dipmap.EngineState) ([]byte, error) {
		return nil, errors.New("svg load failed")
	}

	_, err := d.Dispatch(Command{Name: "map", ChannelID: "chan1", UserID: "u1"})
	is.Err(err)
}

func TestDispatchMap_RejectsOverlayError(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	makeDMSession(d, ch, "chan1")
	// SVG loads fine; overlay injection fails.
	d.overlayFn = func(_ []byte, _ map[string]dipmap.Unit) ([]byte, error) {
		return nil, errors.New("overlay failed")
	}

	_, err := d.Dispatch(Command{Name: "map", ChannelID: "chan1", UserID: "u1"})
	is.Err(err)
}

func TestDispatchMap_RejectsPNGError(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	makeDMSession(d, ch, "chan1")
	// SVG load and overlay succeed; PNG conversion fails (full-board path).
	d.pngFn = func(_ []byte) ([]byte, error) {
		return nil, errors.New("png failed")
	}

	_, err := d.Dispatch(Command{Name: "map", ChannelID: "chan1", UserID: "u1"})
	is.Err(err)
}

func TestDispatchMap_RejectsHighlightError(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	makeDMSession(d, ch, "chan1")
	// SVG load and overlay succeed; highlight fails (zoomed path).
	d.highlightFn = func(_ []byte, _ []string) ([]byte, error) {
		return nil, errors.New("highlight failed")
	}

	_, err := d.Dispatch(Command{Name: "map", Args: []string{"Vienna", "1"}, ChannelID: "chan1", UserID: "u1"})
	is.Err(err)
}

func TestDispatchMap_OverlaysUnitsOnMap(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	sess := makeDMSession(d, ch, "chan1")
	// Seed an engine that reports a unit at "par".
	sess.Eng = &mockEngine{
		phase: "Spring 1901 Movement",
		dump:  []byte(`{}`),
		units: map[string]engine.UnitInfo{"par": {Type: "Army", Nation: "France"}},
	}

	var gotUnits map[string]dipmap.Unit
	d.overlayFn = func(svg []byte, units map[string]dipmap.Unit) ([]byte, error) {
		gotUnits = units
		return svg, nil
	}

	_, err := d.Dispatch(Command{Name: "map", ChannelID: "chan1", UserID: "u1"})
	is.NoErr(err)
	if gotUnits["par"].Type != "Army" || gotUnits["par"].Nation != "France" {
		t.Errorf("expected France Army at par, got %+v", gotUnits["par"])
	}
}

func TestDispatchMap_RejectsZoomError(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	makeDMSession(d, ch, "chan1")
	// SVG load and highlight succeed; zoom render fails.
	d.renderZoomedFn = func(_ dipmap.EngineState, _ []byte, _ []string) ([]byte, error) {
		return nil, errors.New("zoom failed")
	}

	_, err := d.Dispatch(Command{Name: "map", Args: []string{"Vienna", "1"}, ChannelID: "chan1", UserID: "u1"})
	is.Err(err)
}

func TestDispatchMap_UsesCustomGraph(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	makeDMSession(d, ch, "chan1")

	// Set a custom graph on the dispatcher.
	d.graph = customTestGraph{"Vienna": {"Budapest"}}

	resp, err := d.Dispatch(Command{Name: "map", Args: []string{"Vienna", "1"}, ChannelID: "chan1", UserID: "u1"})
	is.NoErr(err)
	is.Equal(resp, "Map posted.")
	is.Equal(len(ch.imgs), 1)
}

// customTestGraph is a simple Graph for use in bot tests.
type customTestGraph map[string][]string

func (g customTestGraph) Edges(t string) []string { return g[t] }

// ---- /history additional coverage ------------------------------------------

func TestDispatchHistory_SkipsNonPhaseResolvedEvents(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	// Valid PhaseResolved first (earliest in history).
	_ = events.Write(ch, "chan1", events.TypePhaseResolved, events.PhaseResolved{
		Phase:         "Spring 1901 Movement",
		StateSnapshot: json.RawMessage(`{}`),
		ResultSummary: json.RawMessage(`"resolved"`),
	})
	// GameCreated (not PhaseResolved) after — hit first by the reverse scan.
	seedGameCreated(ch, "gm1")
	d := newTestDispatcher(ch)

	resp, err := d.Dispatch(Command{Name: "history", Args: []string{"Spring 1901"}, ChannelID: "chan1", UserID: "u1"})
	is.NoErr(err)
	if !containsStr(resp, "resolved") {
		t.Errorf("expected result in response, got: %q", resp)
	}
}

// ---- /help ------------------------------------------------------------------

func TestDispatchHelp_ListsAllCommands(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)

	resp, err := d.Dispatch(Command{Name: "help", ChannelID: "chan1", UserID: "u1"})
	is.NoErr(err)
	if !containsStr(resp, "/newgame") {
		t.Errorf("expected /newgame in help output, got: %q", resp)
	}
	if !containsStr(resp, "/start") {
		t.Errorf("expected /start in help output, got: %q", resp)
	}
}

func TestDispatchHelp_ShowsSpecificCommand(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)

	resp, err := d.Dispatch(Command{Name: "help", Args: []string{"order"}, ChannelID: "chan1", UserID: "u1"})
	is.NoErr(err)
	if !containsStr(resp, "/order") {
		t.Errorf("expected /order in help output, got: %q", resp)
	}
}

func TestDispatchHelp_ShowsSpecificCommandWithSlash(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)

	resp, err := d.Dispatch(Command{Name: "help", Args: []string{"/join"}, ChannelID: "chan1", UserID: "u1"})
	is.NoErr(err)
	if !containsStr(resp, "/join") {
		t.Errorf("expected /join in help output, got: %q", resp)
	}
}

func TestDispatchHelp_RejectsUnknownCommand(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)

	_, err := d.Dispatch(Command{Name: "help", Args: []string{"bogus"}, ChannelID: "chan1", UserID: "u1"})
	is.Err(err)
}

// ---- helpers for draw/concede/GM command tests ------------------------------

// seedStartedGame writes GameCreated, PlayerJoined (for each player in players),
// and GameStarted events to ch, and registers a session in d.sessions keyed by
// channelID. DeadlineHours is 0 to prevent background timers during tests.
// players is a map of userID → nation (same format as session.Session.Players).
func seedStartedGame(d *Dispatcher, ch *mockChannel, channelID, gmID string, players map[string]string) *session.Session {
	_ = events.Write(ch, channelID, events.TypeGameCreated, events.GameCreated{
		Variant: "classical", DeadlineHours: 24, GMUserID: gmID,
	})
	for uid, nation := range players {
		_ = events.Write(ch, channelID, events.TypePlayerJoined, events.PlayerJoined{
			UserID: uid, Nation: nation,
		})
	}
	_ = events.Write(ch, channelID, events.TypeGameStarted, events.GameStarted{
		InitialState: json.RawMessage(`{}`),
	})
	eng := goodEngine()
	sess := session.New(ch, channelID, gmID, "Spring 1901 Movement", players, 0, eng, &mockNotifier{})
	d.sessions[channelID] = sess
	return sess
}

// twoPlayerGame seeds a two-player game (u1→England, u2→France) run by gm1.
func twoPlayerGame(d *Dispatcher, ch *mockChannel) *session.Session {
	players := map[string]string{"u1": "England", "u2": "France"}
	return seedStartedGame(d, ch, "chan1", "gm1", players)
}

// gameCmd builds a non-DM Command for the game channel.
func gameCmd(name, channelID, userID string, args ...string) Command {
	return Command{
		Name:      name,
		Args:      args,
		UserID:    userID,
		ChannelID: channelID,
	}
}

// ---- /draw ------------------------------------------------------------------

func TestDispatchDraw_ProposesDraw_PostsDrawProposed(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	twoPlayerGame(d, ch)

	_, err := d.Dispatch(gameCmd("draw", "chan1", "u1"))
	is.NoErr(err)
	is.Equal(ch.lastEventType(), events.TypeDrawProposed)
}

func TestDispatchDraw_ProposesDraw_ReturnsNonEmptyMessage(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	twoPlayerGame(d, ch)

	resp, err := d.Dispatch(gameCmd("draw", "chan1", "u1"))
	is.NoErr(err)
	if resp == "" {
		t.Error("expected non-empty response")
	}
}

func TestDispatchDraw_VoteYes_PostsDrawVoted(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	twoPlayerGame(d, ch)
	// Seed an existing draw proposal from England.
	_ = events.Write(ch, "chan1", events.TypeDrawProposed, events.DrawProposed{ProposerNation: "England"})

	_, err := d.Dispatch(gameCmd("draw", "chan1", "u2")) // France votes yes
	is.NoErr(err)
	is.Equal(ch.lastEventType(), events.TypeGameEnded) // both voted → game ended
}

func TestDispatchDraw_AllNationsAgree_PostsGameEnded(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	twoPlayerGame(d, ch)
	_ = events.Write(ch, "chan1", events.TypeDrawProposed, events.DrawProposed{ProposerNation: "England"})

	_, err := d.Dispatch(gameCmd("draw", "chan1", "u2"))
	is.NoErr(err)

	var env events.Envelope
	is.NoErr(json.Unmarshal([]byte(ch.msgs[len(ch.msgs)-1]), &env))
	is.Equal(env.Type, events.TypeGameEnded)

	var ge events.GameEnded
	is.NoErr(json.Unmarshal(env.Payload, &ge))
	is.Equal(ge.Result, "draw")
}

func TestDispatchDraw_PartialVote_ReturnsCountMessage(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	// Three players; England proposes; France votes; Turkey still needs to vote.
	players := map[string]string{"u1": "England", "u2": "France", "u3": "Turkey"}
	seedStartedGame(d, ch, "chan1", "gm1", players)
	_ = events.Write(ch, "chan1", events.TypeDrawProposed, events.DrawProposed{ProposerNation: "England"})

	resp, err := d.Dispatch(gameCmd("draw", "chan1", "u2")) // France votes
	is.NoErr(err)
	if !containsStr(resp, "1") {
		t.Errorf("expected remaining nation count in response, got: %q", resp)
	}
}

func TestDispatchDraw_AlreadyVoted_ReturnsMessage(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	twoPlayerGame(d, ch)
	// England proposed and already voted yes.
	_ = events.Write(ch, "chan1", events.TypeDrawProposed, events.DrawProposed{ProposerNation: "England"})

	resp, err := d.Dispatch(gameCmd("draw", "chan1", "u1")) // England tries to vote again
	is.NoErr(err)
	if !containsStr(resp, "already") {
		t.Errorf("expected 'already voted' message, got: %q", resp)
	}
}

func TestDispatchDraw_RejectsIfGameNotStarted(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	seedGameCreated(ch, "gm1")

	_, err := d.Dispatch(gameCmd("draw", "chan1", "u1"))
	is.Err(err)
}

func TestDispatchDraw_RejectsIfGameEnded(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	twoPlayerGame(d, ch)
	_ = events.Write(ch, "chan1", events.TypeGameEnded, events.GameEnded{Result: "solo", Winner: "England"})

	_, err := d.Dispatch(gameCmd("draw", "chan1", "u1"))
	is.Err(err)
}

func TestDispatchDraw_RejectsIfNotPlayer(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	twoPlayerGame(d, ch)

	_, err := d.Dispatch(gameCmd("draw", "chan1", "outsider"))
	is.Err(err)
}

func TestDispatchDraw_RejectsIfChannelFails(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{histErr: errors.New("channel down")}
	d := newTestDispatcher(ch)

	_, err := d.Dispatch(gameCmd("draw", "chan1", "u1"))
	is.Err(err)
}

func TestDispatchDraw_RejectsIfProposalWriteFails(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	twoPlayerGame(d, ch)
	// The first 3 messages (GameCreated, PlayerJoined×2) were posted during seed;
	// GameStarted = 4th. Next post (DrawProposed) should fail.
	ch.postErrAfter = len(ch.msgs)

	_, err := d.Dispatch(gameCmd("draw", "chan1", "u1"))
	is.Err(err)
}

func TestDispatchDraw_RejectsIfVoteWriteFails(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	twoPlayerGame(d, ch)
	_ = events.Write(ch, "chan1", events.TypeDrawProposed, events.DrawProposed{ProposerNation: "England"})
	// Allow the DrawProposed (already written), but fail on DrawVoted.
	ch.postErrAfter = len(ch.msgs)

	_, err := d.Dispatch(gameCmd("draw", "chan1", "u2"))
	is.Err(err)
}

func TestDispatchDraw_RejectsIfGameEndedWriteFails(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	twoPlayerGame(d, ch)
	_ = events.Write(ch, "chan1", events.TypeDrawProposed, events.DrawProposed{ProposerNation: "England"})
	// Allow DrawVoted to succeed (1 more message), but fail on GameEnded.
	ch.postErrAfter = len(ch.msgs) + 1

	_, err := d.Dispatch(gameCmd("draw", "chan1", "u2"))
	is.Err(err)
}

func TestDispatchDraw_SinglePlayerImmediate(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	// Only one nation remains; proposing immediately ends the game.
	players := map[string]string{"u1": "England"}
	seedStartedGame(d, ch, "chan1", "gm1", players)

	_, err := d.Dispatch(gameCmd("draw", "chan1", "u1"))
	is.NoErr(err)
	is.Equal(ch.lastEventType(), events.TypeGameEnded)
}

func TestDispatchDraw_SinglePlayer_GameEndedWriteFails(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	players := map[string]string{"u1": "England"}
	seedStartedGame(d, ch, "chan1", "gm1", players)
	// DrawProposed may succeed (1 more), but GameEnded write must fail.
	ch.postErrAfter = len(ch.msgs) + 1

	_, err := d.Dispatch(gameCmd("draw", "chan1", "u1"))
	is.Err(err)
}

// ---- readState draw/booted/replaced coverage --------------------------------

func TestReadState_TracksDrawVoteYes(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	seedGameCreated(ch, "gm1")
	_ = events.Write(ch, "chan1", events.TypeDrawProposed, events.DrawProposed{ProposerNation: "England"})
	_ = events.Write(ch, "chan1", events.TypeDrawVoted, events.DrawVoted{Nation: "France", Accept: true})
	d := newTestDispatcher(ch)

	state, err := d.readState("chan1")
	is.NoErr(err)
	is.Equal(state.drawVotes["England"], true)
	is.Equal(state.drawVotes["France"], true)
}

func TestReadState_TracksDrawVoteNo(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	seedGameCreated(ch, "gm1")
	_ = events.Write(ch, "chan1", events.TypeDrawProposed, events.DrawProposed{ProposerNation: "England"})
	// England withdraws their vote.
	_ = events.Write(ch, "chan1", events.TypeDrawVoted, events.DrawVoted{Nation: "England", Accept: false})
	d := newTestDispatcher(ch)

	state, err := d.readState("chan1")
	is.NoErr(err)
	is.Equal(state.drawProposed, true)
	is.Equal(len(state.drawVotes), 0) // England's vote was withdrawn
}

func TestReadState_TracksGameEnded(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	seedGameCreated(ch, "gm1")
	_ = events.Write(ch, "chan1", events.TypeGameStarted, events.GameStarted{InitialState: json.RawMessage(`{}`)})
	_ = events.Write(ch, "chan1", events.TypeGameEnded, events.GameEnded{Result: "solo", Winner: "England"})
	d := newTestDispatcher(ch)

	state, err := d.readState("chan1")
	is.NoErr(err)
	is.Equal(state.ended, true)
}

func TestReadState_GameEndedResetsDraw(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	seedGameCreated(ch, "gm1")
	_ = events.Write(ch, "chan1", events.TypeDrawProposed, events.DrawProposed{ProposerNation: "England"})
	_ = events.Write(ch, "chan1", events.TypeGameEnded, events.GameEnded{Result: "draw"})
	d := newTestDispatcher(ch)

	state, err := d.readState("chan1")
	is.NoErr(err)
	is.Equal(state.drawProposed, false)
	is.Equal(len(state.drawVotes), 0)
}

func TestReadState_TracksBootedPlayers(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	seedGameCreated(ch, "gm1")
	_ = events.Write(ch, "chan1", events.TypePlayerJoined, events.PlayerJoined{UserID: "u1", Nation: "England"})
	_ = events.Write(ch, "chan1", events.TypePlayerBooted, events.PlayerBooted{Nation: "England"})
	d := newTestDispatcher(ch)

	state, err := d.readState("chan1")
	is.NoErr(err)
	_, found := state.players["u1"]
	is.Equal(found, false)
	_, nationFound := state.nations["England"]
	is.Equal(nationFound, false)
}

func TestReadState_TracksReplacedPlayers(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	seedGameCreated(ch, "gm1")
	_ = events.Write(ch, "chan1", events.TypePlayerJoined, events.PlayerJoined{UserID: "u1", Nation: "England"})
	_ = events.Write(ch, "chan1", events.TypePlayerReplaced, events.PlayerReplaced{Nation: "England", NewUserID: "u99"})
	d := newTestDispatcher(ch)

	state, err := d.readState("chan1")
	is.NoErr(err)
	_, oldFound := state.players["u1"]
	is.Equal(oldFound, false)
	is.Equal(state.players["u99"], "England")
	is.Equal(state.nations["England"], "u99")
}

func TestReadState_SkipsMalformedDrawProposed(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	seedGameCreated(ch, "gm1")
	seedMalformed(ch, "chan1", events.TypeDrawProposed)
	// A valid DrawProposed after the bad one.
	_ = events.Write(ch, "chan1", events.TypeDrawProposed, events.DrawProposed{ProposerNation: "England"})
	d := newTestDispatcher(ch)

	state, err := d.readState("chan1")
	is.NoErr(err)
	is.Equal(state.drawProposed, true)
}

func TestReadState_SkipsMalformedDrawVoted(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	seedGameCreated(ch, "gm1")
	_ = events.Write(ch, "chan1", events.TypeDrawProposed, events.DrawProposed{ProposerNation: "England"})
	seedMalformed(ch, "chan1", events.TypeDrawVoted)
	d := newTestDispatcher(ch)

	state, err := d.readState("chan1")
	is.NoErr(err)
	is.Equal(len(state.drawVotes), 1) // only England's vote from DrawProposed
}

func TestReadState_SkipsMalformedPlayerBooted(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	seedGameCreated(ch, "gm1")
	_ = events.Write(ch, "chan1", events.TypePlayerJoined, events.PlayerJoined{UserID: "u1", Nation: "England"})
	seedMalformed(ch, "chan1", events.TypePlayerBooted)
	d := newTestDispatcher(ch)

	state, err := d.readState("chan1")
	is.NoErr(err)
	is.Equal(state.players["u1"], "England") // bad event → player still present
}

func TestReadState_SkipsMalformedPlayerReplaced(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	seedGameCreated(ch, "gm1")
	_ = events.Write(ch, "chan1", events.TypePlayerJoined, events.PlayerJoined{UserID: "u1", Nation: "England"})
	seedMalformed(ch, "chan1", events.TypePlayerReplaced)
	d := newTestDispatcher(ch)

	state, err := d.readState("chan1")
	is.NoErr(err)
	is.Equal(state.players["u1"], "England") // bad event → player still original
}

// ---- /concede ---------------------------------------------------------------

func TestDispatchConcede_PostsGameEnded(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	twoPlayerGame(d, ch)

	_, err := d.Dispatch(gameCmd("concede", "chan1", "u1"))
	is.NoErr(err)
	is.Equal(ch.lastEventType(), events.TypeGameEnded)
}

func TestDispatchConcede_GameEndedHasConcessionResult(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	twoPlayerGame(d, ch)

	_, err := d.Dispatch(gameCmd("concede", "chan1", "u1"))
	is.NoErr(err)

	var env events.Envelope
	is.NoErr(json.Unmarshal([]byte(ch.msgs[len(ch.msgs)-1]), &env))
	var ge events.GameEnded
	is.NoErr(json.Unmarshal(env.Payload, &ge))
	is.Equal(ge.Result, "concession")
	is.Equal(ge.Winner, "England")
}

func TestDispatchConcede_RejectsIfNotStarted(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	seedGameCreated(ch, "gm1")

	_, err := d.Dispatch(gameCmd("concede", "chan1", "u1"))
	is.Err(err)
}

func TestDispatchConcede_RejectsIfGameEnded(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	twoPlayerGame(d, ch)
	_ = events.Write(ch, "chan1", events.TypeGameEnded, events.GameEnded{Result: "solo", Winner: "England"})

	_, err := d.Dispatch(gameCmd("concede", "chan1", "u1"))
	is.Err(err)
}

func TestDispatchConcede_RejectsIfNotPlayer(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	twoPlayerGame(d, ch)

	_, err := d.Dispatch(gameCmd("concede", "chan1", "outsider"))
	is.Err(err)
}

func TestDispatchConcede_RejectsIfChannelFails(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{histErr: errors.New("channel down")}
	d := newTestDispatcher(ch)

	_, err := d.Dispatch(gameCmd("concede", "chan1", "u1"))
	is.Err(err)
}

func TestDispatchConcede_RejectsIfWriteFails(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	twoPlayerGame(d, ch)
	ch.postErrAfter = len(ch.msgs)

	_, err := d.Dispatch(gameCmd("concede", "chan1", "u1"))
	is.Err(err)
}

// ---- /pause -----------------------------------------------------------------

func TestDispatchPause_CancelsDeadline(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	twoPlayerGame(d, ch)

	resp, err := d.Dispatch(gameCmd("pause", "chan1", "gm1"))
	is.NoErr(err)
	if resp == "" {
		t.Error("expected non-empty response")
	}
}

func TestDispatchPause_RejectsIfNotGM(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	twoPlayerGame(d, ch)

	_, err := d.Dispatch(gameCmd("pause", "chan1", "u1"))
	is.Err(err)
}

func TestDispatchPause_RejectsIfNoSession(t *testing.T) {
	is := is.New(t)
	d := newTestDispatcher(&mockChannel{})

	_, err := d.Dispatch(gameCmd("pause", "chan1", "gm1"))
	is.Err(err)
}

// ---- /resume ----------------------------------------------------------------

func TestDispatchResume_RestartsDeadline(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	twoPlayerGame(d, ch)

	resp, err := d.Dispatch(gameCmd("resume", "chan1", "gm1"))
	is.NoErr(err)
	if resp == "" {
		t.Error("expected non-empty response")
	}
}

func TestDispatchResume_RejectsIfNotGM(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	twoPlayerGame(d, ch)

	_, err := d.Dispatch(gameCmd("resume", "chan1", "u1"))
	is.Err(err)
}

func TestDispatchResume_RejectsIfNoSession(t *testing.T) {
	is := is.New(t)
	d := newTestDispatcher(&mockChannel{})

	_, err := d.Dispatch(gameCmd("resume", "chan1", "gm1"))
	is.Err(err)
}

// ---- /extend ----------------------------------------------------------------

func TestDispatchExtend_ExtendsDeadline(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	twoPlayerGame(d, ch)

	resp, err := d.Dispatch(gameCmd("extend", "chan1", "gm1", "2h"))
	is.NoErr(err)
	if !containsStr(resp, "2h") {
		t.Errorf("expected duration in response, got: %q", resp)
	}
}

func TestDispatchExtend_RejectsIfNotGM(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	twoPlayerGame(d, ch)

	_, err := d.Dispatch(gameCmd("extend", "chan1", "u1", "2h"))
	is.Err(err)
}

func TestDispatchExtend_RejectsIfNoSession(t *testing.T) {
	is := is.New(t)
	d := newTestDispatcher(&mockChannel{})

	_, err := d.Dispatch(gameCmd("extend", "chan1", "gm1", "2h"))
	is.Err(err)
}

func TestDispatchExtend_RejectsMissingArg(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	twoPlayerGame(d, ch)

	_, err := d.Dispatch(gameCmd("extend", "chan1", "gm1"))
	is.Err(err)
}

func TestDispatchExtend_RejectsInvalidDuration(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	twoPlayerGame(d, ch)

	_, err := d.Dispatch(gameCmd("extend", "chan1", "gm1", "notaduration"))
	is.Err(err)
}

// ---- /force-resolve ---------------------------------------------------------

func TestDispatchForceResolve_CallsAdvanceTurn(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	twoPlayerGame(d, ch)

	resp, err := d.Dispatch(gameCmd("force-resolve", "chan1", "gm1"))
	is.NoErr(err)
	if resp == "" {
		t.Error("expected non-empty response")
	}
	// PhaseResolved must have been posted.
	is.Equal(ch.lastEventType(), events.TypePhaseResolved)
}

func TestDispatchForceResolve_RejectsIfNotGM(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	twoPlayerGame(d, ch)

	_, err := d.Dispatch(gameCmd("force-resolve", "chan1", "u1"))
	is.Err(err)
}

func TestDispatchForceResolve_RejectsIfNoSession(t *testing.T) {
	is := is.New(t)
	d := newTestDispatcher(&mockChannel{})

	_, err := d.Dispatch(gameCmd("force-resolve", "chan1", "gm1"))
	is.Err(err)
}

func TestDispatchForceResolve_RejectsIfAdvanceTurnFails(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	sess := twoPlayerGame(d, ch)
	sess.Eng = &mockEngine{
		phase:      "Spring 1901 Movement",
		dump:       []byte(`{}`),
		resolveErr: errors.New("resolve failed"),
	}

	_, err := d.Dispatch(gameCmd("force-resolve", "chan1", "gm1"))
	is.Err(err)
}

// ---- /boot ------------------------------------------------------------------

func TestDispatchBoot_RemovesPlayer(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	sess := twoPlayerGame(d, ch)

	_, err := d.Dispatch(gameCmd("boot", "chan1", "gm1", "England"))
	is.NoErr(err)
	is.Equal(ch.lastEventType(), events.TypePlayerBooted)
	_, stillThere := sess.Players["u1"]
	is.Equal(stillThere, false)
}

func TestDispatchBoot_RejectsIfNotGM(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	twoPlayerGame(d, ch)

	_, err := d.Dispatch(gameCmd("boot", "chan1", "u1", "France"))
	is.Err(err)
}

func TestDispatchBoot_RejectsIfNoSession(t *testing.T) {
	is := is.New(t)
	d := newTestDispatcher(&mockChannel{})

	_, err := d.Dispatch(gameCmd("boot", "chan1", "gm1", "England"))
	is.Err(err)
}

func TestDispatchBoot_RejectsMissingArg(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	twoPlayerGame(d, ch)

	_, err := d.Dispatch(gameCmd("boot", "chan1", "gm1"))
	is.Err(err)
}

func TestDispatchBoot_RejectsIfNationNotFound(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	twoPlayerGame(d, ch)

	_, err := d.Dispatch(gameCmd("boot", "chan1", "gm1", "Russia"))
	is.Err(err)
}

func TestDispatchBoot_RejectsIfWriteFails(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	twoPlayerGame(d, ch)
	ch.postErrAfter = len(ch.msgs)

	_, err := d.Dispatch(gameCmd("boot", "chan1", "gm1", "England"))
	is.Err(err)
}

// ---- /replace ---------------------------------------------------------------

func TestDispatchReplace_TransfersNation(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	sess := twoPlayerGame(d, ch)

	_, err := d.Dispatch(gameCmd("replace", "chan1", "gm1", "England", "newuser"))
	is.NoErr(err)
	is.Equal(ch.lastEventType(), events.TypePlayerReplaced)
	_, oldGone := sess.Players["u1"]
	is.Equal(oldGone, false)
	is.Equal(sess.Players["newuser"], "England")
}

func TestDispatchReplace_RejectsIfNotGM(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	twoPlayerGame(d, ch)

	_, err := d.Dispatch(gameCmd("replace", "chan1", "u1", "France", "newuser"))
	is.Err(err)
}

func TestDispatchReplace_RejectsIfNoSession(t *testing.T) {
	is := is.New(t)
	d := newTestDispatcher(&mockChannel{})

	_, err := d.Dispatch(gameCmd("replace", "chan1", "gm1", "England", "newuser"))
	is.Err(err)
}

func TestDispatchReplace_RejectsMissingArgs(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	twoPlayerGame(d, ch)

	_, err := d.Dispatch(gameCmd("replace", "chan1", "gm1", "England")) // missing new user
	is.Err(err)
}

func TestDispatchReplace_RejectsIfNationNotFound(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	twoPlayerGame(d, ch)

	_, err := d.Dispatch(gameCmd("replace", "chan1", "gm1", "Russia", "newuser"))
	is.Err(err)
}

func TestDispatchReplace_RejectsIfWriteFails(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	d := newTestDispatcher(ch)
	twoPlayerGame(d, ch)
	ch.postErrAfter = len(ch.msgs)

	_, err := d.Dispatch(gameCmd("replace", "chan1", "gm1", "England", "newuser"))
	is.Err(err)
}
