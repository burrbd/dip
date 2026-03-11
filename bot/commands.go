// Package bot contains the platform-agnostic command router. It receives
// parsed commands from platform adapters, enforces access control (only the
// assigned player may submit orders for their nation), and delegates to the
// session and engine packages.
package bot

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/burrbd/dip/dipmap"
	"github.com/burrbd/dip/engine"
	"github.com/burrbd/dip/events"
	"github.com/burrbd/dip/session"
	"github.com/zond/godip"
	"github.com/zond/godip/variants/classical"
	"github.com/zond/godip/variants/classical/start"
)

// classicalNations is the set of valid nation names in the classical variant.
var classicalNations = map[string]bool{
	"Austria": true,
	"England": true,
	"France":  true,
	"Germany": true,
	"Italy":   true,
	"Russia":  true,
	"Turkey":  true,
}

// Command is a parsed bot command from any platform adapter.
type Command struct {
	Name          string   // command name without leading slash (e.g. "newgame")
	Args          []string // positional arguments
	UserID        string   // platform user identifier
	ChannelID     string   // platform channel identifier (DM channel for DM commands)
	IsDM          bool     // true when the command was sent via direct message
	GameChannelID string   // game channel ID; must be set for DM commands
}

// EngineFactory creates a new game engine for the given Diplomacy variant name.
type EngineFactory func(variant string) (engine.Engine, error)

// Dispatcher routes commands to their handlers and holds in-process sessions.
type Dispatcher struct {
	ch             events.Channel
	notifier       session.Notifier
	loader         session.EngineLoader
	newEng         EngineFactory
	sessions       map[string]*session.Session
	graph          dipmap.Graph                                                // optional board graph for /map neighbourhood queries
	svgFn          func(dipmap.EngineState) ([]byte, error)                   // defaults to dipmap.LoadSVG (raw SVG bytes)
	overlayFn      func([]byte, map[string]dipmap.Unit) ([]byte, error)       // defaults to dipmap.Overlay (unit glyphs)
	imgFn          func([]byte) ([]byte, error)                               // defaults to dipmap.SVGToJPEG (full-board JPEG)
	highlightFn    func([]byte, []string) ([]byte, error)                     // defaults to dipmap.Highlight
	renderZoomedFn func(dipmap.EngineState, []byte, []string) ([]byte, error) // defaults to dipmap.RenderZoomed
}

// New returns a Dispatcher wired to the given dependencies.
func New(ch events.Channel, notifier session.Notifier, loader session.EngineLoader, newEng EngineFactory) *Dispatcher {
	return &Dispatcher{
		ch:             ch,
		notifier:       notifier,
		loader:         loader,
		newEng:         newEng,
		sessions:       make(map[string]*session.Session),
		svgFn:          dipmap.LoadSVG,
		overlayFn:      dipmap.Overlay,
		imgFn:          dipmap.SVGToJPEG,
		highlightFn:    dipmap.Highlight,
		renderZoomedFn: dipmap.RenderZoomed,
	}
}

// Dispatch routes cmd to the correct handler and returns a response text.
func (d *Dispatcher) Dispatch(cmd Command) (string, error) {
	switch cmd.Name {
	case "newgame":
		return d.handleNewGame(cmd)
	case "join":
		return d.handleJoin(cmd)
	case "start":
		return d.handleStart(cmd)
	case "order":
		return d.handleOrder(cmd)
	case "orders":
		return d.handleOrders(cmd)
	case "clear":
		return d.handleClear(cmd)
	case "submit":
		return d.handleSubmit(cmd)
	case "retreat":
		return d.handleRetreat(cmd)
	case "disband":
		return d.handleDisband(cmd)
	case "build":
		return d.handleBuild(cmd)
	case "waive":
		return d.handleWaive(cmd)
	case "status":
		return d.handleStatus(cmd)
	case "history":
		return d.handleHistory(cmd)
	case "map":
		return d.handleMap(cmd)
	case "help":
		return d.handleHelp(cmd)
	case "nations":
		return d.handleNations(cmd)
	case "provinces":
		return d.handleProvinces(cmd)
	case "draw":
		return d.handleDraw(cmd)
	case "concede":
		return d.handleConcede(cmd)
	case "pause":
		return d.handlePause(cmd)
	case "resume":
		return d.handleResume(cmd)
	case "extend":
		return d.handleExtend(cmd)
	case "force-resolve":
		return d.handleForceResolve(cmd)
	case "boot":
		return d.handleBoot(cmd)
	case "replace":
		return d.handleReplace(cmd)
	default:
		return "", fmt.Errorf("bot: unknown command %q", cmd.Name)
	}
}

// gameState is the bot's view of a channel's game state, derived from an event scan.
type gameState struct {
	created       bool
	started       bool
	ended         bool
	gmID          string
	deadlineHours int
	players       map[string]string // userID → nation
	nations       map[string]string // nation → userID
	drawProposed  bool
	drawVotes     map[string]bool // nation → true if voted yes
}

// readState scans the channel event log and returns the current game state.
func (d *Dispatcher) readState(channelID string) (*gameState, error) {
	envs, err := events.Scan(d.ch, channelID)
	if err != nil {
		return nil, fmt.Errorf("bot: scan channel: %w", err)
	}
	gs := &gameState{
		players:       make(map[string]string),
		nations:       make(map[string]string),
		drawVotes:     make(map[string]bool),
		deadlineHours: 24,
	}
	for _, env := range envs {
		switch env.Type {
		case events.TypeGameCreated:
			var gc events.GameCreated
			if err := json.Unmarshal(env.Payload, &gc); err != nil {
				continue
			}
			gs.created = true
			gs.gmID = gc.GMUserID
			if gc.DeadlineHours > 0 {
				gs.deadlineHours = gc.DeadlineHours
			}
		case events.TypeGameStarted:
			gs.started = true
		case events.TypePlayerJoined:
			var pj events.PlayerJoined
			if err := json.Unmarshal(env.Payload, &pj); err != nil {
				continue
			}
			gs.players[pj.UserID] = pj.Nation
			gs.nations[pj.Nation] = pj.UserID
		case events.TypeGameEnded:
			gs.ended = true
			gs.drawProposed = false
			gs.drawVotes = make(map[string]bool)
		case events.TypeDrawProposed:
			var dp events.DrawProposed
			if err := json.Unmarshal(env.Payload, &dp); err != nil {
				continue
			}
			gs.drawProposed = true
			gs.drawVotes[dp.ProposerNation] = true
		case events.TypeDrawVoted:
			var dv events.DrawVoted
			if err := json.Unmarshal(env.Payload, &dv); err != nil {
				continue
			}
			if dv.Accept {
				gs.drawVotes[dv.Nation] = true
			} else {
				delete(gs.drawVotes, dv.Nation)
			}
		case events.TypePlayerBooted:
			var pb events.PlayerBooted
			if err := json.Unmarshal(env.Payload, &pb); err != nil {
				continue
			}
			for uid, nation := range gs.players {
				if nation == pb.Nation {
					delete(gs.players, uid)
					delete(gs.nations, pb.Nation)
					break
				}
			}
		case events.TypePlayerReplaced:
			var pr events.PlayerReplaced
			if err := json.Unmarshal(env.Payload, &pr); err != nil {
				continue
			}
			oldUID := gs.nations[pr.Nation]
			delete(gs.players, oldUID)
			gs.players[pr.NewUserID] = pr.Nation
			gs.nations[pr.Nation] = pr.NewUserID
		}
	}
	return gs, nil
}

