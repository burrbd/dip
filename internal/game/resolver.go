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
			p.Update(unit, m.To)
		}
	}

Loop:
	for {
		p.ConflictHandler(func(units []*board.Unit) {
			sort.Sort(s.ByStrength(units))
			if s.Strength(units[0]) > s.Strength(units[1]) {
				handleDefeats(units[1:], p)
			} else {
				handleBounces(units, p)
			}
		})
		if p.ConflictCount() == 0 {
			break Loop
		}
	}
	return p, nil
}

func handleDefeats(units []*board.Unit, p board.Positions) {
	for _, u := range units {
		if previousPosition(u) {
			p.Update(u, *u.PrevPosition)
		} else {
			u.MustRetreat = true
		}
	}
}

func handleBounces(units []*board.Unit, p board.Positions) {
	for _, u := range units {
		if previousPosition(u) {
			p.Update(u, *u.PrevPosition)
		}
	}
}

func previousPosition(u *board.Unit) bool {
	return u.PrevPosition != nil && *u.PrevPosition != u.Position
}
