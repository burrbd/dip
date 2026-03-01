package engine

// soloThreshold is the number of supply centres required to achieve a solo
// victory in the Classical variant.
const soloThreshold = 18

// SoloWinner returns the nation that has achieved a solo victory (≥18 supply
// centres), or an empty string if no nation has won yet.
//
// godip tracks supply-centre ownership internally; we delegate directly to
// Adjudicator.SoloWinner() which performs the check after each Fall
// Adjustment phase.
func (g *game) SoloWinner() string {
	return string(g.adj.SoloWinner())
}