// handleNewGame processes /newgame [settings].
func (d *Dispatcher) handleNewGame(cmd Command) (string, error) {
	state, err := d.readState(cmd.ChannelID)
	if err != nil {
		return "", err
	}
	if state.created {
		return "", fmt.Errorf("bot: a game already exists in this channel")
	}
	if err := events.Write(d.ch, cmd.ChannelID, events.TypeGameCreated, events.GameCreated{
		Variant:       "classical",
		DeadlineHours: 24,
		GMUserID:      cmd.UserID,
	}); err != nil {
		return "", fmt.Errorf("bot: write GameCreated: %w", err)
	}
	return "Game created. You are the GM. Players can use /join <nation> to claim a nation. Use /start when everyone has joined.", nil
}

// handleJoin processes /join <nation>.
func (d *Dispatcher) handleJoin(cmd Command) (string, error) {
	state, err := d.readState(cmd.ChannelID)
	if err != nil {
		return "", err
	}
	if !state.created {
		return "", fmt.Errorf("bot: no game in this channel; use /newgame first")
	}
	if state.started {
		return "", fmt.Errorf("bot: game has already started")
	}
	if _, already := state.players[cmd.UserID]; already {
		return "", fmt.Errorf("bot: you have already joined as %s", state.players[cmd.UserID])
	}
	if len(cmd.Args) == 0 {
		return "", fmt.Errorf("bot: usage: /join <nation>")
	}
	nation := cmd.Args[0]
	if !classicalNations[nation] {
		return "", fmt.Errorf("bot: unknown nation %q; valid nations are Austria, England, France, Germany, Italy, Russia, Turkey", nation)
	}
	if _, taken := state.nations[nation]; taken {
		return "", fmt.Errorf("bot: nation %q is already taken", nation)
	}
	if err := events.Write(d.ch, cmd.ChannelID, events.TypePlayerJoined, events.PlayerJoined{
		UserID: cmd.UserID,
		Nation: nation,
	}); err != nil {
		return "", fmt.Errorf("bot: write PlayerJoined: %w", err)
	}
	return fmt.Sprintf("Joined as %s.", nation), nil
}

// handleStart processes /start (GM only).
func (d *Dispatcher) handleStart(cmd Command) (string, error) {
	state, err := d.readState(cmd.ChannelID)
	if err != nil {
		return "", err
	}
	if !state.created {
		return "", fmt.Errorf("bot: no game in this channel; use /newgame first")
	}
	if state.started {
		return "", fmt.Errorf("bot: game has already started")
	}
	if cmd.UserID != state.gmID {
		return "", fmt.Errorf("bot: only the GM can start the game")
	}
	n := len(state.players)
	if n < 2 {
		return "", fmt.Errorf("bot: need at least 2 players to start (have %d)", n)
	}
	if n > 7 {
		return "", fmt.Errorf("bot: too many players (max 7, have %d)", n)
	}
	eng, err := d.newEng("classical")
	if err != nil {
		return "", fmt.Errorf("bot: create engine: %w", err)
	}
	snapshot, err := eng.Dump()
	if err != nil {
		return "", fmt.Errorf("bot: dump initial state: %w", err)
	}
	if err := events.Write(d.ch, cmd.ChannelID, events.TypeGameStarted, events.GameStarted{
		InitialState: json.RawMessage(snapshot),
	}); err != nil {
		return "", fmt.Errorf("bot: write GameStarted: %w", err)
	}
	sess := session.New(d.ch, cmd.ChannelID, state.gmID, eng.Phase(), state.players, state.deadlineHours, eng, d.notifier)
	d.sessions[cmd.ChannelID] = sess
	return "Game started! Spring 1901 Movement phase begins. Players, submit your orders via DM.", nil
}

// isMovementPhase returns true if the given phase string is a Movement phase.
func isMovementPhase(phase string) bool {
	return strings.HasSuffix(phase, "Movement")
}

// handleOrder processes /order <order-text> (DM only, Movement phase).
// It parses the order via the engine, validates nation ownership, and stages it.
func (d *Dispatcher) handleOrder(cmd Command) (string, error) {
	if !cmd.IsDM {
		return "", fmt.Errorf("bot: /order must be sent as a direct message to the bot")
	}
	sess, ok := d.sessions[cmd.GameChannelID]
	if !ok || sess == nil {
		return "", fmt.Errorf("bot: no active game found")
	}
	if !isMovementPhase(sess.Phase) {
		return "", fmt.Errorf("bot: /order is only valid during the Movement phase (current: %s)", sess.Phase)
	}
	nation, ok := sess.Players[cmd.UserID]
	if !ok {
		return "", fmt.Errorf("bot: you are not a player in this game")
	}
	if len(cmd.Args) == 0 {
		return "", fmt.Errorf("bot: usage: /order <order-text>")
	}
	orderText := strings.Join(cmd.Args, " ")
	if err := sess.Eng.SubmitOrder(nation, orderText); err != nil {
		return "", fmt.Errorf("bot: invalid order: %w", err)
	}
	sess.StagedOrders[nation] = append(sess.StagedOrders[nation], orderText)
	return fmt.Sprintf("Order staged: %s", orderText), nil
}

// handleOrders processes /orders (DM only) — lists the caller's staged orders.
func (d *Dispatcher) handleOrders(cmd Command) (string, error) {
	if !cmd.IsDM {
		return "", fmt.Errorf("bot: /orders must be sent as a direct message to the bot")
	}
	sess, ok := d.sessions[cmd.GameChannelID]
	if !ok || sess == nil {
		return "", fmt.Errorf("bot: no active game found")
	}
	nation, ok := sess.Players[cmd.UserID]
	if !ok {
		return "", fmt.Errorf("bot: you are not a player in this game")
	}
	orders := sess.StagedOrders[nation]
	if len(orders) == 0 {
		return "No orders staged.", nil
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "Staged orders for %s:\n", nation)
	for i, o := range orders {
		fmt.Fprintf(&sb, "  %d. %s\n", i+1, o)
	}
	return strings.TrimRight(sb.String(), "\n"), nil
}

// handleClear processes /clear [order] (DM only) — removes one or all staged orders.
func (d *Dispatcher) handleClear(cmd Command) (string, error) {
	if !cmd.IsDM {
		return "", fmt.Errorf("bot: /clear must be sent as a direct message to the bot")
	}
	sess, ok := d.sessions[cmd.GameChannelID]
	if !ok || sess == nil {
		return "", fmt.Errorf("bot: no active game found")
	}
	nation, ok := sess.Players[cmd.UserID]
	if !ok {
		return "", fmt.Errorf("bot: you are not a player in this game")
	}
	if len(cmd.Args) == 0 {
		sess.StagedOrders[nation] = nil
		sess.Submitted[nation] = false
		return "All orders cleared.", nil
	}
	target := strings.Join(cmd.Args, " ")
	var filtered []string
	found := false
	for _, o := range sess.StagedOrders[nation] {
		if o == target && !found {
			found = true
			continue
		}
		filtered = append(filtered, o)
	}
	if !found {
		return "", fmt.Errorf("bot: order %q not found", target)
	}
	sess.StagedOrders[nation] = filtered
	sess.Submitted[nation] = false
	return fmt.Sprintf("Order removed: %s", target), nil
}

