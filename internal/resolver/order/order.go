package order

import (
	"errors"

	"github.com/burrbd/kit/graph"

	"log"

	"github.com/burrbd/diplomacy/internal/resolver/order/board"
)

type Move struct {
	Country         string
	UnitType        board.UnitType
	From, To        board.Territory
	MatchedPosition *board.Position
	Result          *Result
}

func (m Move) matchPosition(p *board.Position) bool {
	if p.Territory.ID() == m.From.ID() &&
		p.Unit.Country == m.Country &&
		p.Unit.Type == m.UnitType {
		m.MatchedPosition = p
		return true
	}
	return false
}

type Result struct {
	Success     bool
	Description string
}

type Set struct {
	Positions                          []*board.Position
	ArmyGraph, FleetGraph, ConvoyGraph *graph.Simple
	moves                              []*Move
}

func (s *Set) AddMove(m Move) {
	s.moves = append(s.moves, &m)
}

type Unresolved struct {
	move *Move
	curr *board.Position
}

func (s *Set) Resolve() ([]Move, error) {
	for _, m := range s.moves {
		ok, err := s.ArmyGraph.IsNeighbor(m.From, m.To)
		if err != nil {
			return []Move{}, err
		}
		if !ok {
			m.Result = &Result{Success: !ok}
			continue
		}
		hasPosition := false
		// loop through all board positions to check if there is a
		// valid position for this move
		for _, p := range s.Positions {
			if !m.matchPosition(p) {
				continue
			}
			hasPosition = true
			err := p.Move(m.To, *s.ArmyGraph)
			if err != nil {
				// We do not expect this because we've checked for neighbor
				// above
				return nil, errors.New("something wrong")
			}
			break
		}
		if !hasPosition {
			m.Result = &Result{Success: false}
			continue
		}
	}
	return resolve(s.moves, nil), nil
}

func resolve(moves []*Move, unresolved map[string][]*Move) []Move {
	if unresolved == nil {
		unresolved = make(map[string][]*Move)
	}
	for _, m := range moves {
		if m.Result != nil {
			continue
		}
		currPosition := m.MatchedPosition.Territory.ID()
		unresolved[currPosition] = append(unresolved[currPosition], m)
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
			p.move.Success = success
			results = append(results, *p.move)
		}
	}
	return results
}
