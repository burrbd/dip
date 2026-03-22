package session

import (
	"encoding/json"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/burrbd/dip/engine"
	"github.com/burrbd/dip/events"
	"github.com/cheekybits/is"
)

// ---- mock channel -----------------------------------------------------------

type mockChannel struct {
	mu      sync.Mutex
	msgs    []string
	postErr error
	histErr error
}

func (m *mockChannel) Post(_, text string) error {
	if m.postErr != nil {
		return m.postErr
	}
	m.mu.Lock()
	m.msgs = append(m.msgs, text)
	m.mu.Unlock()
	return nil
}

func (m *mockChannel) History(_ string) ([]string, error) {
	if m.histErr != nil {
		return nil, m.histErr
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.msgs, nil
}

func (m *mockChannel) SendDM(_, _ string) error             { return nil }
func (m *mockChannel) DMHistory(_ string) ([]string, error) { return nil, nil }
func (m *mockChannel) PostImage(_ string, _ []byte) error   { return nil }

func (m *mockChannel) msgCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.msgs)
}

func (m *mockChannel) msgAt(i int) string {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.msgs[i]
}

// ---- mock notifier ----------------------------------------------------------

type mockNotifier struct {
	mu    sync.Mutex
	calls []string
	err   error
}

func (n *mockNotifier) Notify(_, msg string) error {
	n.mu.Lock()
	n.calls = append(n.calls, msg)
	n.mu.Unlock()
	return n.err
}

func (n *mockNotifier) callCount() int {
	n.mu.Lock()
	defer n.mu.Unlock()
	return len(n.calls)
}

// ---- mock engine ------------------------------------------------------------

type mockEngine struct {
	resolveResult engine.ResolutionResult
	resolveErr    error
	advanceErr    error
	dumpData      []byte
	dumpErr       error
	soloWinner    string
	phaseStr      string
	submitErr     error
}

func (e *mockEngine) SubmitOrder(_, _ string) error             { return e.submitErr }
func (e *mockEngine) Resolve() (engine.ResolutionResult, error) { return e.resolveResult, e.resolveErr }
func (e *mockEngine) Advance() error                            { return e.advanceErr }
func (e *mockEngine) SoloWinner() string                        { return e.soloWinner }
func (e *mockEngine) Dump() ([]byte, error)                     { return e.dumpData, e.dumpErr }
func (e *mockEngine) Phase() string                             { return e.phaseStr }
func (e *mockEngine) Dislodgeds() map[string]string             { return make(map[string]string) }
func (e *mockEngine) SupplyCenters() map[string]int             { return make(map[string]int) }
func (e *mockEngine) Units() map[string]engine.UnitInfo         { return make(map[string]engine.UnitInfo) }
func (e *mockEngine) ValidRetreats() map[string][]string        { return make(map[string][]string) }
func (e *mockEngine) BuildOptions() map[string]engine.BuildOption {
	return make(map[string]engine.BuildOption)
}

// ---- helpers ----------------------------------------------------------------

func defaultEng() *mockEngine {
	return &mockEngine{
		resolveResult: engine.ResolutionResult{Phase: "Spring 1901 Movement"},
		dumpData:      []byte(`{}`),
	}
}

func makeLoader(eng engine.Engine) EngineLoader {
	return func(_ []byte) (engine.Engine, error) { return eng, nil }
}

func makeSession(ch events.Channel, eng engine.Engine, notifier Notifier) *Session {
	return &Session{
		ChannelID:     "chan1",
		Phase:         "Spring 1901 Movement",
		StagedOrders:  make(map[string][]string),
		Players:       make(map[string]string),
		Submitted:     make(map[string]bool),
		Eng:           eng,
		ch:            ch,
		notifier:      notifier,
		DeadlineHours: 24,
	}
}

// ---- New constructor tests --------------------------------------------------

func TestNew_SetsFields(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	eng := defaultEng()
	notifier := &mockNotifier{}
	players := map[string]string{"u1": "England", "u2": "France"}

	s := New(ch, "chan1", "gm1", "Spring 1901 Movement", players, 0, eng, notifier)

	is.Equal(s.ChannelID, "chan1")
	is.Equal(s.GMID, "gm1")
	is.Equal(s.Phase, "Spring 1901 Movement")
	is.Equal(s.DeadlineHours, 0)
	is.Equal(s.Players["u1"], "England")
	is.Equal(s.Players["u2"], "France")
	is.Equal(len(s.StagedOrders), 0)
	is.Equal(len(s.Submitted), 0)
}