// handleSubmit processes /submit (DM only) — finalises staged orders and checks
// whether all nations have submitted; if so fires AdvanceTurn immediately.
func (d *Dispatcher) handleSubmit(cmd Command) (string, error) {
	if !cmd.IsDM {
		return "", fmt.Errorf("bot: /submit must be sent as a direct message to the bot")
	}
	sess, ok := d.sessions[cmd.GameChannelID]
	if !ok || sess == nil {
		return "", fmt.Errorf("bot: no active game found")
	}
	if !isMovementPhase(sess.Phase) {
		return "", fmt.Errorf("bot: /submit is only valid during the Movement phase (current: %s)", sess.Phase)
	}
	nation, ok := sess.Players[cmd.UserID]
	if !ok {
		return "", fmt.Errorf("bot: you are not a player in this game")
	}
	if err := events.WriteDM(d.ch, cmd.UserID, events.TypeOrderSubmitted, events.OrderSubmitted{
		UserID: cmd.UserID,
		Nation: nation,
		Orders: sess.StagedOrders[nation],
		Phase:  sess.Phase,
	}); err != nil {
		return "", fmt.Errorf("bot: write OrderSubmitted: %w", err)
	}
	sess.Submitted[nation] = true

	allDone, err := d.allNationsSubmitted(sess)
	if err != nil {
		return "", fmt.Errorf("bot: check submissions: %w", err)
	}
	if allDone {
		if err := sess.AdvanceTurn(); err != nil {
			return "", fmt.Errorf("bot: advance turn: %w", err)
		}
		return "Orders submitted. All nations ready — resolving now!", nil
	}
	return "Orders submitted.", nil
}

// isRetreatPhase returns true if the given phase string is a Retreat phase.
func isRetreatPhase(phase string) bool {
	return strings.HasSuffix(phase, "Retreat")
}

// isAdjustmentPhase returns true if the given phase string is an Adjustment phase.
func isAdjustmentPhase(phase string) bool {
	return strings.HasSuffix(phase, "Adjustment")
}

// handleRetreat processes /retreat <unit_type> <source> <destination> (DM only, Retreat phase).
// It validates the source province has a dislodged unit belonging to the caller's nation,
// then stages a retreat order.
func (d *Dispatcher) handleRetreat(cmd Command) (string, error) {
	if !cmd.IsDM {
		return "", fmt.Errorf("bot: /retreat must be sent as a direct message to the bot")
	}
	sess, ok := d.sessions[cmd.GameChannelID]
	if !ok || sess == nil {
		return "", fmt.Errorf("bot: no active game found")
	}
	if !isRetreatPhase(sess.Phase) {
		return "", fmt.Errorf("bot: /retreat is only valid during the Retreat phase (current: %s)", sess.Phase)
	}
	nation, ok := sess.Players[cmd.UserID]
	if !ok {
		return "", fmt.Errorf("bot: you are not a player in this game")
	}
	if len(cmd.Args) < 3 {
		return "", fmt.Errorf("bot: usage: /retreat <unit_type> <source> <destination>")
	}
	unitType, src, dest := cmd.Args[0], cmd.Args[1], cmd.Args[2]
	src = strings.ToLower(src)
	dislodgeds := sess.Eng.Dislodgeds()
	if n, exists := dislodgeds[src]; !exists || n != nation {
		return "", fmt.Errorf("bot: no dislodged %s unit at %s belonging to %s", unitType, src, nation)
	}
	orderText := fmt.Sprintf("%s %s-%s", unitType, src, dest)
	if err := sess.Eng.SubmitOrder(nation, orderText); err != nil {
		return "", fmt.Errorf("bot: invalid retreat order: %w", err)
	}
	sess.StagedOrders[nation] = append(sess.StagedOrders[nation], orderText)
	if err := events.WriteDM(d.ch, cmd.UserID, events.TypeOrderSubmitted, events.OrderSubmitted{
		UserID: cmd.UserID,
		Nation: nation,
		Orders: []string{orderText},
		Phase:  sess.Phase,
	}); err != nil {
		return "", fmt.Errorf("bot: write OrderSubmitted: %w", err)
	}
	sess.Submitted[nation] = true
	if allRetreatActionsSubmitted(sess) {
		if err := sess.AdvanceTurn(); err != nil {
			return "", fmt.Errorf("bot: advance turn: %w", err)
		}
		return fmt.Sprintf("Retreat order staged: %s. All required orders received — resolving now!", orderText), nil
	}
	return fmt.Sprintf("Retreat order staged: %s", orderText), nil
}

// handleDisband processes /disband <unit_type> <province> (DM only, Retreat or Adjustment phase).
// During the Retreat phase it validates the province has a dislodged unit belonging to the nation.
func (d *Dispatcher) handleDisband(cmd Command) (string, error) {
	if !cmd.IsDM {
		return "", fmt.Errorf("bot: /disband must be sent as a direct message to the bot")
	}
	sess, ok := d.sessions[cmd.GameChannelID]
	if !ok || sess == nil {
		return "", fmt.Errorf("bot: no active game found")
	}
	if !isRetreatPhase(sess.Phase) && !isAdjustmentPhase(sess.Phase) {
		return "", fmt.Errorf("bot: /disband is only valid during Retreat or Adjustment phase (current: %s)", sess.Phase)
	}
	nation, ok := sess.Players[cmd.UserID]
	if !ok {
		return "", fmt.Errorf("bot: you are not a player in this game")
	}
	if len(cmd.Args) < 2 {
		return "", fmt.Errorf("bot: usage: /disband <unit_type> <province>")
	}
	unitType, src := cmd.Args[0], cmd.Args[1]
	src = strings.ToLower(src)
	if isRetreatPhase(sess.Phase) {
		dislodgeds := sess.Eng.Dislodgeds()
		if n, exists := dislodgeds[src]; !exists || n != nation {
			return "", fmt.Errorf("bot: no dislodged %s unit at %s belonging to %s", unitType, src, nation)
		}
	}
	orderText := fmt.Sprintf("%s %s disband", unitType, src)
	if err := sess.Eng.SubmitOrder(nation, orderText); err != nil {
		return "", fmt.Errorf("bot: invalid disband order: %w", err)
	}
	sess.StagedOrders[nation] = append(sess.StagedOrders[nation], orderText)
	if err := events.WriteDM(d.ch, cmd.UserID, events.TypeOrderSubmitted, events.OrderSubmitted{
		UserID: cmd.UserID,
		Nation: nation,
		Orders: []string{orderText},
		Phase:  sess.Phase,
	}); err != nil {
		return "", fmt.Errorf("bot: write OrderSubmitted: %w", err)
	}
	sess.Submitted[nation] = true
	if isRetreatPhase(sess.Phase) {
		if allRetreatActionsSubmitted(sess) {
			if err := sess.AdvanceTurn(); err != nil {
				return "", fmt.Errorf("bot: advance turn: %w", err)
			}
			return fmt.Sprintf("Disband order staged: %s. All required orders received — resolving now!", orderText), nil
		}
	} else {
		if allAdjustmentActionsSubmitted(sess) {
			if err := sess.AdvanceTurn(); err != nil {
				return "", fmt.Errorf("bot: advance turn: %w", err)
			}
			return fmt.Sprintf("Disband order staged: %s. All required orders received — resolving now!", orderText), nil
		}
	}
	return fmt.Sprintf("Disband order staged: %s", orderText), nil
}

