package order

import (
	"github.com/burrbd/diplomacy/internal/game/order/board"
)

type Move struct {
	Country  string
	UnitType board.UnitType
	From, To board.Territory
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

func (s *Set) Strength(u *board.Unit) int {
	strength := 0
	for _, support := range s.MoveSupports {
		if u.PrevPosition != nil &&
			support.Move.From.Abbr == u.PrevPosition.Abbr &&
			support.Move.To.Abbr == u.Position.Abbr &&
			!s.supportCut(*support) {
			strength++
		}
	}
	return strength
}

func (s *Set) supportCut(support MoveSupport) bool {
	for _, cutMove := range s.Moves {
		if cutMove.To == support.By && cutMove.From != support.Move.To {
			return true
		}
	}
	return false
}

func (s Set) ByStrength(units []*board.Unit) strengthSorter {
	return strengthSorter{units, s}
}

type strengthSorter struct {
	units  []*board.Unit
	orders Set
}

func (s strengthSorter) Len() int {
	return len(s.units)
}

func (s strengthSorter) Swap(i, j int) {
	s.units[i], s.units[j] = s.units[j], s.units[i]
}

func (s strengthSorter) Less(i, j int) bool {
	return s.orders.Strength(s.units[i]) > s.orders.Strength(s.units[j])
}
