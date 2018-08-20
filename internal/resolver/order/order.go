package order

import (
	"github.com/burrbd/kit/graph"

	"github.com/burrbd/diplomacy/internal/resolver/order/board"
)

type Move struct {
	Country  string
	UnitType board.UnitType
	From, To board.Territory
}

type Set struct {
	Positions                          []board.Position
	ArmyGraph, FleetGraph, ConvoyGraph *graph.Simple
	moves                              []Move
}

func (s *Set) AddMove(m Move) {
	s.moves = append(s.moves, m)
}

func (s *Set) Resolve() ([]Result, error) {
	res := make([]Result, 0)
	new := make([]board.Position, len(s.Positions))
	copy(new, s.Positions)
	unresolved := make([]Move, 0)
	for _, m := range s.moves {
		ok, err := s.ArmyGraph.IsNeighbor(m.From, m.To)
		if err != nil {
			return Results{}, err
		}
		if !ok {
			res = append(res, Result{Move: m, Success: false})
			continue
		}
		hasMatch := false
		for i := len(s.Positions) - 1; i >= 0; i-- {
			p := s.Positions[i]
			if matchPosition(p, m) {
				new[i].Territory = m.To
				hasMatch = true
				break
			}
		}
		if !hasMatch {
			res = append(res, Result{Move: m, Success: false})
			continue
		}
		unresolved = append(unresolved, m)
	}
	contested := contests(new)
	for _, m := range s.moves {
		if n, ok := contested[m.To.Abbr]; ok && n > 0 {
			res = append(res, Result{Move: m, Success: false})
		} else {
			res = append(res, Result{Move: m, Success: true})
		}
	}
	return res, nil
}

// contests probably throw away
func contests(p []board.Position) map[string]int {
	cnt := make(map[string]int)
	for _, pp := range p {
		id := pp.Territory.ID()
		if _, ok := cnt[id]; !ok {
			cnt[id] = 0
		} else {
			cnt[id]++
		}
	}
	return cnt
}

func matchPosition(p board.Position, m Move) bool {
	return p.Territory.ID() == m.From.ID() &&
		p.Unit.Country == m.Country &&
		p.Unit.Type == m.UnitType
}

type Result struct {
	Move    Move
	Success bool
}

type Results []Result
