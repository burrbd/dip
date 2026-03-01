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

// CancelDeadline stops any pending deadline timer without firing it.
func (s *Session) CancelDeadline() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.timer != nil {
		s.timer.Stop()
		s.timer = nil
	}
}
