// Package events defines the structured event types that are serialised as JSON
// and posted to the chat channel after each game action. The channel's message
// history is the authoritative event log; state is restored by replaying events
// from the last PhaseResolved snapshot.
package events

import "encoding/json"

// EventType identifies the kind of event.
type EventType string

const (
	TypeGameCreated    EventType = "GameCreated"
	TypePlayerJoined   EventType = "PlayerJoined"
	TypeGameStarted    EventType = "GameStarted"
	TypeOrderSubmitted EventType = "OrderSubmitted"
	TypePhaseResolved  EventType = "PhaseResolved"
	TypePhaseSkipped   EventType = "PhaseSkipped"
	TypeNMRRecorded    EventType = "NMRRecorded"
	TypeDrawProposed   EventType = "DrawProposed"
	TypeDrawVoted      EventType = "DrawVoted"
	TypeGameEnded      EventType = "GameEnded"
)

// Envelope wraps a typed event payload for serialisation in the channel.
type Envelope struct {
	Type    EventType       `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// GameCreated is posted when a new game is initialised.
type GameCreated struct {
	Variant       string `json:"variant"`
	DeadlineHours int    `json:"deadline_hours"`
	GMUserID      string `json:"gm_user_id"`
}

// PlayerJoined is posted when a player claims a nation.
type PlayerJoined struct {
	UserID string `json:"user_id"`
	Nation string `json:"nation"`
}

// GameStarted is posted when the GM starts the game; it carries the initial
// godip state snapshot so the bot can restore state after a restart.
type GameStarted struct {
	InitialState json.RawMessage `json:"initial_state"`
}

// OrderSubmitted is posted each time a player submits one or more orders.
type OrderSubmitted struct {
	UserID string   `json:"user_id"`
	Nation string   `json:"nation"`
	Orders []string `json:"orders"`
	Phase  string   `json:"phase"`
}

// PhaseResolved is posted after adjudication; it carries the new godip state
// snapshot and a human-readable result summary.
type PhaseResolved struct {
	Phase         string          `json:"phase"`
	StateSnapshot json.RawMessage `json:"state_snapshot"`
	ResultSummary json.RawMessage `json:"result_summary,omitempty"`
}

// PhaseSkipped is posted when a phase is skipped automatically.
// Reason is one of "no_dislodgements" or "no_sc_delta".
type PhaseSkipped struct {
	Phase  string `json:"phase"`
	Reason string `json:"reason"`
}

// NMRRecorded is posted when a nation does not submit orders in time.
type NMRRecorded struct {
	Nation     string   `json:"nation"`
	Phase      string   `json:"phase"`
	AutoOrders []string `json:"auto_orders"`
}

// DrawProposed is posted when a nation proposes a draw.
type DrawProposed struct {
	ProposerNation string `json:"proposer_nation"`
}

// DrawVoted is posted when a nation votes on a pending draw proposal.
type DrawVoted struct {
	Nation string `json:"nation"`
	Accept bool   `json:"accept"`
}

// GameEnded is posted when the game concludes (solo win, draw, or concession).
// Result is one of "solo", "draw", or "concession".
type GameEnded struct {
	Result     string          `json:"result"`
	Winner     string          `json:"winner,omitempty"`
	FinalState json.RawMessage `json:"final_state,omitempty"`
}
