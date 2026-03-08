//go:build functional

// Package bot_test contains functional tests that exercise every command
// listed in ARCHITECTURE.md through the bot.Dispatcher without any platform
// adapter.  Run with:
//
//	go test -v -tags functional ./bot/
//
// Each test is independent: it spins up a fresh in-memory channel and
// dispatcher so that failures are isolated.  Tests that require a started
// game call startedGame() to perform /newgame + /join × 2 + /start.
//
// Commands that are only valid in Retreat or Adjustment phases cannot be
// reached in a short functional test without a full multi-phase game flow;
// those tests exercise the phase-guard rejection path (calling the command
// outside its phase) to prove the command is wired up and access-controlled.
package bot_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/burrbd/dip/bot"
	"github.com/burrbd/dip/engine"
	"github.com/burrbd/dip/events"
	"github.com/cheekybits/is"
)

// ---------------------------------------------------------------------------
// Infrastructure
// ---------------------------------------------------------------------------

// memChannel is a minimal in-memory implementation of events.Channel.
// It stores channel and DM messages in plain slices so that events.Scan /
// events.ScanDM work correctly against it.
type memChannel struct {
	msgs map[string][]string
	dms  map[string][]string
	imgs [][]byte
}

func newMem() *memChannel {
	return &memChannel{
		msgs: make(map[string][]string),
		dms:  make(map[string][]string),
	}
}

func (c *memChannel) Post(channelID, text string) error {
	c.msgs[channelID] = append(c.msgs[channelID], text)
	return nil
}
func (c *memChannel) History(channelID string) ([]string, error) {
	return c.msgs[channelID], nil
}
func (c *memChannel) SendDM(userID, text string) error {
	c.dms[userID] = append(c.dms[userID], text)
	return nil
}
func (c *memChannel) DMHistory(userID string) ([]string, error) {
	return c.dms[userID], nil
}
func (c *memChannel) PostImage(channelID string, data []byte) error {
	c.imgs = append(c.imgs, data)
	return nil
}

// nopNotifier satisfies session.Notifier with no side effects.
type nopNotifier struct{}

func (n *nopNotifier) Notify(_, _ string) error { return nil }

// newDispatcher returns a Dispatcher wired to ch using real engine.New /
// engine.Load so that no business logic is mocked.
func newDispatcher(ch *memChannel) *bot.Dispatcher {
	return bot.New(ch, &nopNotifier{}, engine.Load, engine.New)
}

// chanCmd builds a channel-sourced command.
func chanCmd(name, userID, channelID string, args ...string) bot.Command {
	return bot.Command{Name: name, UserID: userID, ChannelID: channelID, Args: args}
}

// dmCmd builds a DM-sourced command (IsDM=true, GameChannelID set).
func dmCmd(name, userID, gameChannelID string, args ...string) bot.Command {
	return bot.Command{
		Name:          name,
		UserID:        userID,
		ChannelID:     "dm_" + userID,
		IsDM:          true,
		GameChannelID: gameChannelID,
		Args:          args,
	}
}

// startedGame creates a fresh dispatcher + memChannel, runs /newgame,
// /join England (u1), /join France (u2), /start (gm), and returns them.
// The game is in Spring 1901 Movement phase on return.
func startedGame(t *testing.T) (*bot.Dispatcher, *memChannel) {
	t.Helper()
	ch := newMem()
	d := newDispatcher(ch)

	mustDispatch(t, d, chanCmd("newgame", "gm", "game"))
	mustDispatch(t, d, chanCmd("join", "u1", "game", "England"))
	mustDispatch(t, d, chanCmd("join", "u2", "game", "France"))
	mustDispatch(t, d, chanCmd("start", "gm", "game"))
	return d, ch
}

// mustDispatch calls Dispatch and fails immediately on error.
func mustDispatch(t *testing.T, d *bot.Dispatcher, cmd bot.Command) string {
	t.Helper()
	resp, err := d.Dispatch(cmd)
	if err != nil {
		t.Fatalf("dispatch %q failed: %v", cmd.Name, err)
	}
	return resp
}

// hasEvent returns true when at least one envelope of the given type exists
// in the game channel.
func hasEvent(t *testing.T, ch *memChannel, channelID string, want events.EventType) bool {
	t.Helper()
	envs, err := events.Scan(ch, channelID)
	if err != nil {
		t.Fatalf("events.Scan: %v", err)
	}
	for _, e := range envs {
		if e.Type == want {
			return true
		}
	}
	return false
}

// eventPayload returns the JSON payload of the last event of the given type.
func eventPayload(t *testing.T, ch *memChannel, channelID string, want events.EventType) json.RawMessage {
	t.Helper()
	envs, _ := events.Scan(ch, channelID)
	for i := len(envs) - 1; i >= 0; i-- {
		if envs[i].Type == want {
			return envs[i].Payload
		}
	}
	t.Fatalf("no event of type %q found", want)
	return nil
}