// handleBuild processes /build <unit_type> <province> (DM only, Adjustment phase).
// It stages a build order for the caller's nation.
func (d *Dispatcher) handleBuild(cmd Command) (string, error) {
	if !cmd.IsDM {
		return "", fmt.Errorf("bot: /build must be sent as a direct message to the bot")
	}
	sess, ok := d.sessions[cmd.GameChannelID]
	if !ok || sess == nil {
		return "", fmt.Errorf("bot: no active game found")
	}
	if !isAdjustmentPhase(sess.Phase) {
		return "", fmt.Errorf("bot: /build is only valid during the Adjustment phase (current: %s)", sess.Phase)
	}
	nation, ok := sess.Players[cmd.UserID]
	if !ok {
		return "", fmt.Errorf("bot: you are not a player in this game")
	}
	if len(cmd.Args) < 2 {
		return "", fmt.Errorf("bot: usage: /build <unit_type> <province>")
	}
	unitType, province := cmd.Args[0], cmd.Args[1]
	orderText := fmt.Sprintf("build %s %s", unitType, province)
	if err := sess.Eng.SubmitOrder(nation, orderText); err != nil {
		return "", fmt.Errorf("bot: invalid build order: %w", err)
	}
	sess.StagedOrders[nation] = append(sess.StagedOrders[nation], orderText)
	if err := events.WriteDM(d.ch, cmd.UserID, events.TypeOrderSubmitted, events.OrderSubmitted{
		UserID: cmd.UserID,
		Nation: nation,
		Orders: []string{orderText},
		Phase:  sess.Phase,
	}); err != nil {
		return "", fmt.Errorf("bot: write OrderSubmitted: %w", err)
	}
	sess.Submitted[nation] = true
	if allAdjustmentActionsSubmitted(sess) {
		if err := sess.AdvanceTurn(); err != nil {
			return "", fmt.Errorf("bot: advance turn: %w", err)
		}
		return fmt.Sprintf("Build order staged: %s. All required orders received — resolving now!", orderText), nil
	}
	return fmt.Sprintf("Build order staged: %s", orderText), nil
}

// handleWaive processes /waive (DM only, Adjustment phase).
// It stages a waive order for one available build slot, without engine validation
// since a waive has no associated unit province.
func (d *Dispatcher) handleWaive(cmd Command) (string, error) {
	if !cmd.IsDM {
		return "", fmt.Errorf("bot: /waive must be sent as a direct message to the bot")
	}
	sess, ok := d.sessions[cmd.GameChannelID]
	if !ok || sess == nil {
		return "", fmt.Errorf("bot: no active game found")
	}
	if !isAdjustmentPhase(sess.Phase) {
		return "", fmt.Errorf("bot: /waive is only valid during the Adjustment phase (current: %s)", sess.Phase)
	}
	nation, ok := sess.Players[cmd.UserID]
	if !ok {
		return "", fmt.Errorf("bot: you are not a player in this game")
	}
	sess.StagedOrders[nation] = append(sess.StagedOrders[nation], "Waive")
	if err := events.WriteDM(d.ch, cmd.UserID, events.TypeOrderSubmitted, events.OrderSubmitted{
		UserID: cmd.UserID,
		Nation: nation,
		Orders: []string{"Waive"},
		Phase:  sess.Phase,
	}); err != nil {
		return "", fmt.Errorf("bot: write OrderSubmitted: %w", err)
	}
	sess.Submitted[nation] = true
	if allAdjustmentActionsSubmitted(sess) {
		if err := sess.AdvanceTurn(); err != nil {
			return "", fmt.Errorf("bot: advance turn: %w", err)
		}
		return "Waive order staged. All required orders received — resolving now!", nil
	}
	return "Waive order staged.", nil
}

// handleStatus processes /status — shows current phase, SC counts, and
// order submission status per nation. Requires an active in-memory session.
func (d *Dispatcher) handleStatus(cmd Command) (string, error) {
	sess, ok := d.sessions[cmd.ChannelID]
	if !ok || sess == nil {
		return "", fmt.Errorf("bot: no active game found in this channel")
	}
	scCounts := sess.Eng.SupplyCenters()
	return FormatStatus(sess.Phase, sess.Players, sess.Submitted, scCounts), nil
}

// handleHistory processes /history <turn> — fetches the PhaseResolved event
// matching the given turn string and returns its result summary.
func (d *Dispatcher) handleHistory(cmd Command) (string, error) {
	if len(cmd.Args) == 0 {
		return "", fmt.Errorf("bot: usage: /history <turn>")
	}
	turn := strings.Join(cmd.Args, " ")

	envs, err := events.Scan(d.ch, cmd.ChannelID)
	if err != nil {
		return "", fmt.Errorf("bot: scan history: %w", err)
	}

	// Search in reverse order so the most recent matching phase is preferred.
	for i := len(envs) - 1; i >= 0; i-- {
		if envs[i].Type != events.TypePhaseResolved {
			continue
		}
		var pr events.PhaseResolved
		if err := json.Unmarshal(envs[i].Payload, &pr); err != nil {
			continue
		}
		if strings.Contains(pr.Phase, turn) {
			if len(pr.ResultSummary) > 0 {
				return string(pr.ResultSummary), nil
			}
			return fmt.Sprintf("Phase %s resolved (no result summary).", pr.Phase), nil
		}
	}
	return "", fmt.Errorf("bot: no history found for turn %q", turn)
}

// godipGraph adapts the classical godip board graph to the dipmap.Graph
// interface used for neighbourhood BFS queries.
type godipGraph struct{ g godip.Graph }

// Edges returns all province names directly adjacent to territory.
func (gg godipGraph) Edges(territory string) []string {
	em := gg.g.Edges(godip.Province(territory), false)
	result := make([]string, 0, len(em))
	for p := range em {
		result = append(result, string(p))
	}
	return result
}

// boardGraph returns the Dispatcher's graph. If none is set it falls back
// to the classical variant's board graph so that /map radius queries work
// correctly out of the box.
func (d *Dispatcher) boardGraph() dipmap.Graph {
	if d.graph != nil {
		return d.graph
	}
	return godipGraph{g: start.Graph()}
}

