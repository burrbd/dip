package order

import "github.com/burrbd/dip/game/order/board"

type UnitMatcher interface {
	Match(manager board.Manager) *board.Unit
}

type Move struct {
	Country  string
	UnitType board.UnitType
	From, To board.Territory
	Strength int
}

type MoveSupport struct {
	Country  string
	UnitType board.UnitType
	By       board.Territory
	Move     Move
}

type MoveConvoy struct {
	Country string
	By      board.Territory
	Move    Move
}

type Hold struct {
	Country  string
	UnitType board.UnitType
	At       board.Territory
}

type HoldSupport struct {
	UnitType board.UnitType
	By       board.Territory
	Hold     Hold
}

type Set struct {
	Moves        []Move
	MoveSupports []MoveSupport
	HoldSupports []HoldSupport
}

func (s *Set) AddMove(m Move) {
	s.Moves = append(s.Moves, m)
}

func (s *Set) AddHold(hold Hold) {

}

func (s *Set) AddMoveSupport(sup MoveSupport) {
	s.MoveSupports = append(s.MoveSupports, sup)
}

func (s *Set) AddHoldSupport(sup HoldSupport) {
	s.HoldSupports = append(s.HoldSupports, sup)
}

func (s *Set) AddMoveConvoy(sup MoveConvoy) {

}
