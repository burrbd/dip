package game

import (
	"sort"

	"github.com/burrbd/dip/game/order"
	"github.com/burrbd/dip/game/order/board"
)

// ResolveOrders resolves orders for a turn.
func ResolveOrders(manager board.Manager, orders order.Set) {
	for {
		recomputeStrengths(manager, orders)
		units := manager.Conflict()
		if units == nil {
			return
		}
		var defeated bool
		sort.Sort(UnitPositionsByStrength(manager, units))
		if manager.Position(units[0]).Strength > manager.Position(units[1]).Strength {
			units = units[1:]
			defeated = true
		}
		for _, u := range units {
			atOrigin := manager.AtOrigin(u)
			if atOrigin && defeated {
				manager.SetDefeated(u)
			} else if !atOrigin {
				manager.Bounce(u)
			}
		}
	}
}

// recomputeStrengths recalculates the support strength of every unit that is
// still actively moving, based on who is currently at the supporter's territory.
// A support is considered cut only when a unit is physically present at the
// supporter's territory (Cause == Moved) and did not originate from the support
// target — units that have bounced back to their origin are no longer present
// there and therefore no longer cut the support.
func recomputeStrengths(manager board.Manager, orders order.Set) {
	for unit, pos := range manager.Positions() {
		if pos.Cause != board.Moved {
			continue
		}
		for _, move := range orders.Moves {
			if move.From.Is(manager.Origin(unit)) && move.To.Is(pos.Territory) {
				manager.UpdateStrength(unit, moveSupportStrength(move, orders.MoveSupports, manager))
				break
			}
		}
	}
}

func moveSupportStrength(move order.Move, supports []order.MoveSupport, manager board.Manager) int {
	strength := 0
	for _, sup := range supports {
		if sup.Move.From.Abbr == move.From.Abbr && sup.Move.To.Abbr == move.To.Abbr {
			if !moveSupportCutByPosition(sup, manager) {
				strength++
			}
		}
	}
	return strength
}

// moveSupportCutByPosition reports whether support is cut based on current
// board positions. Support is cut when a unit is currently moving (Cause ==
// Moved) into the supporter's territory from somewhere other than the support
// target (a unit cannot cut the support of a move directed at its own territory).
func moveSupportCutByPosition(sup order.MoveSupport, manager board.Manager) bool {
	for unit, pos := range manager.Positions() {
		if pos.Territory.Is(sup.By) && pos.Cause == board.Moved {
			if manager.Origin(unit).IsNot(sup.Move.To) {
				return true
			}
		}
	}
	return false
}
