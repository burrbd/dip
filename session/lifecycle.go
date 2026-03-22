package session

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/burrbd/dip/engine"
	"github.com/burrbd/dip/events"
)

// AdvanceTurn adjudicates the current phase and advances the game to the next.
//
// It runs: cancel existing timer → resolve staged orders → post PhaseResolved
// event → post order summary → notify players → check for solo winner →
// advance phase → post phase guidance → reset staged orders → start new
// deadline timer.
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

	// Post human-readable order summary (Story 14).
	nations := s.nationList()
	_ = s.ch.Post(s.ChannelID, formatOrderSummary(result, nations))

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

	// Post phase guidance for the new phase (Story 16).
	_ = s.ch.Post(s.ChannelID, s.formatPhaseGuidance(nations))

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
	s.deadlineAt = time.Now().Add(d)
	s.timer = time.AfterFunc(d, s.onDeadline)
	s.mu.Unlock()
}

// onDeadline is the timer callback invoked when the phase deadline expires.
func (s *Session) onDeadline() {
	_ = s.AdvanceTurn()
}

// nationList returns an alphabetically sorted slice of all active nations in
// the session (derived from the Players map).
func (s *Session) nationList() []string {
	seen := make(map[string]bool)
	for _, nation := range s.Players {
		seen[nation] = true
	}
	nations := make([]string, 0, len(seen))
	for n := range seen {
		nations = append(nations, n)
	}
	sort.Strings(nations)
	return nations
}

// formatOrderSummary builds the human-readable adjudication summary posted
// after each phase resolution (Story 14).
func formatOrderSummary(result engine.ResolutionResult, nations []string) string {
	// Group orders by nation.
	byNation := make(map[string][]engine.OrderResult)
	for _, o := range result.Orders {
		byNation[o.Nation] = append(byNation[o.Nation], o)
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "%s %d %s — orders resolved:\n", result.Season, result.Year, result.Phase)

	for _, nation := range nations {
		ords, ok := byNation[nation]
		if !ok {
			// Nation has no units — skip.
			continue
		}
		// Check whether all orders are NMR.
		allNMR := true
		for _, o := range ords {
			if !o.IsNMR {
				allNMR = false
				break
			}
		}
		if allNMR {
			fmt.Fprintf(&sb, "\n%s: no orders submitted\n", nation)
		} else {
			fmt.Fprintf(&sb, "\n%s:\n", nation)
		}
		for _, o := range ords {
			if o.IsNMR {
				fmt.Fprintf(&sb, "  %s: hold (auto)\n", o.Order)
			} else {
				fmt.Fprintf(&sb, "  %s: %s\n", o.Order, o.Outcome)
			}
		}
	}
	return strings.TrimRight(sb.String(), "\n")
}

// formatPhaseGuidance builds the channel message that tells players what they
// need to do in the new phase (Story 16).
func (s *Session) formatPhaseGuidance(nations []string) string {
	phase := s.Phase
	eng := s.Eng

	if strings.HasSuffix(phase, "Movement") {
		return formatMovementGuidance(phase, eng, nations)
	}
	if strings.HasSuffix(phase, "Retreat") {
		return formatRetreatGuidance(phase, eng, nations)
	}
	if strings.HasSuffix(phase, "Adjustment") {
		return formatAdjustmentGuidance(phase, eng, nations)
	}
	return fmt.Sprintf("%s — new phase begun.", phase)
}

// formatMovementGuidance lists each nation's units for a Movement phase.
func formatMovementGuidance(phase string, eng engine.Engine, nations []string) string {
	// Group units by nation.
	byNation := make(map[string][]string)
	for prov, u := range eng.Units() {
		byNation[u.Nation] = append(byNation[u.Nation], prov)
	}
	// Sort province lists for deterministic output.
	for n := range byNation {
		sort.Strings(byNation[n])
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "%s — submit orders for your units.\n", phase)
	for _, nation := range nations {
		provs, ok := byNation[nation]
		if !ok || len(provs) == 0 {
			fmt.Fprintf(&sb, "\n%s: nothing to do this phase\n", nation)
		} else {
			fmt.Fprintf(&sb, "\n%s: %s\n", nation, strings.Join(provs, ", "))
		}
	}
	return strings.TrimRight(sb.String(), "\n")
}

// formatRetreatGuidance lists dislodged units and their valid retreat options
// for a Retreat phase.
func formatRetreatGuidance(phase string, eng engine.Engine, nations []string) string {
	retreats := eng.ValidRetreats() // prov → []dest
	dislodgeds := eng.Dislodgeds()  // prov → nation

	// Group dislodged provinces by nation.
	byNation := make(map[string][]string)
	for prov, nation := range dislodgeds {
		byNation[nation] = append(byNation[nation], prov)
	}
	for n := range byNation {
		sort.Strings(byNation[n])
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "%s — the following units must retreat or disband:\n", phase)
	for _, nation := range nations {
		provs, ok := byNation[nation]
		if !ok || len(provs) == 0 {
			fmt.Fprintf(&sb, "\n%s: nothing to do this phase\n", nation)
			continue
		}
		fmt.Fprintf(&sb, "\n%s:\n", nation)
		for _, prov := range provs {
			dests := retreats[prov]
			if len(dests) == 0 {
				fmt.Fprintf(&sb, "  %s: must disband (no valid retreats)\n", prov)
			} else {
				sort.Strings(dests)
				fmt.Fprintf(&sb, "  %s: retreat to %s — or disband\n", prov, strings.Join(dests, ", "))
			}
		}
	}
	return strings.TrimRight(sb.String(), "\n")
}

// formatAdjustmentGuidance lists build/disband requirements for each nation
// for an Adjustment phase.
func formatAdjustmentGuidance(phase string, eng engine.Engine, nations []string) string {
	opts := eng.BuildOptions() // nation → BuildOption

	var sb strings.Builder
	fmt.Fprintf(&sb, "%s:\n", phase)
	for _, nation := range nations {
		opt, ok := opts[nation]
		if !ok || opt.Delta == 0 {
			fmt.Fprintf(&sb, "\n%s: nothing to do this phase\n", nation)
			continue
		}
		if opt.Delta > 0 {
			homes := make([]string, len(opt.AvailableHomes))
			copy(homes, opt.AvailableHomes)
			sort.Strings(homes)
			fmt.Fprintf(&sb, "\n%s: build %d unit(s) — available home centres: %s\n",
				nation, opt.Delta, strings.Join(homes, ", "))
		} else {
			fmt.Fprintf(&sb, "\n%s: disband %d unit(s)\n", nation, -opt.Delta)
		}
	}
	return strings.TrimRight(sb.String(), "\n")
}