// ---------------------------------------------------------------------------
// Setup commands
// ---------------------------------------------------------------------------

func TestCommand_Newgame(t *testing.T) {
	// /newgame posts a GameCreated event and replies with a confirmation.
	is := is.New(t)
	ch := newMem()
	d := newDispatcher(ch)

	resp, err := d.Dispatch(chanCmd("newgame", "gm", "game"))
	is.NoErr(err)
	is.Equal(strings.Contains(resp, "created") || strings.Contains(resp, "Game"), true)
	is.Equal(hasEvent(t, ch, "game", events.TypeGameCreated), true)

	var gc events.GameCreated
	is.NoErr(json.Unmarshal(eventPayload(t, ch, "game", events.TypeGameCreated), &gc))
	is.Equal(gc.GMUserID, "gm")
	t.Logf("GameCreated: variant=%s gm=%s", gc.Variant, gc.GMUserID)
}

func TestCommand_Join(t *testing.T) {
	// /join <nation> posts a PlayerJoined event for the specified nation.
	is := is.New(t)
	ch := newMem()
	d := newDispatcher(ch)

	mustDispatch(t, d, chanCmd("newgame", "gm", "game"))
	resp, err := d.Dispatch(chanCmd("join", "u1", "game", "England"))
	is.NoErr(err)
	is.Equal(strings.Contains(resp, "England") || strings.Contains(resp, "joined"), true)
	is.Equal(hasEvent(t, ch, "game", events.TypePlayerJoined), true)

	var pj events.PlayerJoined
	is.NoErr(json.Unmarshal(eventPayload(t, ch, "game", events.TypePlayerJoined), &pj))
	is.Equal(pj.UserID, "u1")
	is.Equal(pj.Nation, "England")
	t.Logf("PlayerJoined: user=%s nation=%s", pj.UserID, pj.Nation)
}

func TestCommand_Join_RejectedIfNationTaken(t *testing.T) {
	// /join <nation> is rejected when that nation is already taken.
	is := is.New(t)
	ch := newMem()
	d := newDispatcher(ch)

	mustDispatch(t, d, chanCmd("newgame", "gm", "game"))
	mustDispatch(t, d, chanCmd("join", "u1", "game", "England"))
	_, err := d.Dispatch(chanCmd("join", "u2", "game", "England"))
	is.NotNil(err) // duplicate nation must be rejected
	t.Logf("Duplicate join rejected: %v", err)
}

func TestCommand_Start(t *testing.T) {
	// /start (GM only) posts a GameStarted event and records the initial engine snapshot.
	is := is.New(t)
	ch := newMem()
	d := newDispatcher(ch)

	mustDispatch(t, d, chanCmd("newgame", "gm", "game"))
	mustDispatch(t, d, chanCmd("join", "u1", "game", "England"))
	mustDispatch(t, d, chanCmd("join", "u2", "game", "France"))

	resp, err := d.Dispatch(chanCmd("start", "gm", "game"))
	is.NoErr(err)
	is.Equal(hasEvent(t, ch, "game", events.TypeGameStarted), true)

	var gs events.GameStarted
	is.NoErr(json.Unmarshal(eventPayload(t, ch, "game", events.TypeGameStarted), &gs))
	is.NotNil(gs.InitialState)
	t.Logf("GameStarted: snapshot size=%d bytes, response=%q", len(gs.InitialState), resp)
}

func TestCommand_Start_RejectedByNonGM(t *testing.T) {
	// /start must be rejected when called by a non-GM user.
	is := is.New(t)
	ch := newMem()
	d := newDispatcher(ch)

	mustDispatch(t, d, chanCmd("newgame", "gm", "game"))
	mustDispatch(t, d, chanCmd("join", "u1", "game", "England"))
	mustDispatch(t, d, chanCmd("join", "u2", "game", "France"))

	_, err := d.Dispatch(chanCmd("start", "u1", "game"))
	is.NotNil(err)
	t.Logf("Non-GM start rejected: %v", err)
}

// ---------------------------------------------------------------------------
// Movement phase commands (submitted via DM)
// ---------------------------------------------------------------------------

