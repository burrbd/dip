// Package bot contains the platform-agnostic command router. It receives
// parsed commands from platform adapters, enforces access control (only the
// assigned player may submit orders for their nation), and delegates to the
// session and engine packages.
package bot

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/burrbd/dip/engine"
	"github.com/burrbd/dip/events"
	"github.com/burrbd/dip/session"
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
	ch       events.Channel
	notifier session.Notifier
	loader   session.EngineLoader
	newEng   EngineFactory
	sessions map[string]*session.Session
}

// New returns a Dispatcher wired to the given dependencies.
func New(ch events.Channel, notifier session.Notifier, loader session.EngineLoader, newEng EngineFactory) *Dispatcher {
	return &Dispatcher{
		ch:       ch,
		notifier: notifier,
		loader:   loader,
		newEng:   newEng,
		sessions: make(map[string]*session.Session),
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
	default:
		return "", fmt.Errorf("bot: unknown command %q", cmd.Name)
	}
}

// gameState is the bot's view of a channel's game state, derived from an event scan.
type gameState struct {
	created       bool
	started       bool
	gmID          string
	deadlineHours int
	players       map[string]string // userID → nation
	nations       map[string]string // nation → userID
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
	dislodgeds := sess.Eng.Dislodgeds()
	if n, exists := dislodgeds[src]; !exists || n != nation {
		return "", fmt.Errorf("bot: no dislodged %s unit at %s belonging to %s", unitType, src, nation)
	}
	orderText := fmt.Sprintf("%s %s R %s", unitType, src, dest)
	if err := sess.Eng.SubmitOrder(nation, orderText); err != nil {
		return "", fmt.Errorf("bot: invalid retreat order: %w", err)
	}
	sess.StagedOrders[nation] = append(sess.StagedOrders[nation], orderText)
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
	if isRetreatPhase(sess.Phase) {
		dislodgeds := sess.Eng.Dislodgeds()
		if n, exists := dislodgeds[src]; !exists || n != nation {
			return "", fmt.Errorf("bot: no dislodged %s unit at %s belonging to %s", unitType, src, nation)
		}
	}
	orderText := fmt.Sprintf("%s %s D", unitType, src)
	if err := sess.Eng.SubmitOrder(nation, orderText); err != nil {
		return "", fmt.Errorf("bot: invalid disband order: %w", err)
	}
	sess.StagedOrders[nation] = append(sess.StagedOrders[nation], orderText)
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
	orderText := fmt.Sprintf("%s %s B", unitType, province)
	if err := sess.Eng.SubmitOrder(nation, orderText); err != nil {
		return "", fmt.Errorf("bot: invalid build order: %w", err)
	}
	sess.StagedOrders[nation] = append(sess.StagedOrders[nation], orderText)
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
	return "Waive order staged.", nil
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
