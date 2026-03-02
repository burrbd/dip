// Package bot contains the platform-agnostic command router. It receives
// parsed commands from platform adapters, enforces access control (only the
// assigned player may submit orders for their nation), and delegates to the
// session and engine packages.
package bot

import (
	"encoding/json"
	"fmt"

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
	Name      string   // command name without leading slash (e.g. "newgame")
	Args      []string // positional arguments
	UserID    string   // platform user identifier
	ChannelID string   // platform channel identifier
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