func TestNew_CopiesPlayers(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	players := map[string]string{"u1": "England"}

	s := New(ch, "chan1", "gm1", "Spring 1901 Movement", players, 0, defaultEng(), nil)

	// Mutate source map; session copy must be unaffected.
	players["u2"] = "France"
	is.Equal(len(s.Players), 1)
}

func TestNew_StartsDeadlineWhenPositiveHours(t *testing.T) {
	ch := &mockChannel{}
	s := New(ch, "chan1", "gm1", "Spring 1901 Movement", nil, 24, defaultEng(), nil)
	s.mu.Lock()
	timerSet := s.timer != nil
	s.mu.Unlock()
	if !timerSet {
		t.Error("expected deadline timer to be started when DeadlineHours > 0")
	}
	s.CancelDeadline()
}

// writeGameStarted posts a minimal GameStarted event to ch.
func writeGameStarted(ch *mockChannel) {
	_ = events.Write(ch, "chan1", events.TypeGameStarted,
		events.GameStarted{InitialState: json.RawMessage(`{}`)})
}

// ---- Load tests -------------------------------------------------------------

func TestLoad_RebuildsGMIDAndDeadlineHours(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	_ = events.Write(ch, "chan1", events.TypeGameCreated, events.GameCreated{
		Variant: "classical", DeadlineHours: 24, GMUserID: "gm1",
	})
	writeGameStarted(ch)

	s, err := Load(ch, "chan1", nil, makeLoader(defaultEng()))
	is.NoErr(err)
	is.Equal(s.GMID, "gm1")
	is.Equal(s.DeadlineHours, 24)
}

func TestLoad_RebuildsPlayers(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	_ = events.Write(ch, "chan1", events.TypePlayerJoined, events.PlayerJoined{UserID: "u1", Nation: "England"})
	_ = events.Write(ch, "chan1", events.TypePlayerJoined, events.PlayerJoined{UserID: "u2", Nation: "France"})
	writeGameStarted(ch)

	s, err := Load(ch, "chan1", nil, makeLoader(defaultEng()))
	is.NoErr(err)
	is.Equal(s.Players["u1"], "England")
	is.Equal(s.Players["u2"], "France")
}

func TestLoad_RebuildsPhaseFromPhaseResolved(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	writeGameStarted(ch)
	_ = events.Write(ch, "chan1", events.TypePhaseResolved, events.PhaseResolved{
		Phase: "Fall 1901 Movement", StateSnapshot: json.RawMessage(`{}`),
	})

	s, err := Load(ch, "chan1", nil, makeLoader(defaultEng()))
	is.NoErr(err)
	is.Equal(s.Phase, "Fall 1901 Movement")
}

func TestLoad_StagedOrdersFromEventsAfterSnapshot(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	writeGameStarted(ch)
	_ = events.Write(ch, "chan1", events.TypeOrderSubmitted, events.OrderSubmitted{
		Nation: "England", Orders: []string{"A Lon-Nth", "F Eng-Nth"},
	})
	_ = events.Write(ch, "chan1", events.TypeOrderSubmitted, events.OrderSubmitted{
		Nation: "France", Orders: []string{"A Par-Bur"},
	})

	s, err := Load(ch, "chan1", nil, makeLoader(defaultEng()))
	is.NoErr(err)
	is.Equal(len(s.StagedOrders["England"]), 2)
	is.Equal(s.StagedOrders["England"][0], "A Lon-Nth")
	is.Equal(s.StagedOrders["England"][1], "F Eng-Nth")
	is.Equal(len(s.StagedOrders["France"]), 1)
	is.Equal(s.StagedOrders["France"][0], "A Par-Bur")
}

func TestLoad_OrdersBeforeSnapshotNotIncluded(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	// Order submitted before GameStarted — must not appear in StagedOrders.
	_ = events.Write(ch, "chan1", events.TypeOrderSubmitted, events.OrderSubmitted{
		Nation: "England", Orders: []string{"A Lon-Nth"},
	})
	writeGameStarted(ch)
	// Order after snapshot — must appear.
	_ = events.Write(ch, "chan1", events.TypeOrderSubmitted, events.OrderSubmitted{
		Nation: "France", Orders: []string{"A Par-Bur"},
	})

	s, err := Load(ch, "chan1", nil, makeLoader(defaultEng()))
	is.NoErr(err)
	is.Equal(len(s.StagedOrders["England"]), 0)
	is.Equal(len(s.StagedOrders["France"]), 1)
}