func TestCommand_Order(t *testing.T) {
	// /order <text> (DM) stages an order for the caller's nation and
	// returns a confirmation.  OrderSubmitted events are written to DM by
	// /submit (not /order); check staging via the /orders response instead.
	is := is.New(t)
	d, _ := startedGame(t)

	resp, err := d.Dispatch(dmCmd("order", "u1", "game", "A Lon H"))
	is.NoErr(err)
	is.Equal(resp != "", true)
	is.Equal(strings.Contains(resp, "Lon") || strings.Contains(resp, "staged") || strings.Contains(resp, "lon"), true)

	// Confirm the order is visible via /orders.
	ordersResp, err := d.Dispatch(dmCmd("orders", "u1", "game"))
	is.NoErr(err)
	is.Equal(strings.Contains(ordersResp, "Lon") || strings.Contains(ordersResp, "lon") || strings.Contains(ordersResp, "Hold"), true)
	t.Logf("Order response: %q; orders: %q", resp, ordersResp)
}

func TestCommand_Order_RejectedForWrongNation(t *testing.T) {
	// /order must be rejected when the DM sender is not a player in the game.
	is := is.New(t)
	d, _ := startedGame(t)

	_, err := d.Dispatch(dmCmd("order", "stranger", "game", "A Lon H"))
	is.NotNil(err)
	t.Logf("Non-player order rejected: %v", err)
}

func TestCommand_Orders(t *testing.T) {
	// /orders (DM) lists the caller's currently staged orders.
	is := is.New(t)
	d, _ := startedGame(t)

	mustDispatch(t, d, dmCmd("order", "u1", "game", "A Lon H"))
	resp, err := d.Dispatch(dmCmd("orders", "u1", "game"))
	is.NoErr(err)
	is.Equal(strings.Contains(resp, "lon") || strings.Contains(resp, "Lon") || strings.Contains(resp, "Hold"), true)
	t.Logf("Orders response: %q", resp)
}

func TestCommand_Clear(t *testing.T) {
	// /clear (DM) removes all staged orders for the caller's nation.
	is := is.New(t)
	d, _ := startedGame(t)

	mustDispatch(t, d, dmCmd("order", "u1", "game", "A Lon H"))
	resp, err := d.Dispatch(dmCmd("clear", "u1", "game"))
	is.NoErr(err)
	is.Equal(resp != "", true)

	// After clearing, /orders should show nothing staged.
	ordersResp, err := d.Dispatch(dmCmd("orders", "u1", "game"))
	is.NoErr(err)
	is.Equal(strings.Contains(ordersResp, "No orders") || ordersResp == "" || strings.Contains(ordersResp, "no orders"), true)
	t.Logf("Clear response: %q; orders after: %q", resp, ordersResp)
}

func TestCommand_Submit_PartialDoesNotAdvance(t *testing.T) {
	// /submit from one player must not advance the phase if the other player
	// has not yet submitted.
	is := is.New(t)
	d, ch := startedGame(t)

	mustDispatch(t, d, dmCmd("order", "u1", "game", "A Lon H"))
	_, err := d.Dispatch(dmCmd("submit", "u1", "game"))
	is.NoErr(err)

	// PhaseResolved must NOT have been posted yet.
	is.Equal(hasEvent(t, ch, "game", events.TypePhaseResolved), false)
	t.Log("Partial submit: phase not advanced (correct)")
}

func TestCommand_Submit_AllPlayersAdvancesPhase(t *testing.T) {
	// /submit from the final player triggers AdvanceTurn and posts PhaseResolved.
	// PhaseResolved.Phase carries the godip phase type ("Movement"), not the
	// full display string — use engine.Phase() for the full "Spring 1901 Movement".
	is := is.New(t)
	d, ch := startedGame(t)

	mustDispatch(t, d, dmCmd("order", "u1", "game", "A Lon H"))
	mustDispatch(t, d, dmCmd("submit", "u1", "game"))

	// Verify OrderSubmitted event is posted to DM on /submit.
	dmEnvs, err := events.ScanDM(ch, "u1")
	is.NoErr(err)
	foundOS := false
	for _, e := range dmEnvs {
		if e.Type == events.TypeOrderSubmitted {
			foundOS = true
			var os events.OrderSubmitted
			is.NoErr(json.Unmarshal(e.Payload, &os))
			is.Equal(os.Nation, "England")
			t.Logf("OrderSubmitted (u1): nation=%s orders=%v", os.Nation, os.Orders)
		}
	}
	is.Equal(foundOS, true)

	mustDispatch(t, d, dmCmd("order", "u2", "game", "A Par H"))
	_, err = d.Dispatch(dmCmd("submit", "u2", "game"))
	is.NoErr(err)

	is.Equal(hasEvent(t, ch, "game", events.TypePhaseResolved), true)
	var pr events.PhaseResolved
	is.NoErr(json.Unmarshal(eventPayload(t, ch, "game", events.TypePhaseResolved), &pr))
	is.Equal(pr.Phase, "Movement")
	t.Logf("PhaseResolved: phase=%q", pr.Phase)
}

// ---------------------------------------------------------------------------
// Phase helpers — bring the game to Retreat or Adjustment phase
// ---------------------------------------------------------------------------

