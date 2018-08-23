package order

import (
	"github.com/burrbd/kit/graph"

	"log"

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

type Unresolved struct {
	move *Move
	curr *board.Position
}

func (s *Set) Resolve() ([]Result, error) {
	res := make([]Result, 0)
	unresolved := make(map[string][]Unresolved)
	for _, p := range s.Positions {
		terr := p.Territory.ID()
		unresolved[terr] = []Unresolved{{curr: &p}}
	}
	for _, m := range s.moves {
		ok, err := s.ArmyGraph.IsNeighbor(m.From, m.To)
		if err != nil {
			return Results{}, err
		}
		if !ok {
			res = append(res, Result{Move: m, Success: false})
			continue
		}
		hasPosition := false
		// loop through all board positions to check if there is a
		// valid position for this move
		for _, p := range s.Positions {
			terr := p.Territory.ID()
			if matchPosition(p, m) {
				// delete old position (this feels a bit risky)
				unresolved[terr] = append(unresolved[terr][:0], unresolved[terr][1:]...)

				newTerr := m.To.ID()
				p.Territory = m.To
				if _, ok := unresolved[newTerr]; !ok {
					unresolved[newTerr] = make([]Unresolved, 0)
				}
				unresolved[newTerr] = append(unresolved[newTerr], Unresolved{move: &m, curr: &p})
				hasPosition = true
				break
			}
		}
		if !hasPosition {
			log.Println("what!")
			res = append(res, Result{Move: m, Success: false})
			continue
		}
	}
	return resolve(unresolved, res), nil
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

func resolve(unresolved map[string][]Unresolved, results []Result) []Result {
	log.Printf("start: %+v\n", unresolved)
	newUnresolved := make(map[string][]Unresolved)
	for terr, positions := range unresolved {
		if positions == nil {
			continue
		}
		if len(positions) < 2 {
			newUnresolved[terr] = positions
			continue
		}
		newUnresolved[terr] = make([]Unresolved, 0)
		for _, p := range positions {
			var newPos Unresolved
			if p.move != nil {
				p.curr.Territory = p.move.From
				newPos = Unresolved{move: p.move, curr: p.curr}
			} else {
				newPos = p
			}
			newTerr := newPos.curr.Territory.ID()
			if _, ok := newUnresolved[newTerr]; !ok {
				newUnresolved[newTerr] = make([]Unresolved, 0)
			}
			newUnresolved[newTerr] = append(newUnresolved[newTerr], newPos)
		}
	}
	for _, positions := range newUnresolved {
		if len(positions) > 1 {
			resolve(newUnresolved, results)
		}
	}
	log.Printf("%+v\n", newUnresolved)
	for terr, positions := range newUnresolved {
		for _, p := range positions {
			if p.move == nil {
				continue
			}
			success := false
			if p.move.To.ID() == terr {
				success = true
			}
			results = append(results, Result{Move: *p.move, Success: success})
		}
	}
	return results
}