// handleMap processes /map [territory [n]] — renders the board with unit
// positions and posts it as an image. With territory and n > 0, highlights the
// neighbourhood and crops to a zoomed view. Otherwise the full board is shown.
//
// Pipeline (both paths):
//  1. svgFn   — load raw SVG asset
//  2. overlayFn — inject army/fleet glyphs at province centroids
//  3a. Full board: imgFn — rasterise to JPEG
//  3b. Zoomed:    highlightFn → renderZoomedFn — highlight + crop → JPEG
func (d *Dispatcher) handleMap(cmd Command) (string, error) {
	sess, ok := d.sessions[cmd.ChannelID]
	if !ok || sess == nil {
		return "", fmt.Errorf("bot: no active game found in this channel")
	}

	territory := ""
	if len(cmd.Args) >= 1 {
		territory = cmd.Args[0]
	}
	n := 0
	if len(cmd.Args) >= 2 {
		parsed, err := strconv.Atoi(cmd.Args[1])
		if err != nil {
			return "", fmt.Errorf("bot: invalid radius %q: must be an integer", cmd.Args[1])
		}
		n = parsed
	}

	// Step 1: load SVG.
	svg, err := d.svgFn(sess.Eng)
	if err != nil {
		return "", fmt.Errorf("bot: render map: %w", err)
	}

	// Step 2: overlay unit positions.
	engUnits := sess.Eng.Units()
	units := make(map[string]dipmap.Unit, len(engUnits))
	for p, u := range engUnits {
		units[p] = dipmap.Unit{Type: u.Type, Nation: u.Nation}
	}
	svg, err = d.overlayFn(svg, units)
	if err != nil {
		return "", fmt.Errorf("bot: render map: %w", err)
	}

	// Step 3: rasterise.
	var img []byte
	if territory != "" && n > 0 {
		provinces := dipmap.Neighborhood(d.boardGraph(), territory, n)
		svg, err = d.highlightFn(svg, provinces)
		if err != nil {
			return "", fmt.Errorf("bot: render map: %w", err)
		}
		img, err = d.renderZoomedFn(sess.Eng, svg, provinces)
		if err != nil {
			return "", fmt.Errorf("bot: render map: %w", err)
		}
	} else {
		img, err = d.imgFn(svg)
		if err != nil {
			return "", fmt.Errorf("bot: render map: %w", err)
		}
	}

	if err := d.ch.PostImage(cmd.ChannelID, img); err != nil {
		return "", fmt.Errorf("bot: post map: %w", err)
	}
	return "Map posted.", nil
}

// commandDetail holds the structured help text for a single command.
type commandDetail struct {
	usage       string
	description string
	phase       string
	access      string
	examples    []string
}

// commandDetails maps command names to their detailed help information.
var commandDetails = map[string]commandDetail{
	"newgame": {
		usage:       "/newgame",
		description: "Start a new game in this channel. You become the GM.",
		phase:       "Any (pre-game)",
		access:      "Anyone",
		examples:    []string{"/newgame"},
	},
	"join": {
		usage:       "/join <nation>",
		description: "Join the game as the specified nation.",
		phase:       "Any (pre-game)",
		access:      "Anyone",
		examples:    []string{"/join England", "/join France", "/join Austria"},
	},
	"start": {
		usage:       "/start",
		description: "Start the game. Requires 2–7 players to have joined.",
		phase:       "Any (pre-game)",
		access:      "GM",
		examples:    []string{"/start"},
	},
	"order": {
		usage:       "/order <order-text>",
		description: "Submit a movement order for your nation.",
		phase:       "Movement",
		access:      "Own nation (DM only)",
		examples:    []string{"/order A Vie-Bud", "/order F Lon-NTH", "/order A Par S A Mar-Bur"},
	},
	"orders": {
		usage:       "/orders",
		description: "List your staged orders for the current phase.",
		phase:       "Movement",
		access:      "Own nation (DM only)",
		examples:    []string{"/orders"},
	},
	"clear": {
		usage:       "/clear [order]",
		description: "Clear all staged orders or remove a specific one.",
		phase:       "Movement",
		access:      "Own nation (DM only)",
		examples:    []string{"/clear", "/clear A Vie-Bud"},
	},
	"submit": {
		usage:       "/submit",
		description: "Finalise and submit your orders. If all nations submit, the phase resolves immediately.",
		phase:       "Movement",
		access:      "Own nation (DM only)",
		examples:    []string{"/submit"},
	},
	"retreat": {
		usage:       "/retreat <unit_type> <source> <destination>",
		description: "Retreat a dislodged unit to a valid adjacent province.",
		phase:       "Retreat",
		access:      "Own nation (DM only)",
		examples:    []string{"/retreat F Tri Adr", "/retreat A Bur Par"},
	},
	"disband": {
		usage:       "/disband <unit_type> <province>",
		description: "Disband a unit. In Retreat phase disbands a dislodged unit; in Adjustment phase removes an excess unit.",
		phase:       "Retreat or Adjustment",
		access:      "Own nation (DM only)",
		examples:    []string{"/disband F Tri", "/disband A Mun"},
	},
	"build": {
		usage:       "/build <unit_type> <province>",
		description: "Build a new unit in a home supply centre.",
		phase:       "Adjustment",
		access:      "Own nation (DM only)",
		examples:    []string{"/build F Lon", "/build A Ber"},
	},
	"waive": {
		usage:       "/waive",
		description: "Waive one available build slot.",
		phase:       "Adjustment",
		access:      "Own nation (DM only)",
		examples:    []string{"/waive"},
	},
	"status": {
		usage:       "/status",
		description: "Show current phase, supply centre counts, and order submission status per nation.",
		phase:       "Any",
		access:      "Anyone",
		examples:    []string{"/status"},
	},
	"history": {
		usage:       "/history <turn>",
		description: "Show adjudication results for a past turn.",
		phase:       "Any",
		access:      "Anyone",
		examples:    []string{"/history Spring 1901", "/history Fall 1902"},
	},
	"map": {
		usage:       "/map [territory [n]]",
		description: "Post the board map. With territory and n, highlights that province and all within n hops.",
		phase:       "Any",
		access:      "Anyone",
		examples:    []string{"/map", "/map Vienna 1", "/map vie 2"},
	},
	"help": {
		usage:       "/help [command|rules]",
		description: "List all commands grouped by category, show detailed help for a command, or display game rules.",
		phase:       "Any",
		access:      "Anyone",
		examples:    []string{"/help", "/help order", "/help rules"},
	},
	"nations": {
		usage:       "/nations [nation]",
		description: "List all classical powers with abbreviations and home SCs, or show detail for one nation.",
		phase:       "Any",
		access:      "Anyone",
		examples:    []string{"/nations", "/nations England", "/nations Eng"},
	},
	"provinces": {
		usage:       "/provinces [nation]",
		description: "List all province codes and full names, or filter to a nation's home supply centres.",
		phase:       "Any",
		access:      "Anyone",
		examples:    []string{"/provinces", "/provinces Austria", "/provinces Aus"},
	},
	"draw": {
		usage:       "/draw",
		description: "Propose a draw, or vote yes on an active draw proposal. Game ends when all nations agree.",
		phase:       "Any",
		access:      "Own nation",
		examples:    []string{"/draw"},
	},
	"concede": {
		usage:       "/concede",
		description: "Concede the game immediately.",
		phase:       "Any",
		access:      "Own nation",
		examples:    []string{"/concede"},
	},
	"pause": {
		usage:       "/pause",
		description: "Pause the phase deadline timer.",
		phase:       "Any",
		access:      "GM",
		examples:    []string{"/pause"},
	},
	"resume": {
		usage:       "/resume",
		description: "Resume a paused deadline timer.",
		phase:       "Any",
		access:      "GM",
		examples:    []string{"/resume"},
	},
	"extend": {
		usage:       "/extend <duration>",
		description: "Extend the current deadline by the given duration (e.g. 2h, 30m).",
		phase:       "Any",
		access:      "GM",
		examples:    []string{"/extend 2h", "/extend 30m"},
	},
	"force-resolve": {
		usage:       "/force-resolve",
		description: "Resolve the current phase immediately without waiting for the deadline.",
		phase:       "Any",
		access:      "GM",
		examples:    []string{"/force-resolve"},
	},
	"boot": {
		usage:       "/boot <nation>",
		description: "Remove a player from the game. Their units receive NMR orders going forward.",
		phase:       "Any",
		access:      "GM",
		examples:    []string{"/boot England"},
	},
	"replace": {
		usage:       "/replace <nation> <user>",
		description: "Transfer a nation to a new player.",
		phase:       "Any",
		access:      "GM",
		examples:    []string{"/replace England newplayer"},
	},
}

