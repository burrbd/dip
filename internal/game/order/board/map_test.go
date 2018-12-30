package board_test

import (
	"testing"

	"github.com/cheekybits/is"

	"github.com/burrbd/diplomacy/internal/game/order/board"
)

func TestNewPositionMap_ManyUnitsToTerritory(t *testing.T) {
	is := is.New(t)
	terr := board.Territory{Abbr: "a-territory"}
	u1 := &board.Unit{Territory: terr}
	u2 := &board.Unit{Territory: terr}
	p := board.NewPositionMap([]*board.Unit{u1, u2})
	is.Equal(2, len(p.Units()))
	is.Equal(terr, p.Units()[0].Position().Territory)
	is.Equal(terr, p.Units()[1].Position().Territory)
}

func TestPositionMap_Move_UpdatesPosition(t *testing.T) {
	is := is.New(t)
	prev := board.Territory{Abbr: "prev"}
	next := board.Territory{Abbr: "next"}
	u := &board.Unit{Territory: prev}
	p := board.NewPositionMap([]*board.Unit{u})
	p.Move(u, next, 0)

	is.Equal(prev, u.PrevPosition().Territory)
	is.Equal(next, u.Position().Territory)
}

func TestPositionMap_Conflict(t *testing.T) {
	is := is.New(t)
	t1 := board.Territory{Abbr: "first"}
	u1 := &board.Unit{Territory: t1}
	u2 := &board.Unit{Territory: t1}
	r := board.NewPositionMap([]*board.Unit{u1, u2})

	conflicts := r.Conflict()
	is.Equal(2, len(conflicts))
}

func TestPositionMap_ConflictOnlyWhenInSameTerritory(t *testing.T) {
	is := is.New(t)
	t1 := board.Territory{Abbr: "t1"}
	t2 := board.Territory{Abbr: "t2"}
	u1 := &board.Unit{Territory: t1}
	u2 := &board.Unit{Territory: t1}
	u3 := &board.Unit{Territory: t2}

	r := board.NewPositionMap([]*board.Unit{u1, u2, u3})

	for _, u := range r.Conflict() {
		if u != u1 && u != u2 {
			is.Failf("%v not expected in conflict")
		}
	}
}

func TestPositionMap_Conflict_CounterAttackCausesConflict(t *testing.T) {
	is := is.New(t)
	t1 := board.Territory{Abbr: "t1"}
	t2 := board.Territory{Abbr: "t2"}
	u1 := &board.Unit{Territory: t1}
	u2 := &board.Unit{Territory: t2}

	p := board.NewPositionMap([]*board.Unit{u1, u2})
	p.Move(u1, t2, 0)
	p.Move(u2, t1, 0)

	conflicts := p.Conflict()

	for _, u := range conflicts {
		if u != u1 && u != u2 {
			is.Failf("%v expected in conflict")
		}
	}
}

func TestPositionMap_Bounce_RemovesCounterAttackConflict(t *testing.T) {
	is := is.New(t)
	t1 := board.Territory{Abbr: "first"}
	t2 := board.Territory{Abbr: "second"}
	u := &board.Unit{Territory: t1}
	p := board.NewPositionMap([]*board.Unit{u})
	p.Move(u, t2, 0)
	p.Bounce(u)
	is.Equal(0, len(p.Conflict()))
}

func TestPositionMap_Units(t *testing.T) {
	is := is.New(t)
	first := board.Territory{Abbr: "first"}
	second := board.Territory{Abbr: "second"}
	u1 := &board.Unit{Territory: first}
	u2 := &board.Unit{Territory: second}
	u3 := &board.Unit{Territory: second}
	p := board.NewPositionMap([]*board.Unit{u1, u2, u3})

	for _, u := range p.Units() {
		if u != u1 && u != u2 && u != u3 {
			is.Failf("%v not found in all units", u)
		}
	}
}

func TestPositionMap_Bounce_DoesntChangeReturnedConflictSlice(t *testing.T) {
	is := is.New(t)
	t1 := board.Territory{Abbr: "conflicted"}
	u1 := &board.Unit{Territory: t1}
	u2 := &board.Unit{PhaseHistory: []board.Position{
		{Territory: board.Territory{Abbr: "origin"}},
		{Territory: board.Territory{Abbr: "conflicted"}},
	}}
	r := board.NewPositionMap([]*board.Unit{u1, u2})

	conflicts := r.Conflict()
	is.Equal(2, len(conflicts))

	r.Bounce(conflicts[0])

	is.NotNil(conflicts[0])
	is.NotNil(conflicts[1])
}

func TestPositionMap_Move_DoesntChangeNumberOfUnits(t *testing.T) {
	is := is.New(t)
	t1 := board.Territory{Abbr: "t1"}
	t2 := board.Territory{Abbr: "t2"}
	u1 := &board.Unit{Territory: t1}
	m := board.NewPositionMap([]*board.Unit{u1})
	m.Move(u1, t2, 0)
	is.Equal(1, len(m.Units()))
}

func TestPositionMap_Bounce_DoesntChangeReturnedConflictsSliceWhenCounterAttack(t *testing.T) {
	is := is.New(t)
	t1 := board.Territory{Abbr: "first"}
	t2 := board.Territory{Abbr: "second"}
	u1 := &board.Unit{Territory: t1}
	u2 := &board.Unit{Territory: t2}
	r := board.NewPositionMap([]*board.Unit{u1, u2})
	r.Move(u1, t2, 0)
	r.Move(u2, t1, 0)
	conflicts := r.Conflict()
	is.Equal(2, len(conflicts))

	r.Bounce(conflicts[0])

	is.NotNil(conflicts[0])
	is.NotNil(conflicts[1])
}
