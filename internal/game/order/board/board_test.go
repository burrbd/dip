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
	p.Add("a-territory", u)
	is.NotNil(p.Units["a-territory"])
	is.Equal(u, p.Units["a-territory"][0])
}

func TestPositions_Add_ManyUnitsToTerritory(t *testing.T) {
	is := is.New(t)
	p := board.NewPositions()
	u1 := &board.Unit{}
	u2 := &board.Unit{}
	p.Add("a-territory", u1)
	p.Add("a-territory", u2)
	is.Equal(2, len(p.Units["a-territory"]))
}

func TestPositions_Add_ManyUnitsToDifferentTerritories(t *testing.T) {
	is := is.New(t)
	p := board.NewPositions()
	u1 := &board.Unit{}
	u2 := &board.Unit{}
	p.Add("first", u1)
	p.Add("second", u2)
	is.Equal(u1, p.Units["first"][0])
	is.Equal(u2, p.Units["second"][0])
}

func TestPositions_Del(t *testing.T) {
	is := is.New(t)
	p := board.NewPositions()
	u := &board.Unit{}
	p.Add("a-territory", u)
	p.Del("a-territory", u)
	is.Equal(0, len(p.Units["a-territory"]))
}

func TestPositions_Del_ManyInTerritory(t *testing.T) {
	is := is.New(t)
	p := board.NewPositions()
	u1 := &board.Unit{}
	u2 := &board.Unit{}
	p.Add("terr", u1)
	p.Add("terr", u2)
	p.Del("terr", u2)

	is.NotNil(p.Units["terr"])
	is.Equal(u1, p.Units["terr"][0])
}

func TestPositions_Del_NoneInTerritory_ReturnsError(t *testing.T) {
	is := is.New(t)
	p := board.NewPositions()
	u := &board.Unit{}
	err := p.Del("terr", u)
	is.Err(err)
}