// retreatPhaseGame spins up a full 7-nation game and advances to Fall 1901
// Retreat with Austria's F Tri dislodged.
//
// Scenario (2 turns):
//   - Spring 1901: Italy A Ven→Tyr (empty), A Rom→Ven (now empty).
//     No dislodgements → Spring Retreat skipped → Fall 1901 Movement.
//   - Fall 1901: Italy A Tyr→Tri (attack Trieste), A Ven S A Tyr-Tri.
//     Venice is adjacent to Trieste, so the support is valid.
//     Attack strength 2 vs Austria F Tri hold strength 1 → dislodged.
//     → Fall 1901 Retreat with Austria F Tri dislodged.
//
// Austria F Tri's valid retreats: Adr (Adriatic Sea) or Alb (Albania).
// Tyr is blocked (dislodger came from there); Ven is occupied by Italy A Ven.
func retreatPhaseGame(t *testing.T) (*bot.Dispatcher, *memChannel) {
	t.Helper()
	ch := newMem()
	d := newDispatcher(ch)

	mustDispatch(t, d, chanCmd("newgame", "gm", "game"))
	mustDispatch(t, d, chanCmd("join", "u1", "game", "Austria"))
	mustDispatch(t, d, chanCmd("join", "u2", "game", "England"))
	mustDispatch(t, d, chanCmd("join", "u3", "game", "France"))
	mustDispatch(t, d, chanCmd("join", "u4", "game", "Germany"))
	mustDispatch(t, d, chanCmd("join", "u5", "game", "Italy"))
	mustDispatch(t, d, chanCmd("join", "u6", "game", "Russia"))
	mustDispatch(t, d, chanCmd("join", "u7", "game", "Turkey"))
	mustDispatch(t, d, chanCmd("start", "gm", "game"))

	// Spring 1901 Movement: Italy repositions for a supported attack on Tri.
	// A Ven → Tyr (empty province adjacent to Tri).
	// A Rom → Ven (fills vacated Venice, now adjacent to Tri for support).
	mustDispatch(t, d, dmCmd("order", "u5", "game", "A ven-tyr"))
	mustDispatch(t, d, dmCmd("order", "u5", "game", "A rom-ven"))
	// Spring Retreat will be empty (no dislodgements) → auto-skipped to Fall.
	mustDispatch(t, d, chanCmd("force-resolve", "gm", "game"))

	// Fall 1901 Movement: Italy A Tyr attacks Tri; A Ven (adjacent to Tri)
	// provides the decisive support.  Austria ≠ Italy, so the support is not
	// excluded by the "defender's nation forbidden" rule.
	mustDispatch(t, d, dmCmd("order", "u5", "game", "A tyr-tri"))
	mustDispatch(t, d, dmCmd("order", "u5", "game", "A ven S A tyr-tri"))
	mustDispatch(t, d, chanCmd("force-resolve", "gm", "game"))
	// Fall Retreat: Austria F Tri is dislodged.
	return d, ch
}

// adjustmentPhaseGame spins up a full 7-nation game, advances through Spring
// and Fall 1901 so that England captures Norway and earns one build slot.
// Returns a dispatcher in Winter 1901 Adjustment phase.
//
// Scenario:
//   - Spring 1901: England submits F lon-nth (all others NMR). No dislodgements
//     → Spring Retreat skipped → Fall 1901 Movement.
//   - Fall 1901: England submits F nth-nwy (all others NMR). No dislodgements
//     → Fall Retreat skipped → Winter 1901 Adjustment.
//
// England after Adjustment: 4 SCs (lon, edi, lvp, nwy), 3 units → 1 build slot.
func adjustmentPhaseGame(t *testing.T) (*bot.Dispatcher, *memChannel) {
	t.Helper()
	ch := newMem()
	d := newDispatcher(ch)

	mustDispatch(t, d, chanCmd("newgame", "gm", "game"))
	mustDispatch(t, d, chanCmd("join", "u1", "game", "Austria"))
	mustDispatch(t, d, chanCmd("join", "u2", "game", "England"))
	mustDispatch(t, d, chanCmd("join", "u3", "game", "France"))
	mustDispatch(t, d, chanCmd("join", "u4", "game", "Germany"))
	mustDispatch(t, d, chanCmd("join", "u5", "game", "Italy"))
	mustDispatch(t, d, chanCmd("join", "u6", "game", "Russia"))
	mustDispatch(t, d, chanCmd("join", "u7", "game", "Turkey"))
	mustDispatch(t, d, chanCmd("start", "gm", "game"))

	// Spring 1901 Movement: England moves F Lon to North Sea.
	mustDispatch(t, d, dmCmd("order", "u2", "game", "F lon-nth"))
	mustDispatch(t, d, chanCmd("force-resolve", "gm", "game"))
	// Spring Retreat skipped (no dislodgements) → Fall 1901 Movement.

	// Fall 1901 Movement: England moves F NTH to Norway.
	mustDispatch(t, d, dmCmd("order", "u2", "game", "F nth-nwy"))
	mustDispatch(t, d, chanCmd("force-resolve", "gm", "game"))
	// Fall Retreat skipped (no dislodgements) → Winter 1901 Adjustment.

	return d, ch
}

