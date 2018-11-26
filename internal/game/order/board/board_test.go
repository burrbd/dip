package board_test

import (
	"testing"

	"github.com/cheekybits/is"

	"github.com/burrbd/diplomacy/internal/game/order/board"
)

func TestPositions_Add(t *testing.T) {
	is := is.New(t)
	p := board.NewPositions()
	u := &board.Unit{}
	terr := board.Territory{Abbr: "a-territory"}
	p.Add(terr, u)
	is.NotNil(p.Units["a-territory"])
	is.Equal(u, p.Units["a-territory"][0])
}

func TestPositions_Add_ManyUnitsToTerritory(t *testing.T) {
	is := is.New(t)
	p := board.NewPositions()
	u1 := &board.Unit{}
	u2 := &board.Unit{}
	terr := board.Territory{Abbr: "a-territory"}
	p.Add(terr, u1)
	p.Add(terr, u2)
	is.Equal(2, len(p.Units["a-territory"]))
}

func TestPositions_Add_ManyUnitsToDifferentTerritories(t *testing.T) {
	is := is.New(t)
	p := board.NewPositions()
	u1 := &board.Unit{}
	u2 := &board.Unit{}
	first := board.Territory{Abbr: "first"}
	second := board.Territory{Abbr: "second"}
	p.Add(first, u1)
	p.Add(second, u2)
	is.Equal(u1, p.Units["first"][0])
	is.Equal(u2, p.Units["second"][0])
}

func TestPositions_Del(t *testing.T) {
	is := is.New(t)
	p := board.NewPositions()
	u := &board.Unit{}
	terr := board.Territory{Abbr: "a-territory"}
	p.Add(terr, u)
	err := p.Del(terr, u)
	is.NoErr(err)
	is.Equal(0, len(p.Units["a-territory"]))
}

func TestPositions_Del_ManyInTerritory(t *testing.T) {
	is := is.New(t)
	p := board.NewPositions()
	u1 := &board.Unit{}
	u2 := &board.Unit{}
	terr := board.Territory{Abbr: "terr"}

	p.Add(terr, u1)
	p.Add(terr, u2)
	err := p.Del(terr, u2)
	is.NoErr(err)
	is.NotNil(p.Units["terr"])
	is.Equal(u1, p.Units["terr"][0])
}

func TestPositions_Del_NoneInTerritory_ReturnsError(t *testing.T) {
	is := is.New(t)
	p := board.NewPositions()
	u := &board.Unit{}
	terr := board.Territory{Abbr: "terr"}
	err := p.Del(terr, u)
	is.Err(err)
}

func TestPositions_Update(t *testing.T) {
	is := is.New(t)
	p := board.NewPositions()
	u := &board.Unit{}
	prev := board.Territory{Abbr: "prev"}
	next := board.Territory{Abbr: "next"}
	p.Add(prev, u)
	err := p.Update(prev, next, u)
	is.NoErr(err)
	is.Equal(u, p.Units["next"][0])
	is.Equal(0, len(p.Units["prev"]))
	is.Equal(prev, p.Units["next"][0].PrevPositions[0])
}
