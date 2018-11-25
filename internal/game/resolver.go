package game

import (
	"github.com/burrbd/diplomacy/internal/game/order"
	"github.com/burrbd/diplomacy/internal/game/order/board"
)

type Resolver interface {
	Resolve(set order.Set, positions board.Positions) (board.Positions, error)
}

type MainPhaseResolver struct {
	ArmyGraph, FleetGraph, ConvoyGraph board.Graph
}

func (r MainPhaseResolver) Resolve(s order.Set, p board.Positions) (board.Positions, error) {
	newP := board.NewPositions()
	for _, m := range s.Moves {
		curr := m.From.Abbr
		to := m.To.Abbr
		unit := p.Units[curr][0]
		if neighbour, _ := r.ArmyGraph.IsNeighbour(curr, to); neighbour {
			curr = to
		}
		newP.Add(curr, unit)
	}
	return newP, nil
}
