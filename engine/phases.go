package engine

import "github.com/zond/godip"

// Advance fills NMR orders for any unit without a staged order, advances the
// game to the next phase via godip Next(), and auto-skips empty retreat or
// adjustment phases (phases where no unit requires an order).
func (g *game) Advance() error {
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
// order this phase.
func fillNMR(adj godip.Adjudicator) {
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
}

// isEmptyPhase reports whether the current phase requires no player input:
//   - Retreat phase with no dislodged units
//   - Adjustment phase with no units (nothing to build or disband)
func isEmptyPhase(adj godip.Adjudicator) bool {
	phase := adj.Phase()
	if phase == nil {
		return false
	}
	switch phase.Type() {
	case godip.Retreat:
		return len(adj.Dislodgeds()) == 0
	case godip.Adjustment:
		return len(adj.Units()) == 0
	default:
		return false
	}
}
