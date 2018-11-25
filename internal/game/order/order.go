package order

import (
	"github.com/burrbd/diplomacy/internal/game/order/board"
)

type Move struct {
	Country  string
	UnitType board.UnitType
	From, To board.Territory
}

type Set struct {
	Moves []*Move
}

func (s *Set) AddMove(m Move) {
	s.Moves = append(s.Moves, &m)
}

type Unresolved struct {
	move *Move
}
