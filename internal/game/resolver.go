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
		if ok, _ := r.ArmyGraph.IsNeighbour(m.From.Abbr, m.To.Abbr); ok {
			p.Update(m.From, m.To, unit)
		}
	}

Loop:
	for {
		p.ConflictHandler(func(terr board.Territory, units []*board.Unit) {
			for j := len(units) - 1; j >= 0; j-- {
				if len(units[j].PrevPositions) > 0 {
					p.Update(terr, units[j].PrevPositions[0], units[j])
				}
			}
		})
		if p.ConflictCount() == 0 {
			break Loop
		}
	}

	return p, nil
}
