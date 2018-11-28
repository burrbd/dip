package board_test

import (
	"fmt"
	"testing"

	"github.com/cheekybits/is"

	"github.com/burrbd/diplomacy/internal/game/order/board"
)

func TestPositions_Add(t *testing.T) {
	is := is.New(t)
	terr := board.Territory{Abbr: "a-territory"}
	p := board.NewPositions([]board.Territory{terr})
	u := &board.Unit{Position: terr}
	p.Add(u)
	is.NotNil(p.Units["a-territory"])
	is.Equal(u, p.Units["a-territory"][0])
}

func TestPositions_Add_ManyUnitsToTerritory(t *testing.T) {
	is := is.New(t)
	terr := board.Territory{Abbr: "a-territory"}
	p := board.NewPositions([]board.Territory{terr})
	u1 := &board.Unit{Position: terr}
	u2 := &board.Unit{Position: terr}
	p.Add(u1)
	p.Add(u2)
	is.Equal(2, len(p.Units["a-territory"]))
}

func TestPositions_Add_ManyUnitsToDifferentTerritories(t *testing.T) {
	is := is.New(t)
	first := board.Territory{Abbr: "first"}
	second := board.Territory{Abbr: "second"}
	p := board.NewPositions([]board.Territory{first, second})
	u1 := &board.Unit{Position: first}
	u2 := &board.Unit{Position: second}

	p.Add(u1)
	p.Add(u2)
	is.Equal(u1, p.Units["first"][0])
	is.Equal(u2, p.Units["second"][0])
}

func TestPositions_Del(t *testing.T) {
	is := is.New(t)
	terr := board.Territory{Abbr: "a-territory"}
	p := board.NewPositions([]board.Territory{terr})
	u := &board.Unit{Position: terr}
	p.Add(u)
	err := p.Del(u)
	is.NoErr(err)
	is.Equal(0, len(p.Units["a-territory"]))
}

func TestPositions_Del_ManyInTerritory(t *testing.T) {
	is := is.New(t)
	terr := board.Territory{Abbr: "terr"}
	p := board.NewPositions([]board.Territory{terr})
	u1 := &board.Unit{Position: terr}
	u2 := &board.Unit{Position: terr}

	p.Add(u1)
	p.Add(u2)
	err := p.Del(u2)
	is.NoErr(err)
	is.NotNil(p.Units["terr"])
	is.Equal(u1, p.Units["terr"][0])
}

func TestPositions_Del_NoneInTerritory_ReturnsError(t *testing.T) {
	is := is.New(t)
	terr := board.Territory{Abbr: "terr"}
	p := board.NewPositions([]board.Territory{terr})
	u := &board.Unit{Position: terr}
	err := p.Del(u)
	is.Err(err)
}

func TestPositions_Update(t *testing.T) {
	is := is.New(t)
	prev := board.Territory{Abbr: "prev"}
	next := board.Territory{Abbr: "next"}
	p := board.NewPositions([]board.Territory{prev, next})
	u := &board.Unit{Position: prev}
	p.Add(u)
	err := p.Update(u, next)
	is.NoErr(err)
	is.Equal(u, p.Units["next"][0])
	is.Equal(0, len(p.Units["prev"]))
	is.Equal(prev, *p.Units["next"][0].PrevPosition)
	is.Equal(next, u.Position)
	is.Equal(prev, *u.PrevPosition)
}

func TestPositions_Conflicts(t *testing.T) {
	is := is.New(t)
	t1 := board.Territory{Abbr: "t1"}
	t2 := board.Territory{Abbr: "t2"}
	p := board.NewPositions([]board.Territory{t1, t2})
	u1 := &board.Unit{Position: t1}
	u2 := &board.Unit{Position: t1}
	u3 := &board.Unit{Position: t2}

	p.Add(u1)
	p.Add(u2)

	p.Add(u3)

	p.ConflictHandler(func(terr board.Territory, units []*board.Unit) {
		for _, u := range units {
			u.Country = fmt.Sprintf("hello %s", terr.Abbr)
		}
	})
	is.Equal("hello t1", u1.Country)
	is.Equal("hello t1", u2.Country)
}

func TestPositions_ConflictCount(t *testing.T) {
	is := is.New(t)
	t1 := board.Territory{Abbr: "t1"}
	t2 := board.Territory{Abbr: "t2"}
	t3 := board.Territory{Abbr: "t3"}
	p := board.NewPositions([]board.Territory{t1, t2, t3})
	u1 := &board.Unit{Position: t1}
	u2 := &board.Unit{Position: t1}
	u3 := &board.Unit{Position: t2}
	u4 := &board.Unit{Position: t2}
	u5 := &board.Unit{Position: t3}

	p.Add(u1)
	p.Add(u2)
	p.Add(u3)
	p.Add(u4)
	p.Add(u5)

	is.Equal(2, p.ConflictCount())
}
