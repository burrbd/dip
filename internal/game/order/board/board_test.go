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
	u := &board.Unit{}
	p.Add(terr, u)
	is.NotNil(p.Units["a-territory"])
	is.Equal(u, p.Units["a-territory"][0])
}

func TestPositions_Add_ManyUnitsToTerritory(t *testing.T) {
	is := is.New(t)
	terr := board.Territory{Abbr: "a-territory"}
	p := board.NewPositions([]board.Territory{terr})
	u1 := &board.Unit{}
	u2 := &board.Unit{}
	p.Add(terr, u1)
	p.Add(terr, u2)
	is.Equal(2, len(p.Units["a-territory"]))
}

func TestPositions_Add_ManyUnitsToDifferentTerritories(t *testing.T) {
	is := is.New(t)
	first := board.Territory{Abbr: "first"}
	second := board.Territory{Abbr: "second"}
	p := board.NewPositions([]board.Territory{first, second})
	u1 := &board.Unit{}
	u2 := &board.Unit{}

	p.Add(first, u1)
	p.Add(second, u2)
	is.Equal(u1, p.Units["first"][0])
	is.Equal(u2, p.Units["second"][0])
}

func TestPositions_Del(t *testing.T) {
	is := is.New(t)
	terr := board.Territory{Abbr: "a-territory"}
	p := board.NewPositions([]board.Territory{terr})
	u := &board.Unit{}
	p.Add(terr, u)
	err := p.Del(terr, u)
	is.NoErr(err)
	is.Equal(0, len(p.Units["a-territory"]))
}

func TestPositions_Del_ManyInTerritory(t *testing.T) {
	is := is.New(t)
	terr := board.Territory{Abbr: "terr"}
	p := board.NewPositions([]board.Territory{terr})
	u1 := &board.Unit{}
	u2 := &board.Unit{}

	p.Add(terr, u1)
	p.Add(terr, u2)
	err := p.Del(terr, u2)
	is.NoErr(err)
	is.NotNil(p.Units["terr"])
	is.Equal(u1, p.Units["terr"][0])
}

func TestPositions_Del_NoneInTerritory_ReturnsError(t *testing.T) {
	is := is.New(t)
	terr := board.Territory{Abbr: "terr"}
	p := board.NewPositions([]board.Territory{terr})
	u := &board.Unit{}
	err := p.Del(terr, u)
	is.Err(err)
}

func TestPositions_Update(t *testing.T) {
	is := is.New(t)
	prev := board.Territory{Abbr: "prev"}
	next := board.Territory{Abbr: "next"}
	p := board.NewPositions([]board.Territory{prev, next})
	u := &board.Unit{}
	p.Add(prev, u)
	err := p.Update(prev, next, u)
	is.NoErr(err)
	is.Equal(u, p.Units["next"][0])
	is.Equal(0, len(p.Units["prev"]))
	is.Equal(prev, p.Units["next"][0].PrevPositions[0])
}

func TestPositions_Conflicts(t *testing.T) {
	is := is.New(t)
	t1 := board.Territory{Abbr: "t1"}
	t2 := board.Territory{Abbr: "t2"}
	p := board.NewPositions([]board.Territory{t1, t2})
	u1 := &board.Unit{}
	u2 := &board.Unit{}
	u3 := &board.Unit{}

	p.Add(t1, u1)
	p.Add(t1, u2)

	p.Add(t2, u3)

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
	u1 := &board.Unit{}
	u2 := &board.Unit{}
	u3 := &board.Unit{}
	u4 := &board.Unit{}
	u5 := &board.Unit{}

	p.Add(t1, u1)
	p.Add(t1, u2)
	p.Add(t2, u3)
	p.Add(t2, u4)
	p.Add(t3, u5)

	is.Equal(2, p.ConflictCount())
}
