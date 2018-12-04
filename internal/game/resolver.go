package game

import (
	"sort"

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
			p.Move(unit, m.To)
		}
	}

Loop:
	for {
		p.ConflictHandler(func(units []*board.Unit) {
			defeated := false
			sort.Sort(s.ByStrength(units))
			if s.Strength(units[0]) > s.Strength(units[1]) {
				units, defeated = units[1:], true
			}
			for _, u := range units {
				if u.OriginalPosition() {
					u.Defeated = defeated
				} else {
					p.Bounce(u, *u.PrevPosition)
				}
			}
		})
		if p.ConflictCount() == 0 {
			break Loop
		}
	}
	return p, nil
}
