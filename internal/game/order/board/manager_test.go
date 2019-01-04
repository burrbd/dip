package board_test

import (
	"testing"

	"github.com/cheekybits/is"

	"github.com/burrbd/diplomacy/internal/game/order/board"
)

func TestPositionManager_Move_UpdatesPosition(t *testing.T) {
	is := is.New(t)
	prev := board.Territory{Abbr: "prev"}
	next := board.Territory{Abbr: "next"}
	u := &board.Unit{}
	m := board.NewPositionManager()
	m.AddUnit(u, prev)
	m.Move(u, next, 0)
	is.Equal(prev, u.PrevPosition().Territory)
	is.Equal(next, u.Position().Territory)
}

func TestPositionManager_Conflict(t *testing.T) {
	is := is.New(t)
	u1 := &board.Unit{}
	u2 := &board.Unit{}
	m := board.NewPositionManager()
	m.AddUnit(u1, board.Territory{Abbr: "t1"})
	m.AddUnit(u2, board.Territory{Abbr: "t1"})
	conflicts := m.Conflict()
	is.Equal(2, len(conflicts))
}

func TestPositionManager_ConflictOnlyWhenInSameTerritory(t *testing.T) {
	is := is.New(t)
	u1 := &board.Unit{}
	u2 := &board.Unit{}
	u3 := &board.Unit{}
	m := board.NewPositionManager()
	m.AddUnit(u1, board.Territory{Abbr: "t1"})
	m.AddUnit(u2, board.Territory{Abbr: "t1"})
	m.AddUnit(u3, board.Territory{Abbr: "t2"})
	is.Equal(2, len(m.Conflict()))
	for _, u := range m.Conflict() {
		if u != u1 && u != u2 {
			is.Failf("%v not expected in conflict")
		}
	}
}

func TestPositionManager_Conflict_CounterAttackCausesConflict(t *testing.T) {
	is := is.New(t)
	t1 := board.Territory{Abbr: "t1"}
	t2 := board.Territory{Abbr: "t2"}
	u1 := &board.Unit{}
	u2 := &board.Unit{}
	m := board.NewPositionManager()
	m.AddUnit(u1, t1)
	m.AddUnit(u2, t2)
	m.Move(u1, t2, 0)
	m.Move(u2, t1, 0)
	is.Equal(2, len(m.Conflict()))
	for _, u := range m.Conflict() {
		if u != u1 && u != u2 {
			is.Failf("%v expected in conflict")
		}
	}
}

func TestPositionManager_Bounce_RemovesCounterAttackConflict(t *testing.T) {
	is := is.New(t)
	a := board.Territory{Abbr: "a"}
	b := board.Territory{Abbr: "b"}
	u1 := &board.Unit{}
	u2 := &board.Unit{}
	m := board.NewPositionManager()
	m.AddUnit(u1, a)
	m.AddUnit(u2, b)
	m.Move(u1, b, 0)
	m.Move(u2, a, 0)
	is.Equal(2, len(m.Conflict()))
	m.Bounce(u1)
	m.Bounce(u2)
	is.Equal(0, len(m.Conflict()))
}

func TestPositionManager_Units(t *testing.T) {
	is := is.New(t)
	u1 := &board.Unit{}
	u2 := &board.Unit{}
	u3 := &board.Unit{}
	m := board.NewPositionManager()
	m.AddUnit(u1, board.Territory{Abbr: "terr1"})
	m.AddUnit(u2, board.Territory{Abbr: "terr2"})
	m.AddUnit(u3, board.Territory{Abbr: "terr3"})
	for _, u := range m.Units() {
		if u != u1 && u != u2 && u != u3 {
			is.Failf("%v not found in all units", u)
		}
	}
}

func TestPositionManager_Move_DoesntChangeNumberOfUnits(t *testing.T) {
	is := is.New(t)
	u1 := &board.Unit{}
	m := board.NewPositionManager()
	m.AddUnit(u1, board.Territory{Abbr: "a-terr"})
	m.Move(u1, board.Territory{Abbr: "b-terr"}, 0)
	is.Equal(1, len(m.Units()))
}

func TestPositionManager_Bounce_DoesntChangeReturnedConflictSlice(t *testing.T) {
	is := is.New(t)
	terr := board.Territory{Abbr: "a terr"}
	u1 := &board.Unit{}
	u2 := &board.Unit{}
	m := board.NewPositionManager()
	m.AddUnit(u1, terr)
	m.AddUnit(u2, terr)
	conflicts := m.Conflict()
	is.Equal(2, len(conflicts))
	m.Bounce(conflicts[0])
	is.NotNil(conflicts[0])
	is.NotNil(conflicts[1])
}

func TestPositionManager_Bounce_DoesntChangeReturnedConflictsSliceWhenCounterAttack(t *testing.T) {
	is := is.New(t)
	a := board.Territory{Abbr: "a-terr"}
	b := board.Territory{Abbr: "b-terr"}
	u1 := &board.Unit{}
	u2 := &board.Unit{}
	m := board.NewPositionManager()
	m.AddUnit(u1, a)
	m.AddUnit(u2, b)
	m.Move(u1, b, 0)
	m.Move(u2, a, 0)
	conflicts := m.Conflict()
	is.Equal(2, len(conflicts))
	m.Bounce(conflicts[0])
	is.NotNil(conflicts[0])
	is.NotNil(conflicts[1])
}

func TestPositionManager_AddUnit_UpdatesPhaseHistory(t *testing.T) {
	is := is.New(t)
	terr := board.Territory{Abbr: "terr"}
	u := &board.Unit{}
	m := board.NewPositionManager()
	m.AddUnit(u, terr)
	is.Equal(board.Originated, u.PhaseHistory[0].Cause)
	is.Equal(terr, u.PhaseHistory[0].Territory)
}

func TestPositionManager_AddUnit_AddsToUnits(t *testing.T) {
	is := is.New(t)
	terr := board.Territory{Abbr: "terr"}
	u := &board.Unit{}
	m := board.NewPositionManager()
	m.AddUnit(u, terr)
	is.Equal(u, m.Units()[0])
}
