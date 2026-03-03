// Package bot contains the platform-agnostic command router. It receives
// parsed commands from platform adapters, enforces access control (only the
// assigned player may submit orders for their nation), and delegates to the
// session and engine packages.
package bot

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/burrbd/dip/dipmap"
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
	ch              events.Channel
	notifier        session.Notifier
	loader          session.EngineLoader
	newEng          EngineFactory
	sessions        map[string]*session.Session
	graph           dipmap.Graph                                               // optional board graph for /map neighbourhood queries
	renderFn        func(dipmap.EngineState) ([]byte, error)                  // defaults to dipmap.Render (full-board PNG)
	svgFn           func(dipmap.EngineState) ([]byte, error)                  // defaults to dipmap.LoadSVG (raw SVG bytes)
	highlightFn     func([]byte, []string) ([]byte, error)                    // defaults to dipmap.Highlight
	renderZoomedFn  func(dipmap.EngineState, []byte, []string) ([]byte, error) // defaults to dipmap.RenderZoomed
}

// New returns a Dispatcher wired to the given dependencies.
func New(ch events.Channel, notifier session.Notifier, loader session.EngineLoader, newEng EngineFactory) *Dispatcher {
	return &Dispatcher{
		ch:             ch,
		notifier:       notifier,
		loader:         loader,
		newEng:         newEng,
		sessions:       make(map[string]*session.Session),
		renderFn:       dipmap.Render,
		svgFn:          dipmap.LoadSVG,
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

// boardGraph returns the Dispatcher's graph or EmptyGraph if none is set.
func (d *Dispatcher) boardGraph() dipmap.Graph {
	if d.graph != nil {
		return d.graph
	}
	return dipmap.EmptyGraph{}
}

// handleMap processes /map [territory [n]] — renders the board and posts it
// as an image. With territory and n > 0, highlights the neighbourhood and
// returns a zoomed crop. With no territory or n == 0, returns the full board.
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

	var img []byte
	if territory != "" && n > 0 {
		// Zoomed path: load SVG → highlight neighbourhood → render zoomed PNG.
		provinces := dipmap.Neighborhood(d.boardGraph(), territory, n)
		svg, err := d.svgFn(sess.Eng)
		if err != nil {
			return "", fmt.Errorf("bot: render map: %w", err)
		}
		highlighted, err := d.highlightFn(svg, provinces)
		if err != nil {
			return "", fmt.Errorf("bot: render map: %w", err)
		}
		var err2 error
		img, err2 = d.renderZoomedFn(sess.Eng, highlighted, provinces)
		if err2 != nil {
			return "", fmt.Errorf("bot: render map: %w", err2)
		}
	} else {
		// Full-board path: render the whole map as PNG.
		var err error
		img, err = d.renderFn(sess.Eng)
		if err != nil {
			return "", fmt.Errorf("bot: render map: %w", err)
		}
	}

	if err := d.ch.PostImage(cmd.ChannelID, img); err != nil {
		return "", fmt.Errorf("bot: post map: %w", err)
	}
	return "Map posted.", nil
}

// commandHelp maps command names to their help text.
var commandHelp = map[string]string{
	"newgame":       "/newgame — Start a new game in this channel. You become the GM.",
	"join":          "/join <nation> — Join the game as a nation (e.g. England, France).",
	"start":         "/start — Start the game (GM only). Requires 2–7 players.",
	"order":         "/order <order-text> — Submit a movement order (DM only, Movement phase).",
	"orders":        "/orders — List your staged orders for the current phase (DM only).",
	"clear":         "/clear [order] — Clear all staged orders or a specific one (DM only).",
	"submit":        "/submit — Finalise and submit your orders (DM only, Movement phase).",
	"retreat":       "/retreat <unit_type> <source> <destination> — Retreat a dislodged unit (DM only, Retreat phase).",
	"disband":       "/disband <unit_type> <province> — Disband a unit (DM only, Retreat or Adjustment phase).",
	"build":         "/build <unit_type> <province> — Build a unit (DM only, Adjustment phase).",
	"waive":         "/waive — Waive a build slot (DM only, Adjustment phase).",
	"status":        "/status — Show current phase, SC counts, and submission status.",
	"history":       "/history <turn> — Show adjudication results for a past turn (e.g. /history Spring 1901).",
	"map":           "/map [territory [n]] — Post the board map. With args, highlights territory and all provinces within n hops.",
	"help":          "/help [command] — List all commands or show help for a specific command.",
	"draw":          "/draw — Propose a draw; subsequent /draw calls from other nations are yes votes.",
	"concede":       "/concede — Concede the game immediately.",
	"pause":         "/pause — Pause the phase deadline timer (GM only).",
	"resume":        "/resume — Resume a paused deadline timer (GM only).",
	"extend":        "/extend <duration> — Extend the current deadline (GM only, e.g. /extend 2h).",
	"force-resolve": "/force-resolve — Resolve the current phase immediately (GM only).",
	"boot":          "/boot <nation> — Remove a player from the game (GM only).",
	"replace":       "/replace <nation> <user> — Transfer a nation to a new player (GM only).",
}

// commandList defines the canonical display order for /help.
var commandList = []string{
	"newgame", "join", "start",
	"order", "orders", "clear", "submit",
	"retreat", "disband", "build", "waive",
	"status", "history", "map", "help",
	"draw", "concede",
	"pause", "resume", "extend", "force-resolve", "boot", "replace",
}

// handleHelp processes /help [command] — lists all commands or shows detailed
// usage for a specific command.
func (d *Dispatcher) handleHelp(cmd Command) (string, error) {
	if len(cmd.Args) == 0 {
		var sb strings.Builder
		fmt.Fprintln(&sb, "Available commands:")
		for _, name := range commandList {
			fmt.Fprintf(&sb, "  %s\n", commandHelp[name])
		}
		return strings.TrimRight(sb.String(), "\n"), nil
	}

	name := cmd.Args[0]
	if strings.HasPrefix(name, "/") {
		name = name[1:]
	}
	if text, ok := commandHelp[name]; ok {
		return text, nil
	}
	return "", fmt.Errorf("bot: unknown command %q; use /help for a list", cmd.Args[0])
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