// helpCategories defines the seven categories and their command members in display order.
var helpCategories = []struct {
	name     string
	commands []string
}{
	{"Setup", []string{"newgame", "join", "start"}},
	{"Movement", []string{"order", "orders", "clear", "submit"}},
	{"Retreat", []string{"retreat", "disband"}},
	{"Adjustment", []string{"build", "disband", "waive"}},
	{"Info", []string{"status", "history", "map", "help", "nations", "provinces"}},
	{"Draw", []string{"draw", "concede"}},
	{"GM", []string{"pause", "resume", "extend", "force-resolve", "boot", "replace"}},
}

// commandList defines the canonical display order for /help (used for coverage checks).
var commandList = []string{
	"newgame", "join", "start",
	"order", "orders", "clear", "submit",
	"retreat", "disband", "build", "waive",
	"status", "history", "map", "help", "nations", "provinces",
	"draw", "concede",
	"pause", "resume", "extend", "force-resolve", "boot", "replace",
}

// helpRules is the condensed game rules overview returned by /help rules.
const helpRules = `Diplomacy — Quick Rules

Powers: Austria, England, France, Germany, Italy, Russia, Turkey (7 classical powers)
Win condition: Control 18 of 34 supply centres (SCs).

Phase sequence (repeating):
  Spring Movement → Spring Retreat → Fall Movement → Fall Retreat → Winter Adjustment → repeat

Orders (Movement phase, via DM):
  Move:         A Vie-Bud         (army in Vienna moves to Budapest)
  Hold:         A Vie H           (army holds position)
  Support hold: A Tri S A Vie     (Trieste supports Vienna's hold)
  Support move: A Tri S A Vie-Bud (Trieste supports Vienna's attack)
  Convoy:       F ADR C A Vie-Gre (fleet convoyes army across sea)

NMR (No Moves Received): unsubmitted orders become holds; unordered retreat units
are auto-disbanded; unordered build slots are waived.

Draw: any player may propose a draw with /draw; all remaining nations must agree
with /draw for the game to end in a draw.
Concede: a player may end the game immediately with /concede.`