func TestLoad_OrdersResetOnNewSnapshot(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	writeGameStarted(ch)
	_ = events.Write(ch, "chan1", events.TypeOrderSubmitted, events.OrderSubmitted{
		Nation: "England", Orders: []string{"A Lon-Nth"},
	})
	// A new PhaseResolved snapshot resets staged orders.
	_ = events.Write(ch, "chan1", events.TypePhaseResolved, events.PhaseResolved{
		Phase: "Fall 1901 Movement", StateSnapshot: json.RawMessage(`{}`),
	})
	// Only this order should appear in StagedOrders.
	_ = events.Write(ch, "chan1", events.TypeOrderSubmitted, events.OrderSubmitted{
		Nation: "France", Orders: []string{"A Par-Bur"},
	})

	s, err := Load(ch, "chan1", nil, makeLoader(defaultEng()))
	is.NoErr(err)
	is.Equal(len(s.StagedOrders["England"]), 0)
	is.Equal(len(s.StagedOrders["France"]), 1)
}

func TestLoad_ReturnsErrorWhenChannelFails(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{histErr: errors.New("channel down")}
	_, err := Load(ch, "chan1", nil, makeLoader(defaultEng()))
	is.Err(err)
}

func TestLoad_ReturnsErrorWhenNoSnapshot(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	_ = events.Write(ch, "chan1", events.TypePlayerJoined, events.PlayerJoined{UserID: "u1", Nation: "England"})
	// No GameStarted or PhaseResolved — Rebuild will fail.
	_, err := Load(ch, "chan1", nil, makeLoader(defaultEng()))
	is.Err(err)
}

func TestLoad_ReturnsErrorWhenEngineFails(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	writeGameStarted(ch)
	badLoader := func(_ []byte) (engine.Engine, error) {
		return nil, errors.New("engine load failed")
	}
	_, err := Load(ch, "chan1", nil, badLoader)
	is.Err(err)
}

func TestLoad_SkipsMalformedGameCreated(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	// Inject a malformed GameCreated payload.
	bad := events.Envelope{Type: events.TypeGameCreated, Payload: json.RawMessage(`"bad"`)}
	data, _ := json.Marshal(bad)
	ch.msgs = append(ch.msgs, string(data))
	writeGameStarted(ch)

	s, err := Load(ch, "chan1", nil, makeLoader(defaultEng()))
	is.NoErr(err)
	is.Equal(s.GMID, "") // malformed payload → GMID stays empty
}

func TestLoad_SkipsMalformedPlayerJoined(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	bad := events.Envelope{Type: events.TypePlayerJoined, Payload: json.RawMessage(`"bad"`)}
	data, _ := json.Marshal(bad)
	ch.msgs = append(ch.msgs, string(data))
	writeGameStarted(ch)

	s, err := Load(ch, "chan1", nil, makeLoader(defaultEng()))
	is.NoErr(err)
	is.Equal(len(s.Players), 0)
}

func TestLoad_SkipsMalformedPhaseResolved(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	writeGameStarted(ch)
	bad := events.Envelope{Type: events.TypePhaseResolved, Payload: json.RawMessage(`"bad"`)}
	data, _ := json.Marshal(bad)
	ch.msgs = append(ch.msgs, string(data))
	// A valid PhaseResolved after the malformed one.
	_ = events.Write(ch, "chan1", events.TypePhaseResolved, events.PhaseResolved{
		Phase: "Fall 1901 Movement", StateSnapshot: json.RawMessage(`{}`),
	})

	s, err := Load(ch, "chan1", nil, makeLoader(defaultEng()))
	is.NoErr(err)
	is.Equal(s.Phase, "Fall 1901 Movement")
}

func TestLoad_SkipsMalformedOrderSubmitted(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	writeGameStarted(ch)
	bad := events.Envelope{Type: events.TypeOrderSubmitted, Payload: json.RawMessage(`"bad"`)}
	data, _ := json.Marshal(bad)
	ch.msgs = append(ch.msgs, string(data))

	s, err := Load(ch, "chan1", nil, makeLoader(defaultEng()))
	is.NoErr(err)
	is.Equal(len(s.StagedOrders), 0)
}

// ---- AdvanceTurn tests ------------------------------------------------------

func TestAdvanceTurn_PostsPhaseResolved(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	s := makeSession(ch, defaultEng(), nil)

	err := s.AdvanceTurn()
	s.CancelDeadline()

	is.NoErr(err)
	// PhaseResolved event + order summary + phase guidance = 3 channel posts.
	is.Equal(ch.msgCount(), 3)

	var env events.Envelope
	is.NoErr(json.Unmarshal([]byte(ch.msgAt(0)), &env))
	is.Equal(env.Type, events.TypePhaseResolved)
}

