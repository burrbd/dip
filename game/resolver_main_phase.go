package game

import (
	"sort"

	"github.com/burrbd/dip/game/order"
	"github.com/burrbd/dip/game/order/board"
)

// ResolveOrders resolves all orders for a turn using a simultaneous
// fixed-point algorithm.
//
// Each outer iteration:
//  1. Stabilise support strengths: recompute move strengths based on which
//     units are tentative winners at supporter territories, iterate until
//     no strength changes (inner fixed-point).
//  2. Collect all conflict groups simultaneously and apply every outcome at
//     once (bounces / defeats), so that no single conflict resolution can
//     corrupt the inputs to another.
//  3. Repeat until there are no conflicts left.
func ResolveOrders(manager board.Manager, orders order.Set) {
	for {
		stabilizeStrengths(manager, orders)
		if !resolvePass(manager) {
			return
		}
	}
}

// stabilizeStrengths iterates the inner fixed-point: recompute all move
// strengths based on tentative winners, stop when nothing changes.
func stabilizeStrengths(manager board.Manager, orders order.Set) {
	for i := 0; i < 20; i++ {
		prev := captureStrengths(manager)
		refreshStrengths(manager, orders)
		if !strengthsChanged(manager, prev) {
			return
		}
	}
}

// refreshStrengths rewrites the strength of every unit that is currently
// moving, using winner-based support-cut detection.
func refreshStrengths(manager board.Manager, orders order.Set) {
	winners := tentativeWinners(manager)
	for unit, pos := range manager.Positions() {
		if pos.Cause != board.Moved {
			continue
		}
		for _, move := range orders.Moves {
			if move.From.Is(manager.Origin(unit)) && move.To.Is(pos.Territory) {
				manager.UpdateStrength(unit, countMoveSupports(move, orders.MoveSupports, manager, winners))
				break
			}
		}
	}
}

// tentativeWinners returns a map from territory abbreviation to the unit that
// is the tentative winner there (nil key means a tie — no winner).
//
// A unit is the tentative winner at a territory when it has strictly greater
// strength than every other unit in that conflict group and it moved there
// (Cause == Moved). Holding units do not "win" in the attacking sense and
// cannot cut support on their own territory.
//
// Uncontested movers (no conflict at their destination) are also tentative
// winners, since nothing is stopping them.
func tentativeWinners(manager board.Manager) map[string]*board.Unit {
	conflictedTerrs := make(map[string]bool)
	winners := make(map[string]*board.Unit)

	for _, group := range manager.AllConflicts() {
		for _, u := range group {
			conflictedTerrs[manager.Position(u).Territory.Abbr] = true
		}
		sort.Sort(UnitPositionsByStrength(manager, group))
		if manager.Position(group[0]).Strength > manager.Position(group[1]).Strength {
			w := group[0]
			if manager.Position(w).Cause == board.Moved {
				winners[manager.Position(w).Territory.Abbr] = w
			}
		}
		// tie — no winner recorded for this group
	}

	// Uncontested movers are also tentative winners at their destination.
	for u, pos := range manager.Positions() {
		if pos.Cause == board.Moved && !conflictedTerrs[pos.Territory.Abbr] {
			winners[pos.Territory.Abbr] = u
		}
	}

	return winners
}

// countMoveSupports counts how many MoveSupport orders for this move are not
// cut, using winner-based cut detection.
func countMoveSupports(move order.Move, supports []order.MoveSupport, manager board.Manager, winners map[string]*board.Unit) int {
	count := 0
	for _, sup := range supports {
		if sup.Move.From.Abbr == move.From.Abbr && sup.Move.To.Abbr == move.To.Abbr {
			if !isSupportCutByWinner(sup, manager, winners) {
				count++
			}
		}
	}
	return count
}

// isSupportCutByWinner reports whether a MoveSupport is cut given the current
// set of tentative winners. Support is cut only when a unit that WINS (not
// merely attacks) the supporter's territory originated outside the support
// target — a unit that only ties at the territory is not a winner and
// therefore cannot cut the support.
func isSupportCutByWinner(sup order.MoveSupport, manager board.Manager, winners map[string]*board.Unit) bool {
	winner, ok := winners[sup.By.Abbr]
	if !ok || winner == nil {
		return false
	}
	pos := manager.Position(winner)
	if pos.Cause != board.Moved {
		// The unit at the supporter's territory is holding there, not attacking.
		return false
	}
	return manager.Origin(winner).IsNot(sup.Move.To)
}

// resolvePass applies all current conflict outcomes simultaneously.
// Returns true if any unit was bounced or defeated.
func resolvePass(manager board.Manager) bool {
	groups := manager.AllConflicts()
	if len(groups) == 0 {
		return false
	}

	type outcome struct{ defeated bool }
	pending := make(map[*board.Unit]outcome)

	for _, units := range groups {
		sort.Sort(UnitPositionsByStrength(manager, units))
		decisive := manager.Position(units[0]).Strength > manager.Position(units[1]).Strength
		losers := units
		if decisive {
			losers = units[1:]
		}
		for _, u := range losers {
			if _, exists := pending[u]; exists {
				continue // already scheduled from another conflict group
			}
			atOrigin := manager.AtOrigin(u)
			if atOrigin && decisive {
				pending[u] = outcome{true}
			} else if !atOrigin {
				pending[u] = outcome{false}
			}
		}
	}

	if len(pending) == 0 {
		return false
	}
	for u, o := range pending {
		if o.defeated {
			manager.SetDefeated(u)
		} else {
			manager.Bounce(u)
		}
	}
	return true
}

func captureStrengths(manager board.Manager) map[*board.Unit]int {
	prev := make(map[*board.Unit]int)
	for u, pos := range manager.Positions() {
		prev[u] = pos.Strength
	}
	return prev
}

func strengthsChanged(manager board.Manager, prev map[*board.Unit]int) bool {
	for u, pos := range manager.Positions() {
		if pos.Strength != prev[u] {
			return true
		}
	}
	return false
}
