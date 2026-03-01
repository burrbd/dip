package session

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/burrbd/dip/events"
)

// AdvanceTurn adjudicates the current phase and advances the game to the next.
//
// It runs: cancel existing timer → resolve staged orders → post PhaseResolved
// event → notify players → check for solo winner → advance phase → reset staged
// orders → start new deadline timer.
func (s *Session) AdvanceTurn() error {
	s.CancelDeadline()

	result, err := s.Eng.Resolve()
	if err != nil {
		return fmt.Errorf("session: resolve: %w", err)
	}

	snapshot, err := s.Eng.Dump()
	if err != nil {
		return fmt.Errorf("session: dump state: %w", err)
	}

	summary, _ := json.Marshal(result)
	if err := events.Write(s.ch, s.ChannelID, events.TypePhaseResolved, events.PhaseResolved{
		Phase:         result.Phase,
		StateSnapshot: snapshot,
		ResultSummary: summary,
	}); err != nil {
		return fmt.Errorf("session: write PhaseResolved: %w", err)
	}

	if s.notifier != nil {
		msg := fmt.Sprintf("Phase %s resolved. %d orders adjudicated.", result.Phase, len(result.Orders))
		_ = s.notifier.Notify(s.ChannelID, msg)
	}

	// Check for solo winner before advancing to the next phase.
	if winner := s.Eng.SoloWinner(); winner != "" {
		finalState, _ := s.Eng.Dump()
		return events.Write(s.ch, s.ChannelID, events.TypeGameEnded, events.GameEnded{
			Result:     "solo",
			Winner:     winner,
			FinalState: finalState,
		})
	}

	if err := s.Eng.Advance(); err != nil {
		return fmt.Errorf("session: advance: %w", err)
	}

	s.StagedOrders = make(map[string][]string)
	s.Submitted = make(map[string]bool)
	s.Phase = s.Eng.Phase()

	s.startDeadline()
	return nil
}

// startDeadline starts the deadline timer using s.DeadlineHours. When it
// fires, onDeadline is called automatically. Does nothing if DeadlineHours ≤ 0.
func (s *Session) startDeadline() {
	if s.DeadlineHours <= 0 {
		return
	}
	d := time.Duration(s.DeadlineHours) * time.Hour
	s.mu.Lock()
	s.timer = time.AfterFunc(d, s.onDeadline)
	s.mu.Unlock()
}

// onDeadline is the timer callback invoked when the phase deadline expires.
func (s *Session) onDeadline() {
	_ = s.AdvanceTurn()
}