func TestAdvanceTurn_PhaseResolvedContainsSnapshot(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	eng := &mockEngine{
		resolveResult: engine.ResolutionResult{Phase: "Spring 1901 Movement"},
		dumpData:      []byte(`{"phase":"Spring 1901 Movement"}`),
	}
	s := makeSession(ch, eng, nil)

	is.NoErr(s.AdvanceTurn())
	s.CancelDeadline()

	var env events.Envelope
	is.NoErr(json.Unmarshal([]byte(ch.msgAt(0)), &env))
	var pr events.PhaseResolved
	is.NoErr(json.Unmarshal(env.Payload, &pr))
	is.Equal(pr.Phase, "Spring 1901 Movement")
	is.Equal(string(pr.StateSnapshot), `{"phase":"Spring 1901 Movement"}`)
}

func TestAdvanceTurn_CallsNotifier(t *testing.T) {
	is := is.New(t)
	notifier := &mockNotifier{}
	ch := &mockChannel{}
	s := makeSession(ch, defaultEng(), notifier)

	is.NoErr(s.AdvanceTurn())
	s.CancelDeadline()

	is.Equal(notifier.callCount(), 1)
}

func TestAdvanceTurn_StartsDeadlineTimer(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	s := makeSession(ch, defaultEng(), nil)

	is.NoErr(s.AdvanceTurn())
	s.mu.Lock()
	timerSet := s.timer != nil
	s.mu.Unlock()
	is.Equal(timerSet, true)
	s.CancelDeadline()
}

func TestAdvanceTurn_CancelsExistingTimer(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	s := makeSession(ch, defaultEng(), nil)

	// Plant a long-running timer; AdvanceTurn must cancel it.
	fired := false
	s.mu.Lock()
	s.timer = time.AfterFunc(time.Hour, func() { fired = true })
	s.mu.Unlock()

	is.NoErr(s.AdvanceTurn())
	s.CancelDeadline()

	is.Equal(fired, false)
}

func TestAdvanceTurn_UpdatesPhase(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	eng := &mockEngine{
		resolveResult: engine.ResolutionResult{Phase: "Spring 1901 Movement"},
		dumpData:      []byte(`{}`),
		phaseStr:      "Fall 1901 Movement",
	}
	s := makeSession(ch, eng, nil)

	is.NoErr(s.AdvanceTurn())
	s.CancelDeadline()

	is.Equal(s.Phase, "Fall 1901 Movement")
}

func TestAdvanceTurn_ResetsStagedOrdersAndSubmitted(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	s := makeSession(ch, defaultEng(), nil)
	s.StagedOrders["England"] = []string{"A Lon-Nth"}
	s.Submitted["England"] = true

	is.NoErr(s.AdvanceTurn())
	s.CancelDeadline()

	is.Equal(len(s.StagedOrders), 0)
	is.Equal(len(s.Submitted), 0)
}

func TestAdvanceTurn_DetectsSoloWinner_PostsGameEnded(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	eng := &mockEngine{
		resolveResult: engine.ResolutionResult{Phase: "Fall 1906 Adjustment"},
		dumpData:      []byte(`{}`),
		soloWinner:    "England",
	}
	s := makeSession(ch, eng, nil)

	is.NoErr(s.AdvanceTurn())

	// Three messages: PhaseResolved event, order summary, GameEnded.
	is.Equal(ch.msgCount(), 3)
	var env events.Envelope
	is.NoErr(json.Unmarshal([]byte(ch.msgAt(2)), &env))
	is.Equal(env.Type, events.TypeGameEnded)

	var ge events.GameEnded
	is.NoErr(json.Unmarshal(env.Payload, &ge))
	is.Equal(ge.Result, "solo")
	is.Equal(ge.Winner, "England")
}

func TestAdvanceTurn_SoloWinner_NoTimerStarted(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	eng := &mockEngine{
		resolveResult: engine.ResolutionResult{Phase: "Fall 1906 Adjustment"},
		dumpData:      []byte(`{}`),
		soloWinner:    "France",
	}
	s := makeSession(ch, eng, nil)

	is.NoErr(s.AdvanceTurn())

	// Game is over; no deadline timer should have been started.
	s.mu.Lock()
	timerNil := s.timer == nil
	s.mu.Unlock()
	is.Equal(timerNil, true)
}

func TestAdvanceTurn_ResolveError(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	eng := &mockEngine{resolveErr: errors.New("resolve failed")}
	s := makeSession(ch, eng, nil)

	err := s.AdvanceTurn()
	is.Err(err)
}

