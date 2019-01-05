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

func (m Move) Match(manager board.Manager) *board.Unit {
	return nil
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
	Pos      board.Territory
}

type Set struct {
	Moves        []*Move
	MoveSupports []*MoveSupport
}

func (s *Set) AddMove(m Move) {
	s.Moves = append(s.Moves, &m)
}

func (s *Set) AddMoveSupport(sup MoveSupport) {
	s.MoveSupports = append(s.MoveSupports, &sup)
}

func (s *Set) AddMoveConvoy(sup MoveConvoy) {

}

func (s *Set) AddHold(hold Hold) {

}
