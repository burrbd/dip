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