func TestAdvanceTurn_DumpError(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	eng := &mockEngine{
		resolveResult: engine.ResolutionResult{Phase: "Spring 1901 Movement"},
		dumpErr:       errors.New("dump failed"),
	}
	s := makeSession(ch, eng, nil)

	err := s.AdvanceTurn()
	is.Err(err)
}

func TestAdvanceTurn_PostError(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{postErr: errors.New("channel unavailable")}
	s := makeSession(ch, defaultEng(), nil)

	err := s.AdvanceTurn()
	is.Err(err)
}

func TestAdvanceTurn_AdvanceError(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	eng := &mockEngine{
		resolveResult: engine.ResolutionResult{Phase: "Spring 1901 Movement"},
		dumpData:      []byte(`{}`),
		advanceErr:    errors.New("advance failed"),
	}
	s := makeSession(ch, eng, nil)

	err := s.AdvanceTurn()
	is.Err(err)
}

// ---- CancelDeadline tests ---------------------------------------------------

func TestCancelDeadline_StopsTimer(t *testing.T) {
	s := &Session{}
	s.mu.Lock()
	s.timer = time.AfterFunc(time.Hour, func() {})
	s.mu.Unlock()
	s.CancelDeadline()
	s.mu.Lock()
	timerNil := s.timer == nil
	s.mu.Unlock()
	if !timerNil {
		t.Error("expected timer to be nil after CancelDeadline")
	}
}

func TestCancelDeadline_NilTimerIsNoOp(t *testing.T) {
	// Must not panic when timer is nil.
	s := &Session{}
	s.CancelDeadline()
}

// ---- Deadline timer fires AdvanceTurn ---------------------------------------

func TestDeadlineTimer_FiresAdvanceTurn(t *testing.T) {
	ch := &mockChannel{}
	s := makeSession(ch, defaultEng(), nil)
	// DeadlineHours = 0 so the startDeadline inside AdvanceTurn is a no-op.
	s.DeadlineHours = 0

	var wg sync.WaitGroup
	wg.Add(1)

	// Plant a 1 ms timer pointing at onDeadline — the real callback used in
	// production. Hold the mutex while writing s.timer to satisfy the race
	// detector (the callback goroutine reads s.timer via CancelDeadline).
	s.mu.Lock()
	s.timer = time.AfterFunc(time.Millisecond, func() {
		defer wg.Done()
		s.onDeadline()
	})
	s.mu.Unlock()

	// Wait for the callback to complete fully before asserting.
	wg.Wait()

	// At least the PhaseResolved event must have been posted.
	if ch.msgCount() < 1 {
		t.Error("expected at least 1 message after deadline timer fired onDeadline")
	}
}

// ---- startDeadline / onDeadline tests ---------------------------------------

func TestStartDeadline_ZeroHoursDoesNotStartTimer(t *testing.T) {
	s := &Session{DeadlineHours: 0}
	s.startDeadline()
	if s.timer != nil {
		t.Error("expected no timer when DeadlineHours is 0")
	}
}

func TestStartDeadline_PositiveHoursSetsTimer(t *testing.T) {
	s := &Session{DeadlineHours: 24}
	s.startDeadline()
	s.mu.Lock()
	timerSet := s.timer != nil
	s.mu.Unlock()
	if !timerSet {
		t.Error("expected timer to be set when DeadlineHours > 0")
	}
	s.CancelDeadline()
}

// TestOnDeadline_CallsAdvanceTurn verifies that onDeadline (the timer callback)
// invokes AdvanceTurn. This covers the callback path without needing real timers.
func TestOnDeadline_CallsAdvanceTurn(t *testing.T) {
	is := is.New(t)
	ch := &mockChannel{}
	s := makeSession(ch, defaultEng(), nil)

	s.onDeadline()
	s.CancelDeadline()

	// PhaseResolved event + order summary + phase guidance = 3 channel posts.
	is.Equal(ch.msgCount(), 3)
}

// ---- RestartDeadline tests --------------------------------------------------

func TestRestartDeadline_StartsTimerWithRemainingTime(t *testing.T) {
	s := &Session{DeadlineHours: 24}
	s.mu.Lock()
	s.deadlineAt = time.Now().Add(time.Hour) // 1 hour remaining
	s.mu.Unlock()

	s.RestartDeadline()

	s.mu.Lock()
	timerSet := s.timer != nil
	s.mu.Unlock()
	if !timerSet {
		t.Error("expected timer to be set when deadlineAt is in the future")
	}
	s.CancelDeadline()
}

