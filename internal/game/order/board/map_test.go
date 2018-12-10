package board_test

import (
	"testing"

	"github.com/cheekybits/is"

	"github.com/burrbd/diplomacy/internal/game/order/board"
)

func TestPositionMap_Add(t *testing.T) {
	is := is.New(t)
	terr := board.Territory{Abbr: "a-territory"}
	p := board.NewPositionMap()
	u := &board.Unit{Position: terr}
	p.Add(u)
	is.NotNil(p.Units["a-territory"])
	is.Equal(u, p.Units["a-territory"][0])
}

func TestPositionMap_Add_ManyUnitsToTerritory(t *testing.T) {
	is := is.New(t)
	terr := board.Territory{Abbr: "a-territory"}
	p := board.NewPositionMap()
	u1 := &board.Unit{Position: terr}
	u2 := &board.Unit{Position: terr}
	p.Add(u1)
	p.Add(u2)
	is.Equal(2, len(p.Units["a-territory"]))
}

func TestPositionMap_Add_ManyUnitsToDifferentTerritories(t *testing.T) {
	is := is.New(t)
	first := board.Territory{Abbr: "first"}
	second := board.Territory{Abbr: "second"}
	p := board.NewPositionMap()
	u1 := &board.Unit{Position: first}
	u2 := &board.Unit{Position: second}

	p.Add(u1)
	p.Add(u2)
	is.Equal(u1, p.Units["first"][0])
	is.Equal(u2, p.Units["second"][0])
}

func TestPositionMap_Del(t *testing.T) {
	is := is.New(t)
	terr := board.Territory{Abbr: "a-territory"}
	p := board.NewPositionMap()
	u := &board.Unit{Position: terr}
	p.Add(u)
	p.Del(u)
	is.Equal(0, len(p.Units["a-territory"]))
}

func TestPositionMap_Del_ManyInTerritory(t *testing.T) {
	is := is.New(t)
	terr := board.Territory{Abbr: "terr"}
	p := board.NewPositionMap()
	u1 := &board.Unit{Position: terr}
	u2 := &board.Unit{Position: terr}

	p.Add(u1)
	p.Add(u2)
	p.Del(u2)
	is.NotNil(p.Units["terr"])
	is.Equal(u1, p.Units["terr"][0])
}

func TestPositionMap_Move(t *testing.T) {
	is := is.New(t)
	prev := board.Territory{Abbr: "prev"}
	next := board.Territory{Abbr: "next"}
	p := board.NewPositionMap()
	u := &board.Unit{Position: prev}
	p.Add(u)
	p.Move(u, next)
	is.Equal(u, p.Units["next"][0])
	is.Equal(0, len(p.Units["prev"]))
	is.Equal(prev, *p.Units["next"][0].PrevPosition)
	is.Equal(next, u.Position)
	is.Equal(prev, *u.PrevPosition)
}

func TestPositionMap_Conflicts(t *testing.T) {
	is := is.New(t)
	t1 := board.Territory{Abbr: "t1"}
	t2 := board.Territory{Abbr: "t2"}
	r := board.NewPositionMap()
	u1 := &board.Unit{Position: t1}
	u2 := &board.Unit{Position: t1}
	u3 := &board.Unit{Position: t2}

	r.Add(u1)
	r.Add(u2)
	r.Add(u3)

	is.Equal([]*board.Unit{u2, u1}, r.GetConflict())
}

func TestPositionMap_Bounce_Removes_CounterAttackConflict(t *testing.T) {
	is := is.New(t)
	t1 := board.Territory{Abbr: "first"}
	t2 := board.Territory{Abbr: "second"}
	p := board.NewPositionMap()
	u := &board.Unit{Position: t1}
	p.Add(u)
	p.Move(u, t2)
	p.Bounce(u, t1)
	is.Equal(0, len(p.CounterConflicts["firstsecond"]))
}

func TestPositionMap_GetConflicts_CounterAttack_CausesConflict(t *testing.T) {
	is := is.New(t)
	t1 := board.Territory{Abbr: "t1"}
	t2 := board.Territory{Abbr: "t2"}
	p := board.NewPositionMap()
	u1 := &board.Unit{Position: t1}
	u2 := &board.Unit{Position: t2}

	p.Add(u1)
	p.Add(u2)
	p.Move(u1, t2)
	p.Move(u2, t1)

	conflicts := p.GetConflict()

	for _, u := range conflicts {
		if u != u1 && u != u2 {
			is.Failf("%v expected to be in conflict")
		}
	}
}

func TestPositionMap_AllUnits(t *testing.T) {
	is := is.New(t)
	first := board.Territory{Abbr: "first"}
	second := board.Territory{Abbr: "second"}
	p := board.NewPositionMap()
	u1 := &board.Unit{Position: first}
	u2 := &board.Unit{Position: second}
	u3 := &board.Unit{Position: second}

	p.Add(u1)
	p.Add(u2)
	p.Add(u3)

	for _, u := range p.AllUnits() {
		if u != u1 && u != u2 && u != u3 {
			is.Failf("%v not found in all units", u)
		}
	}
}

func TestPositionMap_GetConflict(t *testing.T) {
	is := is.New(t)
	t1 := board.Territory{Abbr: "first"}
	r := board.NewPositionMap()
	u1 := &board.Unit{Position: t1}
	u2 := &board.Unit{Position: t1}
	r.Add(u1)
	r.Add(u2)

	conflicts := r.GetConflict()
	is.Equal(2, len(conflicts))
}

func TestPositionMap_Bounce_DoesntChangeConflictsSlice(t *testing.T) {
	is := is.New(t)
	t1 := board.Territory{Abbr: "first"}
	t2 := board.Territory{Abbr: "second"}
	r := board.NewPositionMap()
	u1 := &board.Unit{Position: t1}
	u2 := &board.Unit{Position: t1}
	r.Add(u1)
	r.Add(u2)

	conflicts := r.GetConflict()
	is.Equal(2, len(conflicts))

	r.Bounce(conflicts[0], t2)

	is.NotNil(conflicts[0])
	is.NotNil(conflicts[1])
}

func TestPositionMap_Bounce_DoesntChangeConflictsSliceWhenCounterAttack(t *testing.T) {
	is := is.New(t)
	t1 := board.Territory{Abbr: "first"}
	t2 := board.Territory{Abbr: "second"}
	r := board.NewPositionMap()
	u1 := &board.Unit{Position: t1}
	u2 := &board.Unit{Position: t2}
	r.Add(u1)
	r.Add(u2)
	r.Move(u1, t2)
	r.Move(u2, t1)
	conflicts := r.GetConflict()
	is.Equal(2, len(conflicts))

	r.Bounce(conflicts[0], t2)

	is.NotNil(conflicts[0])
	is.NotNil(conflicts[1])
}
