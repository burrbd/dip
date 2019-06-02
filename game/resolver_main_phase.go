package game

import (
	"sort"

	"github.com/burrbd/dip/game/order/board"
)

// ResolveOrders resolves orders for a turn.
func ResolveOrders(manager board.Manager) {
	for {
		units := manager.Conflict()
		if units == nil {
			return
		}
		var mustLose bool
		sort.Sort(UnitPositionsByStrength(manager, units))
		if manager.Position(units[0]).Strength > manager.Position(units[1]).Strength {
			units, mustLose = units[1:], true
		}
		for _, u := range units {
			if !manager.AtOrigin(u) {
				manager.Bounce(u)
			} else if mustLose {
				manager.SetDefeated(u)
			}
		}
	}
}