func TestRestartDeadline_UsesFullDeadlineWhenDeadlineAtInPast(t *testing.T) {
	s := &Session{DeadlineHours: 24}
	s.mu.Lock()
	s.deadlineAt = time.Now().Add(-time.Hour) // already in the past
	s.mu.Unlock()

	s.RestartDeadline()

	s.mu.Lock()
	timerSet := s.timer != nil
	s.mu.Unlock()
	if !timerSet {
		t.Error("expected timer to be set using full DeadlineHours fallback")
	}
	s.CancelDeadline()
}

func TestRestartDeadline_NoOpWhenDeadlineHoursZeroAndDeadlineAtZero(t *testing.T) {
	s := &Session{DeadlineHours: 0}
	s.RestartDeadline()
	s.mu.Lock()
	timerNil := s.timer == nil
	s.mu.Unlock()
	if !timerNil {
		t.Error("expected no timer when DeadlineHours=0 and deadlineAt is zero")
	}
}

func TestRestartDeadline_CancelsExistingTimer(t *testing.T) {
	is := is.New(t)
	s := &Session{DeadlineHours: 24}
	s.mu.Lock()
	s.deadlineAt = time.Now().Add(time.Hour)
	s.mu.Unlock()

	fired := false
	s.mu.Lock()
	s.timer = time.AfterFunc(time.Hour, func() { fired = true })
	s.mu.Unlock()

	s.RestartDeadline()
	s.CancelDeadline()

	is.Equal(fired, false)
}

func TestRestartDeadline_NoTimerCancelWhenNilTimer(t *testing.T) {
	// Should not panic when called with no existing timer.
	s := &Session{DeadlineHours: 24}
	s.mu.Lock()
	s.deadlineAt = time.Now().Add(time.Hour)
	s.mu.Unlock()
	// s.timer is nil here
	s.RestartDeadline()
	s.CancelDeadline()
}

// ---- ExtendDeadline tests ---------------------------------------------------

func TestExtendDeadline_ExtendsByDurationFromExistingDeadline(t *testing.T) {
	s := &Session{DeadlineHours: 24}
	s.mu.Lock()
	s.deadlineAt = time.Now().Add(time.Hour)
	s.mu.Unlock()

	s.ExtendDeadline(30 * time.Minute)

	s.mu.Lock()
	timerSet := s.timer != nil
	s.mu.Unlock()
	if !timerSet {
		t.Error("expected timer to be set after extending deadline")
	}
	s.CancelDeadline()
}

func TestExtendDeadline_FallsBackToDeadlineHoursWhenDeadlineAtZero(t *testing.T) {
	s := &Session{DeadlineHours: 24}
	// deadlineAt is zero value

	s.ExtendDeadline(time.Hour)

	s.mu.Lock()
	timerSet := s.timer != nil
	s.mu.Unlock()
	if !timerSet {
		t.Error("expected timer to be set when deadlineAt is zero but DeadlineHours > 0")
	}
	s.CancelDeadline()
}

func TestExtendDeadline_NoOpWhenDeadlineHoursZeroAndDeadlineAtZero(t *testing.T) {
	s := &Session{DeadlineHours: 0}
	// deadlineAt is zero
	s.ExtendDeadline(time.Hour)
	s.mu.Lock()
	timerNil := s.timer == nil
	s.mu.Unlock()
	if !timerNil {
		t.Error("expected no timer when DeadlineHours=0 and deadlineAt is zero")
	}
}

func TestExtendDeadline_NoTimerWhenExtendedDeadlineInPast(t *testing.T) {
	s := &Session{DeadlineHours: 0}
	s.mu.Lock()
	s.deadlineAt = time.Now().Add(-2 * time.Hour)
	s.mu.Unlock()

	s.ExtendDeadline(-time.Hour) // still in the past after extension

	s.mu.Lock()
	timerNil := s.timer == nil
	s.mu.Unlock()
	if !timerNil {
		t.Error("expected no timer when extended deadline is still in the past")
	}
}

func TestExtendDeadline_CancelsExistingTimer(t *testing.T) {
	is := is.New(t)
	s := &Session{DeadlineHours: 24}
	s.mu.Lock()
	s.deadlineAt = time.Now().Add(time.Hour)
	s.mu.Unlock()

	fired := false
	s.mu.Lock()
	s.timer = time.AfterFunc(time.Hour, func() { fired = true })
	s.mu.Unlock()

	s.ExtendDeadline(30 * time.Minute)
	s.CancelDeadline()

	is.Equal(fired, false)
}

// ---- lifecycle helper tests -------------------------------------------------

func TestNationList_ReturnsUniqueSorted(t *testing.T) {
	is := is.New(t)
	s := &Session{
		Players: map[string]string{
			"u1": "England",
			"u2": "France",
			"u3": "England", // duplicate nation
		},
	}
	nations := s.nationList()
	is.Equal(nations, []string{"England", "France"})
}