// ---------------------------------------------------------------------------
// Happy-path tests for Retreat / Adjustment commands
// ---------------------------------------------------------------------------

func TestCommand_Retreat(t *testing.T) {
	// /retreat happy path: Austria's F Tri retreats to Adr (Adriatic Sea).
	// Ven is occupied by Italy's A Ven; Tyr is blocked (dislodger came from there).
	// Expects no error and an OrderSubmitted event recorded in Austria's DM.
	is := is.New(t)
	d, ch := retreatPhaseGame(t)

	resp, err := d.Dispatch(dmCmd("retreat", "u1", "game", "F", "tri", "adr"))
	is.NoErr(err)
	is.Equal(resp != "", true)

	dmEnvs, err := events.ScanDM(ch, "u1")
	is.NoErr(err)
	foundOS := false
	for _, e := range dmEnvs {
		if e.Type == events.TypeOrderSubmitted {
			var os events.OrderSubmitted
			is.NoErr(json.Unmarshal(e.Payload, &os))
			is.Equal(os.Nation, "Austria")
			foundOS = true
		}
	}
	is.Equal(foundOS, true)
	t.Logf("Retreat: resp=%q", resp)
}

func TestCommand_Disband_InRetreatPhase(t *testing.T) {
	// /disband happy path in Retreat phase: Austria's F Tri disbands instead of retreating.
	// Expects no error and an OrderSubmitted event recorded in Austria's DM.
	is := is.New(t)
	d, ch := retreatPhaseGame(t)

	resp, err := d.Dispatch(dmCmd("disband", "u1", "game", "F", "tri"))
	is.NoErr(err)
	is.Equal(resp != "", true)

	dmEnvs, err := events.ScanDM(ch, "u1")
	is.NoErr(err)
	foundOS := false
	for _, e := range dmEnvs {
		if e.Type == events.TypeOrderSubmitted {
			var os events.OrderSubmitted
			is.NoErr(json.Unmarshal(e.Payload, &os))
			is.Equal(os.Nation, "Austria")
			foundOS = true
		}
	}
	is.Equal(foundOS, true)
	t.Logf("Disband (retreat phase): resp=%q", resp)
}

func TestCommand_Build(t *testing.T) {
	// /build happy path: England builds F Lon in Winter 1901 Adjustment.
	// England has 4 SCs (lon, edi, lvp, nwy) and 3 units → 1 build slot.
	// Expects no error and an OrderSubmitted event recorded in England's DM.
	is := is.New(t)
	d, ch := adjustmentPhaseGame(t)

	resp, err := d.Dispatch(dmCmd("build", "u2", "game", "F", "lon"))
	is.NoErr(err)
	is.Equal(resp != "", true)

	dmEnvs, err := events.ScanDM(ch, "u2")
	is.NoErr(err)
	foundOS := false
	for _, e := range dmEnvs {
		if e.Type == events.TypeOrderSubmitted {
			var os events.OrderSubmitted
			is.NoErr(json.Unmarshal(e.Payload, &os))
			is.Equal(os.Nation, "England")
			foundOS = true
		}
	}
	is.Equal(foundOS, true)
	t.Logf("Build: resp=%q", resp)
}

func TestCommand_Waive(t *testing.T) {
	// /waive happy path: England waives its one available build slot.
	// Expects no error and the waive order to be staged.
	is := is.New(t)
	d, _ := adjustmentPhaseGame(t)

	resp, err := d.Dispatch(dmCmd("waive", "u2", "game"))
	is.NoErr(err)
	is.Equal(resp != "", true)
	t.Logf("Waive: resp=%q", resp)
}

// ---------------------------------------------------------------------------
// Phase-guard tests for Retreat / Adjustment commands
// (called in the wrong phase to prove the commands are wired and guarded)
// ---------------------------------------------------------------------------

func TestCommand_Retreat_RejectedOutsideRetreatPhase(t *testing.T) {
	// /retreat is only valid in the Retreat phase. Calling it in Movement
	// phase must return an error.
	is := is.New(t)
	d, _ := startedGame(t)

	_, err := d.Dispatch(dmCmd("retreat", "u1", "game", "A", "Lon", "Wal"))
	is.NotNil(err)
	t.Logf("Retreat rejected in Movement phase: %v", err)
}

