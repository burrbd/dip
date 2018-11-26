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
	for _, m := range s.Moves {
		unit := p.Units[m.From.Abbr][0]
		if neighbour, _ := r.ArmyGraph.IsNeighbour(m.From.Abbr, m.To.Abbr); neighbour {
			p.Update(m.From, m.To, unit)
		}
	}
	for i, units := range p.Units {
		if len(units) > 1 {
			for _, u := range units {
				if len(u.PrevPositions) == 1 {
					terr, _ := p.Territory(i)
					p.Update(terr, u.PrevPositions[0], u)
				}
			}
		}
	}
	return p, nil
}