// handleHelp processes /help [command|rules] — lists all commands grouped by category,
// shows detailed usage for a specific command, or returns the rules overview.
func (d *Dispatcher) handleHelp(cmd Command) (string, error) {
	if len(cmd.Args) == 0 {
		var sb strings.Builder
		for _, cat := range helpCategories {
			fmt.Fprintf(&sb, "%s:", cat.name)
			cmds := make([]string, 0, len(cat.commands))
			for _, c := range cat.commands {
				cmds = append(cmds, "/"+c)
			}
			fmt.Fprintf(&sb, "  %s\n", strings.Join(cmds, ", "))
		}
		return strings.TrimRight(sb.String(), "\n"), nil
	}

	arg := cmd.Args[0]
	if strings.HasPrefix(arg, "/") {
		arg = arg[1:]
	}

	if arg == "rules" {
		return helpRules, nil
	}

	det, ok := commandDetails[arg]
	if !ok {
		return "", fmt.Errorf("bot: unknown command %q; use /help for a list", cmd.Args[0])
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%s\n", det.usage)
	fmt.Fprintf(&sb, "  %s\n", det.description)
	fmt.Fprintf(&sb, "  Phase:   %s\n", det.phase)
	fmt.Fprintf(&sb, "  Access:  %s\n", det.access)
	fmt.Fprintf(&sb, "  Examples:\n")
	for _, ex := range det.examples {
		fmt.Fprintf(&sb, "    %s\n", ex)
	}
	return strings.TrimRight(sb.String(), "\n"), nil
}

// nationInfo holds static data about one classical power.
type nationInfo struct {
	name       string
	abbrev     string
	homeSCs    []string // province codes (lowercase), sorted
	startUnits []string // e.g. "F Edinburgh", "A Liverpool"
}

// abbrevToNation maps 3-letter abbreviations (lowercase) to full nation names.
var abbrevToNation = map[string]string{
	"eng": "England",
	"fra": "France",
	"ger": "Germany",
	"ita": "Italy",
	"aus": "Austria",
	"rus": "Russia",
	"tur": "Turkey",
}

// buildNationInfo constructs the nationInfo table from godip classical data.
// This is called once and cached.
func buildNationInfo() []nationInfo {
	scs := start.SupplyCenters()
	units := start.Units()
	longNames := classical.ClassicalVariant.ProvinceLongNames

	// Collect home SCs and starting units per nation.
	homeSCsMap := make(map[godip.Nation][]string)
	for prov, nation := range scs {
		homeSCsMap[nation] = append(homeSCsMap[nation], string(prov))
	}

	startUnitsMap := make(map[godip.Nation][]string)
	for prov, unit := range units {
		utype := "A"
		if unit.Type == godip.Fleet {
			utype = "F"
		}
		// Use the base province for the long name (strip /sc, /nc suffixes).
		baseProv := strings.Split(string(prov), "/")[0]
		longName := longNames[godip.Province(baseProv)]
		startUnitsMap[unit.Nation] = append(startUnitsMap[unit.Nation], utype+" "+longName)
	}

	orderedNations := []string{"England", "France", "Germany", "Italy", "Austria", "Russia", "Turkey"}
	abbrevs := map[string]string{
		"England": "Eng", "France": "Fra", "Germany": "Ger",
		"Italy": "Ita", "Austria": "Aus", "Russia": "Rus", "Turkey": "Tur",
	}

	result := make([]nationInfo, 0, len(orderedNations))
	for _, name := range orderedNations {
		n := godip.Nation(name)
		scsForNation := homeSCsMap[n]
		sort.Strings(scsForNation)
		unitsForNation := startUnitsMap[n]
		sort.Strings(unitsForNation)
		result = append(result, nationInfo{
			name:       name,
			abbrev:     abbrevs[name],
			homeSCs:    scsForNation,
			startUnits: unitsForNation,
		})
	}
	return result
}

// nationInfoTable is the lazily initialised nation data table.
var nationInfoTable []nationInfo

// getNationInfoTable returns the singleton nation info table, building it once.
func getNationInfoTable() []nationInfo {
	if nationInfoTable == nil {
		nationInfoTable = buildNationInfo()
	}
	return nationInfoTable
}

// resolveNation returns the full nation name for a given input (full name or
// abbreviation, case-insensitive). Returns empty string if not found.
func resolveNation(input string) string {
	// Try full name (title-cased).
	titled := strings.Title(strings.ToLower(input)) //nolint:staticcheck
	if classicalNations[titled] {
		return titled
	}
	// Try abbreviation.
	abbrev := strings.ToLower(input)
	if full, ok := abbrevToNation[abbrev]; ok {
		return full
	}
	return ""
}

// handleNations processes /nations [nation] — lists all classical powers or
// shows detail for one nation.
func (d *Dispatcher) handleNations(cmd Command) (string, error) {
	table := getNationInfoTable()
	longNames := classical.ClassicalVariant.ProvinceLongNames

	if len(cmd.Args) == 0 {
		// Table of all nations.
		var sb strings.Builder
		fmt.Fprintf(&sb, "%-10s %-7s %s\n", "Nation", "Abbrev", "Home SCs")
		for _, ni := range table {
			// Format home SCs as "Edinburgh (edi), London (lon), ..."
			scParts := make([]string, 0, len(ni.homeSCs))
			for _, sc := range ni.homeSCs {
				longName := longNames[godip.Province(sc)]
				scParts = append(scParts, fmt.Sprintf("%s (%s)", longName, sc))
			}
			fmt.Fprintf(&sb, "%-10s %-7s %s\n", ni.name, ni.abbrev, strings.Join(scParts, ", "))
		}
		return strings.TrimRight(sb.String(), "\n"), nil
	}

	// Detail for one nation.
	nationName := resolveNation(strings.Join(cmd.Args, " "))
	if nationName == "" {
		return "", fmt.Errorf("bot: unknown nation %q; valid names: Austria (Aus), England (Eng), France (Fra), Germany (Ger), Italy (Ita), Russia (Rus), Turkey (Tur)", strings.Join(cmd.Args, " "))
	}

	var ni nationInfo
	for _, n := range table {
		if n.name == nationName {
			ni = n
			break
		}
	}

	// Format home SCs.
	scParts := make([]string, 0, len(ni.homeSCs))
	for _, sc := range ni.homeSCs {
		longName := longNames[godip.Province(sc)]
		scParts = append(scParts, fmt.Sprintf("%s (%s)", longName, sc))
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "%s (%s)\n", ni.name, ni.abbrev)
	fmt.Fprintf(&sb, "  Home SCs:       %s\n", strings.Join(scParts, ", "))
	fmt.Fprintf(&sb, "  Starting units: %s\n", strings.Join(ni.startUnits, ", "))
	fmt.Fprintf(&sb, "  Win condition:  Control 18 of 34 supply centres")
	return sb.String(), nil
}

// handleProvinces processes /provinces [nation] — lists all province codes with
// full names, or filters to a nation's home SCs.
func (d *Dispatcher) handleProvinces(cmd Command) (string, error) {
	longNames := classical.ClassicalVariant.ProvinceLongNames

	if len(cmd.Args) == 0 {
		// Full alphabetical list.
		codes := make([]string, 0, len(longNames))
		for prov := range longNames {
			codes = append(codes, string(prov))
		}
		sort.Strings(codes)

		var sb strings.Builder
		fmt.Fprintln(&sb, "Province reference (Classical Diplomacy):")
		for _, code := range codes {
			fmt.Fprintf(&sb, "  %-7s — %s\n", code, longNames[godip.Province(code)])
		}
		return strings.TrimRight(sb.String(), "\n"), nil
	}

	// Filter to home SCs for the given nation.
	nationName := resolveNation(strings.Join(cmd.Args, " "))
	if nationName == "" {
		return "", fmt.Errorf("bot: unknown nation %q; valid names: Austria (Aus), England (Eng), France (Fra), Germany (Ger), Italy (Ita), Russia (Rus), Turkey (Tur)", strings.Join(cmd.Args, " "))
	}

	scs := start.SupplyCenters()
	var codes []string
	for prov, nation := range scs {
		if string(nation) == nationName {
			codes = append(codes, string(prov))
		}
	}
	sort.Strings(codes)

	var sb strings.Builder
	fmt.Fprintf(&sb, "%s home provinces:\n", nationName)
	for _, code := range codes {
		fmt.Fprintf(&sb, "  %-7s — %s\n", code, longNames[godip.Province(code)])
	}
	return strings.TrimRight(sb.String(), "\n"), nil
}

// handleDraw processes /draw — proposes a draw or casts a yes vote on an
// active draw proposal. When all remaining nations have voted yes, posts
// GameEnded with result="draw".
func (d *Dispatcher) handleDraw(cmd Command) (string, error) {
	state, err := d.readState(cmd.ChannelID)
	if err != nil {
		return "", err
	}
	if !state.started || state.ended {
		return "", fmt.Errorf("bot: no active game in this channel")
	}
	nation, ok := state.players[cmd.UserID]
	if !ok {
		return "", fmt.Errorf("bot: you are not a player in this game")
	}

	if !state.drawProposed {
		// First call: post the proposal.
		if err := events.Write(d.ch, cmd.ChannelID, events.TypeDrawProposed, events.DrawProposed{
			ProposerNation: nation,
		}); err != nil {
			return "", fmt.Errorf("bot: write DrawProposed: %w", err)
		}
		// If this is the only remaining nation the draw resolves immediately.
		if len(state.nations) == 1 {
			var finalState json.RawMessage
			if sess := d.sessions[cmd.ChannelID]; sess != nil {
				finalState, _ = sess.Eng.Dump()
			}
			if err := events.Write(d.ch, cmd.ChannelID, events.TypeGameEnded, events.GameEnded{
				Result: "draw", FinalState: finalState,
			}); err != nil {
				return "", fmt.Errorf("bot: write GameEnded: %w", err)
			}
			return "Draw agreed. Game over!", nil
		}
		return fmt.Sprintf("Draw proposed by %s. All nations must use /draw to accept.", nation), nil
	}

	// Draw already active — cast a yes vote.
	if state.drawVotes[nation] {
		return "You have already voted for this draw.", nil
	}
	if err := events.Write(d.ch, cmd.ChannelID, events.TypeDrawVoted, events.DrawVoted{
		Nation: nation, Accept: true,
	}); err != nil {
		return "", fmt.Errorf("bot: write DrawVoted: %w", err)
	}
	state.drawVotes[nation] = true

	// Check whether all remaining nations (by nation name) have now voted yes.
	allVoted := true
	for n := range state.nations {
		if !state.drawVotes[n] {
			allVoted = false
			break
		}
	}
	if allVoted {
		var finalState json.RawMessage
		if sess := d.sessions[cmd.ChannelID]; sess != nil {
			finalState, _ = sess.Eng.Dump()
		}
		if err := events.Write(d.ch, cmd.ChannelID, events.TypeGameEnded, events.GameEnded{
			Result: "draw", FinalState: finalState,
		}); err != nil {
			return "", fmt.Errorf("bot: write GameEnded: %w", err)
		}
		return "All nations agree. The game ends in a draw!", nil
	}

	remaining := len(state.nations) - len(state.drawVotes)
	return fmt.Sprintf("%s votes yes for the draw. Waiting for %d more nation(s).", nation, remaining), nil
}

// handleConcede processes /concede — the caller surrenders, ending the game
// immediately with result="concession".
func (d *Dispatcher) handleConcede(cmd Command) (string, error) {
	state, err := d.readState(cmd.ChannelID)
	if err != nil {
		return "", err
	}
	if !state.started || state.ended {
		return "", fmt.Errorf("bot: no active game in this channel")
	}
	nation, ok := state.players[cmd.UserID]
	if !ok {
		return "", fmt.Errorf("bot: you are not a player in this game")
	}

	var finalState json.RawMessage
	if sess := d.sessions[cmd.ChannelID]; sess != nil {
		finalState, _ = sess.Eng.Dump()
	}
	if err := events.Write(d.ch, cmd.ChannelID, events.TypeGameEnded, events.GameEnded{
		Result: "concession", Winner: nation, FinalState: finalState,
	}); err != nil {
		return "", fmt.Errorf("bot: write GameEnded: %w", err)
	}
	return fmt.Sprintf("%s concedes. Game over!", nation), nil
}

// handlePause processes /pause (GM only) — cancels the deadline timer.
func (d *Dispatcher) handlePause(cmd Command) (string, error) {
	sess, ok := d.sessions[cmd.ChannelID]
	if !ok || sess == nil {
		return "", fmt.Errorf("bot: no active game found in this channel")
	}
	if cmd.UserID != sess.GMID {
		return "", fmt.Errorf("bot: only the GM can pause the game")
	}
	sess.CancelDeadline()
	return "Game paused. Use /resume to restart the deadline.", nil
}

// handleResume processes /resume (GM only) — restarts a paused deadline timer.
func (d *Dispatcher) handleResume(cmd Command) (string, error) {
	sess, ok := d.sessions[cmd.ChannelID]
	if !ok || sess == nil {
		return "", fmt.Errorf("bot: no active game found in this channel")
	}
	if cmd.UserID != sess.GMID {
		return "", fmt.Errorf("bot: only the GM can resume the game")
	}
	sess.RestartDeadline()
	return "Game resumed. Deadline restarted.", nil
}

// handleExtend processes /extend <duration> (GM only) — adds time to the
// current phase deadline. Duration uses Go's time.ParseDuration format (e.g. "2h", "30m").
func (d *Dispatcher) handleExtend(cmd Command) (string, error) {
	sess, ok := d.sessions[cmd.ChannelID]
	if !ok || sess == nil {
		return "", fmt.Errorf("bot: no active game found in this channel")
	}
	if cmd.UserID != sess.GMID {
		return "", fmt.Errorf("bot: only the GM can extend the deadline")
	}
	if len(cmd.Args) == 0 {
		return "", fmt.Errorf("bot: usage: /extend <duration>")
	}
	dur, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return "", fmt.Errorf("bot: invalid duration %q: %w", cmd.Args[0], err)
	}
	sess.ExtendDeadline(dur)
	return fmt.Sprintf("Deadline extended by %s.", dur), nil
}