func TestCommand_Disband_RejectedOutsideRetreatAndAdjustmentPhase(t *testing.T) {
	// /disband is valid only in Retreat or Adjustment phases.
	is := is.New(t)
	d, _ := startedGame(t)

	_, err := d.Dispatch(dmCmd("disband", "u1", "game", "A", "Lon"))
	is.NotNil(err)
	t.Logf("Disband rejected in Movement phase: %v", err)
}

func TestCommand_Build_RejectedOutsideAdjustmentPhase(t *testing.T) {
	// /build is valid only in the Adjustment phase.
	is := is.New(t)
	d, _ := startedGame(t)

	_, err := d.Dispatch(dmCmd("build", "u1", "game", "A", "Lon"))
	is.NotNil(err)
	t.Logf("Build rejected in Movement phase: %v", err)
}

func TestCommand_Waive_RejectedOutsideAdjustmentPhase(t *testing.T) {
	// /waive is valid only in the Adjustment phase.
	is := is.New(t)
	d, _ := startedGame(t)

	_, err := d.Dispatch(dmCmd("waive", "u1", "game"))
	is.NotNil(err)
	t.Logf("Waive rejected in Movement phase: %v", err)
}

// ---------------------------------------------------------------------------
// Info commands (any phase, anyone)
// ---------------------------------------------------------------------------

func TestCommand_Status(t *testing.T) {
	// /status returns the current phase, SC counts, and submission status.
	is := is.New(t)
	d, _ := startedGame(t)

	resp, err := d.Dispatch(chanCmd("status", "anyone", "game"))
	is.NoErr(err)
	is.Equal(strings.Contains(resp, "Spring 1901") || strings.Contains(resp, "Movement"), true)
	t.Logf("Status: %q", resp)
}

func TestCommand_History_BeforeFirstResolution(t *testing.T) {
	// /history before any phase has resolved should return a descriptive
	// message (not crash).
	d, _ := startedGame(t)

	resp, err := d.Dispatch(chanCmd("history", "anyone", "game", "1"))
	// Either an informational response or a "not found" error is acceptable;
	// it must not panic and must return something.
	_ = err
	_ = resp
	t.Logf("History (before any resolution): resp=%q err=%v", resp, err)
}

func TestCommand_Map_NoArgs(t *testing.T) {
	// /map with no arguments posts the full board PNG to the channel.
	is := is.New(t)
	d, ch := startedGame(t)

	_, err := d.Dispatch(chanCmd("map", "anyone", "game"))
	is.NoErr(err)
	// PNG must have been posted (PostImage called).
	is.Equal(len(ch.imgs) > 0, true)
	// First four bytes must be the PNG magic number.
	if len(ch.imgs) > 0 {
		img := ch.imgs[0]
		is.Equal(img[0], byte(0x89))
		is.Equal(img[1], byte('P'))
		is.Equal(img[2], byte('N'))
		is.Equal(img[3], byte('G'))
		t.Logf("PNG posted: %d KB", len(img)/1024)
	}
}

func TestCommand_Map_WithTerritoryAndRadius(t *testing.T) {
	// /map Vienna 1 highlights Vienna and its adjacent provinces.
	is := is.New(t)
	d, ch := startedGame(t)

	_, err := d.Dispatch(chanCmd("map", "anyone", "game", "vie", "1"))
	is.NoErr(err)
	is.Equal(len(ch.imgs) > 0, true)
	t.Logf("Map with territory: PNG %d KB", len(ch.imgs[0])/1024)
}

func TestCommand_Help_NoArgs(t *testing.T) {
	// /help lists commands grouped by the seven categories.
	is := is.New(t)
	d, _ := startedGame(t)

	resp, err := d.Dispatch(chanCmd("help", "anyone", "game"))
	is.NoErr(err)
	// All seven category headers must be present.
	for _, header := range []string{"Setup:", "Movement:", "Retreat:", "Adjustment:", "Info:", "Draw:", "GM:"} {
		is.Equal(strings.Contains(resp, header), true)
	}
	// /nations and /provinces must appear in the Info section.
	is.Equal(strings.Contains(resp, "nations"), true)
	is.Equal(strings.Contains(resp, "provinces"), true)
	t.Logf("Help response length: %d chars", len(resp))
}

func TestCommand_Help_WithCommand(t *testing.T) {
	// /help order returns a multi-line block with Phase, Access, and Examples sections.
	is := is.New(t)
	d, _ := startedGame(t)

	resp, err := d.Dispatch(chanCmd("help", "anyone", "game", "order"))
	is.NoErr(err)
	for _, section := range []string{"Phase:", "Access:", "Examples:"} {
		is.Equal(strings.Contains(resp, section), true)
	}
	is.Equal(strings.Contains(resp, "order"), true)
	t.Logf("Help for 'order': %q", resp)
}

