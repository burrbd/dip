package game

import (
	"sort"

	"github.com/burrbd/dip/game/order/board"
)

// ResolveOrders resolves all orders for a turn using a simultaneous
// fixed-point algorithm.
//
// Support strengths are computed once by OrderHandler.ApplyOrders using
// order-based support-cut detection (DATC standard: any attack on a
// supporter's territory cuts its support, regardless of whether the attack
// succeeds). This function then loops, applying all conflict outcomes
// simultaneously each pass, until no conflicts remain.
func ResolveOrders(manager board.Manager) {
	for {
		if !resolvePass(manager) {
			return
		}
	}
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
