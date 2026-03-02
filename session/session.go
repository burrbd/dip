// Package session manages the lifecycle of a single Diplomacy game session
// within a chat channel. It owns the deadline timer, staged orders, player→nation
// mapping, GM identity, and phase transitions.
package session

import (
	"sync"
	"time"

	"github.com/burrbd/dip/engine"
	"github.com/burrbd/dip/events"
)

// Notifier sends messages to game participants.
type Notifier interface {
	Notify(channelID, message string) error
}

// Session represents the state of a single Diplomacy game within a channel.
type Session struct {
	ChannelID     string
	Phase         string
	StagedOrders  map[string][]string // nation → staged order texts for the current phase
	Players       map[string]string   // userID → nation
	Submitted     map[string]bool     // nation → true if orders are finalised
	GMID          string
	DeadlineHours int
	Eng           engine.Engine

	mu       sync.Mutex
	ch       events.Channel
	notifier Notifier
	timer    *time.Timer
}

// New creates a Session with all required dependencies and starts the deadline
// timer. It is called by the bot after posting a GameStarted event to wire up
// the in-process deadline manager for the new game.
func New(ch events.Channel, channelID, gmID, phase string, players map[string]string, deadlineHours int, eng engine.Engine, notifier Notifier) *Session {
	s := &Session{
		ChannelID:     channelID,
		Phase:         phase,
		StagedOrders:  make(map[string][]string),
		Players:       make(map[string]string),
		Submitted:     make(map[string]bool),
		GMID:          gmID,
		DeadlineHours: deadlineHours,
		Eng:           eng,
		ch:            ch,
		notifier:      notifier,
	}
	for k, v := range players {
		s.Players[k] = v
	}
	s.startDeadline()
	return s
}

// CancelDeadline stops any pending deadline timer without firing it.
func (s *Session) CancelDeadline() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.timer != nil {
		s.timer.Stop()
		s.timer = nil
	}
}