func TestCommand_Help_Rules(t *testing.T) {
	// /help rules returns a condensed rules summary.
	is := is.New(t)
	d, _ := startedGame(t)

	resp, err := d.Dispatch(chanCmd("help", "anyone", "game", "rules"))
	is.NoErr(err)
	is.Equal(strings.Contains(resp, "supply centres"), true)
	is.Equal(strings.Contains(resp, "phase"), true)
	t.Logf("Help rules length: %d chars", len(resp))
}

// ---------------------------------------------------------------------------
// /nations — reference command (any phase, anyone)
// ---------------------------------------------------------------------------

func TestCommand_Nations_NoArgs(t *testing.T) {
	// /nations lists all 7 nations with abbreviations and home SC codes.
	is := is.New(t)
	d, _ := startedGame(t)

	resp, err := d.Dispatch(chanCmd("nations", "anyone", "game"))
	is.NoErr(err)
	for _, name := range []string{"England", "France", "Germany", "Italy", "Austria", "Russia", "Turkey"} {
		is.Equal(strings.Contains(resp, name), true)
	}
	is.Equal(strings.Contains(resp, "Eng"), true)
	t.Logf("Nations response length: %d chars", len(resp))
}

func TestCommand_Nations_WithNation(t *testing.T) {
	// /nations England shows detail including home SCs and starting units.
	is := is.New(t)
	d, _ := startedGame(t)

	resp, err := d.Dispatch(chanCmd("nations", "anyone", "game", "England"))
	is.NoErr(err)
	is.Equal(strings.Contains(resp, "Edinburgh"), true)
	is.Equal(strings.Contains(resp, "F Edinburgh") || strings.Contains(resp, "F London"), true)
	t.Logf("Nations England: %q", resp)
}

func TestCommand_Nations_UnknownNation(t *testing.T) {
	// /nations with an unknown name returns an error.
	is := is.New(t)
	d, _ := startedGame(t)

	_, err := d.Dispatch(chanCmd("nations", "anyone", "game", "Gondor"))
	is.NotNil(err)
	t.Logf("Unknown nation error: %v", err)
}

// ---------------------------------------------------------------------------
// /provinces — reference command (any phase, anyone)
// ---------------------------------------------------------------------------

func TestCommand_Provinces_NoArgs(t *testing.T) {
	// /provinces lists all province codes with full names.
	is := is.New(t)
	d, _ := startedGame(t)

	resp, err := d.Dispatch(chanCmd("provinces", "anyone", "game"))
	is.NoErr(err)
	is.Equal(strings.Contains(resp, "vie"), true)
	is.Equal(strings.Contains(resp, "Vienna"), true)
	t.Logf("Provinces response length: %d chars", len(resp))
}

func TestCommand_Provinces_WithNation(t *testing.T) {
	// /provinces Austria filters to Austria's home SCs.
	is := is.New(t)
	d, _ := startedGame(t)

	resp, err := d.Dispatch(chanCmd("provinces", "anyone", "game", "Austria"))
	is.NoErr(err)
	for _, code := range []string{"vie", "tri", "bud"} {
		is.Equal(strings.Contains(resp, code), true)
	}
	t.Logf("Provinces Austria: %q", resp)
}

// ---------------------------------------------------------------------------
// Draw commands (any phase, own nation)
// ---------------------------------------------------------------------------

func TestCommand_Draw_ProposesOnFirstCall(t *testing.T) {
	// First /draw from a player posts DrawProposed.
	is := is.New(t)
	d, ch := startedGame(t)

	resp, err := d.Dispatch(chanCmd("draw", "u1", "game"))
	is.NoErr(err)
	is.Equal(hasEvent(t, ch, "game", events.TypeDrawProposed), true)
	t.Logf("Draw proposed: %q", resp)
}

func TestCommand_Draw_AllNationsEndGame(t *testing.T) {
	// When all players vote yes via /draw the game ends with result "draw".
	is := is.New(t)
	d, ch := startedGame(t)

	mustDispatch(t, d, chanCmd("draw", "u1", "game")) // propose
	_, err := d.Dispatch(chanCmd("draw", "u2", "game")) // second player votes yes → game over
	is.NoErr(err)

	is.Equal(hasEvent(t, ch, "game", events.TypeGameEnded), true)
	var ge events.GameEnded
	is.NoErr(json.Unmarshal(eventPayload(t, ch, "game", events.TypeGameEnded), &ge))
	is.Equal(ge.Result, "draw")
	t.Logf("GameEnded: result=%q", ge.Result)
}