// handleForceResolve processes /force-resolve (GM only) — triggers AdvanceTurn
// immediately without waiting for the deadline.
func (d *Dispatcher) handleForceResolve(cmd Command) (string, error) {
	sess, ok := d.sessions[cmd.ChannelID]
	if !ok || sess == nil {
		return "", fmt.Errorf("bot: no active game found in this channel")
	}
	if cmd.UserID != sess.GMID {
		return "", fmt.Errorf("bot: only the GM can force-resolve the current phase")
	}
	if err := sess.AdvanceTurn(); err != nil {
		return "", fmt.Errorf("bot: force-resolve: %w", err)
	}
	return "Phase force-resolved.", nil
}

// handleBoot processes /boot <nation> (GM only) — removes a player from the
// game. Their units receive NMR orders each turn going forward.
func (d *Dispatcher) handleBoot(cmd Command) (string, error) {
	sess, ok := d.sessions[cmd.ChannelID]
	if !ok || sess == nil {
		return "", fmt.Errorf("bot: no active game found in this channel")
	}
	if cmd.UserID != sess.GMID {
		return "", fmt.Errorf("bot: only the GM can boot players")
	}
	if len(cmd.Args) == 0 {
		return "", fmt.Errorf("bot: usage: /boot <nation>")
	}
	nation := cmd.Args[0]

	// Find the userID for the given nation.
	var userID string
	for uid, n := range sess.Players {
		if n == nation {
			userID = uid
			break
		}
	}
	if userID == "" {
		return "", fmt.Errorf("bot: nation %q not found in this game", nation)
	}

	if err := events.Write(d.ch, cmd.ChannelID, events.TypePlayerBooted, events.PlayerBooted{
		Nation: nation,
	}); err != nil {
		return "", fmt.Errorf("bot: write PlayerBooted: %w", err)
	}
	delete(sess.Players, userID)
	return fmt.Sprintf("%s has been booted from the game.", nation), nil
}

// handleReplace processes /replace <nation> <user> (GM only) — transfers a
// nation to a new player identified by user.
func (d *Dispatcher) handleReplace(cmd Command) (string, error) {
	sess, ok := d.sessions[cmd.ChannelID]
	if !ok || sess == nil {
		return "", fmt.Errorf("bot: no active game found in this channel")
	}
	if cmd.UserID != sess.GMID {
		return "", fmt.Errorf("bot: only the GM can replace players")
	}
	if len(cmd.Args) < 2 {
		return "", fmt.Errorf("bot: usage: /replace <nation> <user>")
	}
	nation, newUserID := cmd.Args[0], cmd.Args[1]

	// Find the old userID for the given nation.
	var oldUserID string
	for uid, n := range sess.Players {
		if n == nation {
			oldUserID = uid
			break
		}
	}
	if oldUserID == "" {
		return "", fmt.Errorf("bot: nation %q not found in this game", nation)
	}

	if err := events.Write(d.ch, cmd.ChannelID, events.TypePlayerReplaced, events.PlayerReplaced{
		Nation: nation, NewUserID: newUserID,
	}); err != nil {
		return "", fmt.Errorf("bot: write PlayerReplaced: %w", err)
	}
	delete(sess.Players, oldUserID)
	sess.Players[newUserID] = nation
	return fmt.Sprintf("%s is now playing as %s.", newUserID, nation), nil
}

// allAdjustmentActionsSubmitted returns true when every nation that has a
// non-zero supply-centre delta (needs builds or disbands) has set
// sess.Submitted[nation] = true. Nations whose SC and unit counts are equal
// have no required orders and are not checked.
func allAdjustmentActionsSubmitted(sess *session.Session) bool {
	scCounts := sess.Eng.SupplyCenters() // map[nation]int
	unitCounts := map[string]int{}
	for _, u := range sess.Eng.Units() {
		unitCounts[u.Nation]++
	}
	for _, nation := range sess.Players {
		if scCounts[nation] != unitCounts[nation] && !sess.Submitted[nation] {
			return false
		}
	}
	return true
}

// allRetreatActionsSubmitted returns true when every nation that has at least
// one dislodged unit has set sess.Submitted[nation] = true.
func allRetreatActionsSubmitted(sess *session.Session) bool {
	nationsDislodged := map[string]bool{}
	for _, nation := range sess.Eng.Dislodgeds() {
		nationsDislodged[nation] = true
	}
	for nation := range nationsDislodged {
		if !sess.Submitted[nation] {
			return false
		}
	}
	return true
}

// allNationsSubmitted reads each player's DM thread to check whether every
// nation has an OrderSubmitted event for the current phase.
func (d *Dispatcher) allNationsSubmitted(sess *session.Session) (bool, error) {
	for userID, nation := range sess.Players {
		msgs, err := d.ch.DMHistory(userID)
		if err != nil {
			return false, fmt.Errorf("bot: dm history for %s: %w", nation, err)
		}
		found := false
		for _, msg := range msgs {
			var env events.Envelope
			if err := json.Unmarshal([]byte(msg), &env); err != nil {
				continue
			}
			if env.Type != events.TypeOrderSubmitted {
				continue
			}
			var os events.OrderSubmitted
			if err := json.Unmarshal(env.Payload, &os); err != nil {
				continue
			}
			if os.Phase == sess.Phase && os.Nation == nation {
				found = true
				break
			}
		}
		if !found {
			return false, nil
		}
	}
	return true, nil
}
