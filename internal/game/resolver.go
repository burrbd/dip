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
			p.Update(unit, m.To)
		}
	}

Loop:
	for {
		p.ConflictHandler(func(terr board.Territory, units []*board.Unit) {
			unitStrength, prevUnitStrength := 0, 0
			for j := len(units) - 1; j >= 0; j-- {
				unitStrength = s.Strength(units[j])
				if units[j].PrevPosition != nil && unitStrength <= prevUnitStrength {
					p.Update(units[j], *units[j].PrevPosition)
				}
				prevUnitStrength = unitStrength
			}
		})
		if p.ConflictCount() == 0 {
			break Loop
		}
	}

	return p, nil
}