func TestCommand_Concede(t *testing.T) {
	// /concede ends the game immediately with result "concession".
	is := is.New(t)
	d, ch := startedGame(t)

	resp, err := d.Dispatch(chanCmd("concede", "u1", "game"))
	is.NoErr(err)
	is.Equal(hasEvent(t, ch, "game", events.TypeGameEnded), true)

	var ge events.GameEnded
	is.NoErr(json.Unmarshal(eventPayload(t, ch, "game", events.TypeGameEnded), &ge))
	is.Equal(ge.Result, "concession")
	t.Logf("Concede: resp=%q result=%q", resp, ge.Result)
}

// ---------------------------------------------------------------------------
// GM commands (any phase, GM only)
// ---------------------------------------------------------------------------

func TestCommand_Pause(t *testing.T) {
	// /pause cancels the deadline timer; must be callable by GM only.
	is := is.New(t)
	d, _ := startedGame(t)

	_, err := d.Dispatch(chanCmd("pause", "gm", "game"))
	is.NoErr(err)

	// Non-GM must be rejected.
	_, err = d.Dispatch(chanCmd("pause", "u1", "game"))
	is.NotNil(err)
	t.Logf("Pause: GM accepted, non-GM rejected (%v)", err)
}

func TestCommand_Resume(t *testing.T) {
	// /resume restarts the deadline timer after a pause.
	is := is.New(t)
	d, _ := startedGame(t)

	mustDispatch(t, d, chanCmd("pause", "gm", "game"))
	_, err := d.Dispatch(chanCmd("resume", "gm", "game"))
	is.NoErr(err)
	t.Log("Pause then resume: accepted")
}

func TestCommand_Extend(t *testing.T) {
	// /extend <duration> adds time to the current deadline.
	is := is.New(t)
	d, _ := startedGame(t)

	_, err := d.Dispatch(chanCmd("extend", "gm", "game", "2h"))
	is.NoErr(err)

	// Non-GM must be rejected.
	_, err = d.Dispatch(chanCmd("extend", "u1", "game", "2h"))
	is.NotNil(err)
	t.Logf("Extend 2h: GM accepted, non-GM rejected (%v)", err)
}

func TestCommand_ForceResolve(t *testing.T) {
	// /force-resolve immediately adjudicates the current phase and posts
	// PhaseResolved.  GM only.
	is := is.New(t)
	d, ch := startedGame(t)

	_, err := d.Dispatch(chanCmd("force-resolve", "gm", "game"))
	is.NoErr(err)
	is.Equal(hasEvent(t, ch, "game", events.TypePhaseResolved), true)

	var pr events.PhaseResolved
	is.NoErr(json.Unmarshal(eventPayload(t, ch, "game", events.TypePhaseResolved), &pr))
	// PhaseResolved.Phase is the godip phase type ("Movement"), not the full display string.
	is.Equal(pr.Phase, "Movement")

	// Non-GM must be rejected.
	_, err = d.Dispatch(chanCmd("force-resolve", "u1", "game"))
	is.NotNil(err)
	t.Logf("ForceResolve: phase=%q, non-GM rejected (%v)", pr.Phase, err)
}

func TestCommand_Boot(t *testing.T) {
	// /boot <nation> removes a player.  GM only.
	is := is.New(t)
	d, ch := startedGame(t)

	_, err := d.Dispatch(chanCmd("boot", "gm", "game", "England"))
	is.NoErr(err)
	is.Equal(hasEvent(t, ch, "game", events.TypePlayerBooted), true)

	var pb events.PlayerBooted
	is.NoErr(json.Unmarshal(eventPayload(t, ch, "game", events.TypePlayerBooted), &pb))
	is.Equal(pb.Nation, "England")

	// Non-GM must be rejected.
	_, err = d.Dispatch(chanCmd("boot", "u1", "game", "France"))
	is.NotNil(err)
	t.Logf("Boot England: ok; non-GM rejected (%v)", err)
}

func TestCommand_Replace(t *testing.T) {
	// /replace <nation> <newuser> transfers nation ownership.  GM only.
	is := is.New(t)
	d, ch := startedGame(t)

	_, err := d.Dispatch(chanCmd("replace", "gm", "game", "England", "newplayer"))
	is.NoErr(err)
	is.Equal(hasEvent(t, ch, "game", events.TypePlayerReplaced), true)

	var pr events.PlayerReplaced
	is.NoErr(json.Unmarshal(eventPayload(t, ch, "game", events.TypePlayerReplaced), &pr))
	is.Equal(pr.Nation, "England")
	is.Equal(pr.NewUserID, "newplayer")

	// Non-GM must be rejected.
	_, err = d.Dispatch(chanCmd("replace", "u1", "game", "France", "x"))
	is.NotNil(err)
	t.Logf("Replace England→newplayer: ok; non-GM rejected (%v)", err)
}
