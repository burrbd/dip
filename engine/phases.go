package engine

import "github.com/zond/godip"

// Advance fills NMR orders for any unit without a staged order, advances the
// game to the next phase via godip Next(), and auto-skips empty retreat or
// adjustment phases (phases where no unit requires an order).
//
// When Resolve() has already called Next() (g.advanced == true), Advance()
// clears the flag and only handles empty-phase skipping.
func (g *game) Advance() error {
	if g.advanced {
		g.advanced = false
		// State was already advanced by Resolve(); only skip empty phases.
		for isEmptyPhase(g.adj) {
			fillNMR(g.adj)
			next, err := g.adj.Next()
			if err != nil {
				return err
			}
			g.adj = next
		}
		return nil
	}

	fillNMR(g.adj)

	next, err := g.adj.Next()
	if err != nil {
		return err
	}
	g.adj = next

	// Skip empty retreat phases (no dislodged units) and empty adjustment
	// phases (no builds/disbands available) by advancing again.
	for isEmptyPhase(g.adj) {
		fillNMR(g.adj)
		next, err = g.adj.Next()
		if err != nil {
			return err
		}
		g.adj = next
	}
	return nil
}

// fillNMR stages a default order for every unit that has not yet received an
// order this phase. During the Retreat phase it also fills default (disband)
// orders for dislodged units that have no retreat order.
func fillNMR(adj gameState) {
	phase := adj.Phase()
	if phase == nil {
		return
	}
	orders := adj.Orders()
	for prov := range adj.Units() {
		if _, hasOrder := orders[prov]; !hasOrder {
			if def := phase.DefaultOrder(prov); def != nil {
				adj.SetOrder(prov, def)
			}
		}
	}
	if phase.Type() == godip.Retreat {
		for prov := range adj.Dislodgeds() {
			if _, hasOrder := orders[prov]; !hasOrder {
				if def := phase.DefaultOrder(prov); def != nil {
					adj.SetOrder(prov, def)
				}
			}
		}
	}
}

// isEmptyPhase reports whether the current phase requires no player input:
//   - Retreat phase with no dislodged units
//   - Adjustment phase where every nation's supply-centre count equals its
//     unit count (no builds or disbands required)
func isEmptyPhase(adj gameState) bool {
	phase := adj.Phase()
	if phase == nil {
		return false
	}
	switch phase.Type() {
	case godip.Retreat:
		return len(adj.Dislodgeds()) == 0
	case godip.Adjustment:
		// Count supply centres per nation.
		scCount := map[godip.Nation]int{}
		for _, nation := range adj.SupplyCenters() {
			scCount[nation]++
		}
		// Count units per nation.
		unitCount := map[godip.Nation]int{}
		for _, unit := range adj.Units() {
			unitCount[unit.Nation]++
		}
		// Empty when no nation has a mismatch (build or disband required).
		for nation, scs := range scCount {
			if unitCount[nation] != scs {
				return false
			}
		}
		for nation, units := range unitCount {
			if scCount[nation] != units {
				return false
			}
		}
		return true
	default:
		return false
	}
}
