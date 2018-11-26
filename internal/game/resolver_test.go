package game_test

import (
	"testing"

	"github.com/burrbd/diplomacy/internal/game"

	"github.com/cheekybits/is"

	"github.com/burrbd/diplomacy/internal/game/order"
	"github.com/burrbd/diplomacy/internal/game/order/board"
)

var (
	bud = board.Territory{Abbr: "bud", Name: "Budapest"}
	gal = board.Territory{Abbr: "gal", Name: "Galicia"}
	vie = board.Territory{Abbr: "vie", Name: "Vienna"}
	boh = board.Territory{Abbr: "boh", Name: "Vienna"}
	lon = board.Territory{Abbr: "lon", Name: "London"}
)

type mockGraph struct {
	IsNeighbourFunc func(t1, t2 string) (bool, error)
}

func (g mockGraph) IsNeighbour(t1, t2 string) (bool, error) {
	return g.IsNeighbourFunc(t1, t2)
}

func newPositions() board.Positions {
	return board.NewPositions([]board.Territory{bud, gal, vie, boh, lon})

}

func TestMainPhaseResolver_Resolve_HandlesMoveAndReturnsNewPositions(t *testing.T) {
	is := is.New(t)
	graph := mockGraph{
		IsNeighbourFunc: func(t1, t2 string) (bool, error) { return true, nil },
	}
	resolver := game.MainPhaseResolver{ArmyGraph: graph}
	unit := &board.Unit{}
	positions := newPositions()
	positions.Add(bud, unit)

	orders := order.Set{}
	orders.AddMove(order.Move{From: bud, To: gal})

	resolved, err := resolver.Resolve(orders, positions)

	is.NoErr(err)
	is.Equal(unit, resolved.Units["gal"][0])
	is.Equal(0, len(resolved.Units["bud"]))
}

func TestMainPhaseResolver_Resolve_HandlesAnotherMoveAndReturnsNewPositions(t *testing.T) {
	is := is.New(t)
	graph := mockGraph{
		IsNeighbourFunc: func(t1, t2 string) (bool, error) { return true, nil },
	}

	resolver := game.MainPhaseResolver{ArmyGraph: graph}

	unit := &board.Unit{}

	positions := newPositions()
	positions.Add(gal, unit)

	orders := order.Set{}
	orders.AddMove(order.Move{From: gal, To: bud})

	resolved, err := resolver.Resolve(orders, positions)

	is.NoErr(err)
	is.Equal(unit, resolved.Units["bud"][0])
	is.Equal(0, len(resolved.Units["gal"]))
}

func TestMainPhaseResolver_Resolve_OnlyMovesToNeighbouringTerritory(t *testing.T) {
	is := is.New(t)
	graph := mockGraph{
		IsNeighbourFunc: func(t1, t2 string) (bool, error) { return false, nil },
	}
	resolver := game.MainPhaseResolver{ArmyGraph: graph}

	unit := &board.Unit{}

	positions := newPositions()
	positions.Add(gal, unit)

	orders := order.Set{}
	orders.AddMove(order.Move{From: gal, To: lon})

	resolved, err := resolver.Resolve(orders, positions)

	is.NoErr(err)
	is.Nil(resolved.Units["lon"])
	is.Equal(unit, resolved.Units["gal"][0])
}

func TestMainPhaseResolver_Resolve_DoesNotMoveToOccupiedTerritory(t *testing.T) {
	is := is.New(t)
	graph := mockGraph{
		IsNeighbourFunc: func(t1, t2 string) (bool, error) { return true, nil },
	}
	resolver := game.MainPhaseResolver{ArmyGraph: graph}

	u1 := &board.Unit{}
	u2 := &board.Unit{}

	positions := newPositions()
	positions.Add(gal, u1)
	positions.Add(bud, u2)

	orders := order.Set{}
	orders.AddMove(order.Move{From: gal, To: bud})

	resolved, err := resolver.Resolve(orders, positions)

	is.NoErr(err)
	is.Equal(u1, resolved.Units["gal"][0])
	is.Equal(u2, resolved.Units["bud"][0])
}