func TestNationList_Empty(t *testing.T) {
	is := is.New(t)
	s := &Session{Players: map[string]string{}}
	nations := s.nationList()
	is.Equal(len(nations), 0)
}

func TestFormatOrderSummary_AllNMR(t *testing.T) {
	is := is.New(t)
	result := engine.ResolutionResult{
		Phase:  "Spring 1901 Movement",
		Season: "Spring",
		Year:   1901,
		Orders: []engine.OrderResult{
			{Nation: "England", Order: "A lon H", IsNMR: true},
		},
	}
	out := formatOrderSummary(result, []string{"England"})
	is.Equal(out != "", true)
	// Should contain "no orders submitted" for an all-NMR nation.
	if !containsStr(out, "no orders submitted") {
		t.Fatalf("expected 'no orders submitted' in %q", out)
	}
}

func TestFormatOrderSummary_MixedOrders(t *testing.T) {
	is := is.New(t)
	result := engine.ResolutionResult{
		Phase:  "Spring 1901 Movement",
		Season: "Spring",
		Year:   1901,
		Orders: []engine.OrderResult{
			{Nation: "England", Order: "A lon-bel", IsNMR: false, Outcome: "success"},
			{Nation: "England", Order: "F nth H", IsNMR: true, Outcome: ""},
		},
	}
	out := formatOrderSummary(result, []string{"England"})
	is.Equal(out != "", true)
	// Should show the real order with outcome and the NMR order as "(auto)".
	if !containsStr(out, "success") {
		t.Fatalf("expected 'success' in %q", out)
	}
	if !containsStr(out, "(auto)") {
		t.Fatalf("expected '(auto)' in %q", out)
	}
}

func TestFormatOrderSummary_NationWithNoOrders(t *testing.T) {
	is := is.New(t)
	result := engine.ResolutionResult{
		Phase:  "Spring 1901 Movement",
		Season: "Spring",
		Year:   1901,
		Orders: []engine.OrderResult{},
	}
	out := formatOrderSummary(result, []string{"England"})
	// England has no orders at all — should be skipped.
	is.Equal(out != "", true)
}

func TestFormatPhaseGuidance_MovementPhase(t *testing.T) {
	is := is.New(t)
	eng := &mockEngine{phaseStr: "Spring 1901 Movement"}
	s := &Session{Phase: "Spring 1901 Movement", Eng: eng}
	out := s.formatPhaseGuidance([]string{"England"})
	if !containsStr(out, "submit orders") {
		t.Fatalf("expected 'submit orders' in %q", out)
	}
	is.Equal(out != "", true)
}

func TestFormatPhaseGuidance_RetreatPhase(t *testing.T) {
	is := is.New(t)
	eng := &mockEngine{phaseStr: "Spring 1901 Retreat"}
	s := &Session{Phase: "Spring 1901 Retreat", Eng: eng}
	out := s.formatPhaseGuidance([]string{"England"})
	if !containsStr(out, "must retreat") {
		t.Fatalf("expected 'must retreat' in %q", out)
	}
	is.Equal(out != "", true)
}

func TestFormatPhaseGuidance_AdjustmentPhase(t *testing.T) {
	is := is.New(t)
	eng := &mockEngine{phaseStr: "Fall 1901 Adjustment"}
	s := &Session{Phase: "Fall 1901 Adjustment", Eng: eng}
	out := s.formatPhaseGuidance([]string{"England"})
	is.Equal(out != "", true)
}

func TestFormatPhaseGuidance_UnknownPhase(t *testing.T) {
	is := is.New(t)
	eng := &mockEngine{phaseStr: "Some 1901 Unknown"}
	s := &Session{Phase: "Some 1901 Unknown", Eng: eng}
	out := s.formatPhaseGuidance([]string{"England"})
	if !containsStr(out, "new phase begun") {
		t.Fatalf("expected 'new phase begun' in %q", out)
	}
	is.Equal(out != "", true)
}

func TestFormatMovementGuidance_NationWithUnits(t *testing.T) {
	is := is.New(t)
	eng := &mockEngineWithUnits{
		units: map[string]engine.UnitInfo{
			"lon": {Nation: "England"},
			"lvp": {Nation: "England"},
		},
	}
	out := formatMovementGuidance("Spring 1901 Movement", eng, []string{"England"})
	if !containsStr(out, "England") {
		t.Fatalf("expected 'England' in %q", out)
	}
	if !containsStr(out, "lon") {
		t.Fatalf("expected 'lon' in %q", out)
	}
	is.Equal(out != "", true)
}

