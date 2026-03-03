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

	mu         sync.Mutex
	ch         events.Channel
	notifier   Notifier
	timer      *time.Timer
	deadlineAt time.Time // absolute UTC time when the current phase deadline fires
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

// RestartDeadline re-enables a paused deadline timer. It restarts from the
// remaining time in s.deadlineAt. If s.deadlineAt is in the past, a fresh
// DeadlineHours-duration timer is started. No-op when DeadlineHours is 0 and
// deadlineAt is unset.
func (s *Session) RestartDeadline() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.timer != nil {
		s.timer.Stop()
		s.timer = nil
	}
	remaining := time.Until(s.deadlineAt)
	if remaining <= 0 {
		if s.DeadlineHours <= 0 {
			return
		}
		remaining = time.Duration(s.DeadlineHours) * time.Hour
		s.deadlineAt = time.Now().Add(remaining)
	}
	s.timer = time.AfterFunc(remaining, s.onDeadline)
}

// ExtendDeadline adds d to the current deadline and resets the timer. If
// s.deadlineAt is unset but DeadlineHours > 0, the current time plus
// DeadlineHours is used as the base before extending. No-op when DeadlineHours
// is 0 and deadlineAt is unset.
func (s *Session) ExtendDeadline(d time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.timer != nil {
		s.timer.Stop()
		s.timer = nil
	}
	if s.deadlineAt.IsZero() {
		if s.DeadlineHours <= 0 {
			return
		}
		s.deadlineAt = time.Now().Add(time.Duration(s.DeadlineHours) * time.Hour)
	}
	s.deadlineAt = s.deadlineAt.Add(d)
	if remaining := time.Until(s.deadlineAt); remaining > 0 {
		s.timer = time.AfterFunc(remaining, s.onDeadline)
	}
}
