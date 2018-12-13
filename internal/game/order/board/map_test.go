package board_test

import (
	"testing"

	"github.com/cheekybits/is"

	"github.com/burrbd/diplomacy/internal/game/order/board"
)

func TestNewPositionMap_ManyUnitsToTerritory(t *testing.T) {
	is := is.New(t)
	terr := board.Territory{Abbr: "a-territory"}
	u1 := &board.Unit{Position: terr}
	u2 := &board.Unit{Position: terr}
	p := board.NewPositionMap([]*board.Unit{u1, u2})
	is.Equal(2, len(p.Units()))
	is.Equal("a-territory", p.Units()[0].Position.Abbr)
	is.Equal("a-territory", p.Units()[1].Position.Abbr)
}

func TestPositionMap_Move_UpdatesPosition(t *testing.T) {
	is := is.New(t)
	prev := board.Territory{Abbr: "prev"}
	next := board.Territory{Abbr: "next"}
	u := &board.Unit{Position: prev}
	p := board.NewPositionMap([]*board.Unit{u})
	p.Move(u, next)

	is.Equal(prev, *u.PrevPosition)
	is.Equal(next, u.Position)
}

func TestPositionMap_Conflict(t *testing.T) {
	is := is.New(t)
	t1 := board.Territory{Abbr: "first"}
	u1 := &board.Unit{Position: t1}
	u2 := &board.Unit{Position: t1}
	r := board.NewPositionMap([]*board.Unit{u1, u2})

	conflicts := r.Conflict()
	is.Equal(2, len(conflicts))
}

func TestPositionMap_ConflictOnlyWhenInSameTerritory(t *testing.T) {
	is := is.New(t)
	t1 := board.Territory{Abbr: "t1"}
	t2 := board.Territory{Abbr: "t2"}
	u1 := &board.Unit{Position: t1}
	u2 := &board.Unit{Position: t1}
	u3 := &board.Unit{Position: t2}

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
	u1 := &board.Unit{Position: t1}
	u2 := &board.Unit{Position: t2}

	p := board.NewPositionMap([]*board.Unit{u1, u2})
	p.Move(u1, t2)
	p.Move(u2, t1)

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
	u := &board.Unit{Position: t1}
	p := board.NewPositionMap([]*board.Unit{u})
	p.Move(u, t2)
	p.Bounce(u, t1)
	is.Equal(0, len(p.Conflict()))
}

func TestPositionMap_Units(t *testing.T) {
	is := is.New(t)
	first := board.Territory{Abbr: "first"}
	second := board.Territory{Abbr: "second"}
	u1 := &board.Unit{Position: first}
	u2 := &board.Unit{Position: second}
	u3 := &board.Unit{Position: second}
	p := board.NewPositionMap([]*board.Unit{u1, u2, u3})

	for _, u := range p.Units() {
		if u != u1 && u != u2 && u != u3 {
			is.Failf("%v not found in all units", u)
		}
	}
}

func TestPositionMap_Bounce_DoesntChangeReturnedConflictSlice(t *testing.T) {
	is := is.New(t)
	t1 := board.Territory{Abbr: "first"}
	t2 := board.Territory{Abbr: "second"}
	u1 := &board.Unit{Position: t1}
	u2 := &board.Unit{Position: t1}
	r := board.NewPositionMap([]*board.Unit{u1, u2})

	conflicts := r.Conflict()
	is.Equal(2, len(conflicts))

	r.Bounce(conflicts[0], t2)

	is.NotNil(conflicts[0])
	is.NotNil(conflicts[1])
}

func TestPositionMap_Bounce_DoesntChangeReturnedConflictsSliceWhenCounterAttack(t *testing.T) {
	is := is.New(t)
	t1 := board.Territory{Abbr: "first"}
	t2 := board.Territory{Abbr: "second"}
	u1 := &board.Unit{Position: t1}
	u2 := &board.Unit{Position: t2}
	r := board.NewPositionMap([]*board.Unit{u1, u2})
	r.Move(u1, t2)
	r.Move(u2, t1)
	conflicts := r.Conflict()
	is.Equal(2, len(conflicts))

	r.Bounce(conflicts[0], t2)

	is.NotNil(conflicts[0])
	is.NotNil(conflicts[1])
}