func TestFormatMovementGuidance_NationWithNoUnits(t *testing.T) {
	is := is.New(t)
	eng := &mockEngineWithUnits{units: map[string]engine.UnitInfo{}}
	out := formatMovementGuidance("Spring 1901 Movement", eng, []string{"England"})
	if !containsStr(out, "nothing to do") {
		t.Fatalf("expected 'nothing to do' in %q", out)
	}
	is.Equal(out != "", true)
}

func TestFormatRetreatGuidance_NationWithRetreats(t *testing.T) {
	is := is.New(t)
	eng := &mockEngineWithRetreats{
		dislodgeds:    map[string]string{"vie": "Austria"},
		validRetreats: map[string][]string{"vie": {"bud", "tri"}},
	}
	out := formatRetreatGuidance("Spring 1901 Retreat", eng, []string{"Austria"})
	if !containsStr(out, "Austria") {
		t.Fatalf("expected 'Austria' in %q", out)
	}
	if !containsStr(out, "retreat to") {
		t.Fatalf("expected 'retreat to' in %q", out)
	}
	is.Equal(out != "", true)
}

func TestFormatRetreatGuidance_UnitMustDisband(t *testing.T) {
	is := is.New(t)
	eng := &mockEngineWithRetreats{
		dislodgeds:    map[string]string{"vie": "Austria"},
		validRetreats: map[string][]string{"vie": {}},
	}
	out := formatRetreatGuidance("Spring 1901 Retreat", eng, []string{"Austria"})
	if !containsStr(out, "must disband") {
		t.Fatalf("expected 'must disband' in %q", out)
	}
	is.Equal(out != "", true)
}

func TestFormatRetreatGuidance_NationNothingToDo(t *testing.T) {
	is := is.New(t)
	eng := &mockEngineWithRetreats{
		dislodgeds:    map[string]string{},
		validRetreats: map[string][]string{},
	}
	out := formatRetreatGuidance("Spring 1901 Retreat", eng, []string{"England"})
	if !containsStr(out, "nothing to do") {
		t.Fatalf("expected 'nothing to do' in %q", out)
	}
	is.Equal(out != "", true)
}

func TestFormatAdjustmentGuidance_BuildCase(t *testing.T) {
	is := is.New(t)
	eng := &mockEngineWithBuild{
		opts: map[string]engine.BuildOption{
			"England": {Delta: 1, AvailableHomes: []string{"edi"}},
		},
	}
	out := formatAdjustmentGuidance("Fall 1901 Adjustment", eng, []string{"England"})
	if !containsStr(out, "build") {
		t.Fatalf("expected 'build' in %q", out)
	}
	is.Equal(out != "", true)
}

func TestFormatAdjustmentGuidance_DisbandCase(t *testing.T) {
	is := is.New(t)
	eng := &mockEngineWithBuild{
		opts: map[string]engine.BuildOption{
			"France": {Delta: -1},
		},
	}
	out := formatAdjustmentGuidance("Fall 1901 Adjustment", eng, []string{"France"})
	if !containsStr(out, "disband") {
		t.Fatalf("expected 'disband' in %q", out)
	}
	is.Equal(out != "", true)
}

func TestFormatAdjustmentGuidance_NothingToDo(t *testing.T) {
	is := is.New(t)
	eng := &mockEngineWithBuild{
		opts: map[string]engine.BuildOption{},
	}
	out := formatAdjustmentGuidance("Fall 1901 Adjustment", eng, []string{"England"})
	if !containsStr(out, "nothing to do") {
		t.Fatalf("expected 'nothing to do' in %q", out)
	}
	is.Equal(out != "", true)
}

// ---- local mock helpers for lifecycle tests ---------------------------------

func containsStr(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstr(s, substr))
}

func containsSubstr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

type mockEngineWithUnits struct {
	mockEngine
	units map[string]engine.UnitInfo
}

func (e *mockEngineWithUnits) Units() map[string]engine.UnitInfo { return e.units }

type mockEngineWithRetreats struct {
	mockEngine
	dislodgeds    map[string]string
	validRetreats map[string][]string
}

func (e *mockEngineWithRetreats) Dislodgeds() map[string]string   { return e.dislodgeds }
func (e *mockEngineWithRetreats) ValidRetreats() map[string][]string { return e.validRetreats }

type mockEngineWithBuild struct {
	mockEngine
	opts map[string]engine.BuildOption
}

func (e *mockEngineWithBuild) BuildOptions() map[string]engine.BuildOption { return e.opts }
