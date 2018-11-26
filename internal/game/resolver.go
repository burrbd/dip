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
		old, curr := m.From, m.From
		to := m.To
		unit := p.Units[curr.Abbr][0]
		if neighbour, _ := r.ArmyGraph.IsNeighbour(curr.Abbr, to.Abbr); neighbour {
			if unit.PrevPositions == nil {
				unit.PrevPositions = make([]board.Territory, 0, 1)
			}
			curr = to
			unit.PrevPositions = append(unit.PrevPositions, old)
		}
		p.Add(curr.Abbr, unit)
		p.Del(old.Abbr, unit)
	}
	for terr, units := range p.Units {
		if len(units) > 1 {
			for _, u := range units {
				if len(u.PrevPositions) == 1 {
					p.Add(u.PrevPositions[0].Abbr, u)
					p.Del(terr, u)
				}
			}
		}
	}
	return p, nil
}
