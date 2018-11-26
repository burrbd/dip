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

func TestMainPhaseResolver_Resolve_HandlesMoveAndReturnsNewPositions(t *testing.T) {
	is := is.New(t)
	graph := mockGraph{
		IsNeighbourFunc: func(t1, t2 string) (bool, error) { return true, nil },
	}

	resolver := game.MainPhaseResolver{ArmyGraph: graph}
	unit := &board.Unit{}
	positions := board.NewPositions()
	positions.Add("bud", unit)

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

	positions := board.NewPositions()
	positions.Add("gal", unit)

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

	positions := board.NewPositions()
	positions.Add("gal", unit)

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

	positions := board.NewPositions()
	positions.Add("gal", u1)
	positions.Add("bud", u2)

	orders := order.Set{}
	orders.AddMove(order.Move{From: gal, To: bud})

	resolved, err := resolver.Resolve(orders, positions)

	is.NoErr(err)
	is.Equal(u1, resolved.Units["gal"][0])
	is.Equal(u2, resolved.Units["bud"][0])
}

//
//func TestResolveArmyMovesWrongCountry(t *testing.T) {
//	is := is.New(t)
//	g := armyGraph()
//	bud := board.Territory{Abbr: "bud", Name: "Budapest"}
//	gal := board.Territory{Abbr: "gal", Name: "Galicia"}
//	unit := board.Unit{Country: "France", Type: board.Army}
//	set := order.Set{ArmyGraph: g, Positions: []board.Position{{Territory: bud, Unit: unit}}}
//	set.AddMove(order.Move{Country: "France", UnitType: board.Army, From: bud, To: gal})
//	r, err := set.Resolve()
//	is.NoErr(err)
//	is.True(r[0].Success)
//}
//
//func TestResolveArmyMovesToOccupiedSpace(t *testing.T) {
//	is := is.New(t)
//	g := armyGraph()
//	bud := board.Territory{Abbr: "bud", Name: "Budapest"}
//	gal := board.Territory{Abbr: "gal", Name: "Galicia"}
//	ubud := board.Unit{Country: "France", Type: board.Army}
//	ugal := board.Unit{Country: "Austria-Hungary", Type: board.Army}
//	set := order.Set{
//		ArmyGraph: g,
//		Positions: []board.Position{
//			{Territory: bud, Unit: ubud},
//			{Territory: gal, Unit: ugal},
//		}}
//	set.AddMove(order.Move{Country: "France", UnitType: board.Army, From: bud, To: gal})
//	r, err := set.Resolve()
//	is.NoErr(err)
//	is.Equal(1, len(r))
//	is.False(r[0].Success)
//}
//
//func TestResolveArmyMovesToOccTerrAndOccTerrMoves(t *testing.T) {
//	is := is.New(t)
//	g := armyGraph()
//	bud := board.Territory{Abbr: "bud", Name: "Budapest"}
//	gal := board.Territory{Abbr: "gal", Name: "Galicia"}
//	vie := board.Territory{Abbr: "vie", Name: "Vienna"}
//	ubud := board.Unit{Country: "France", Type: board.Army}
//	ugal := board.Unit{Country: "Austria-Hungary", Type: board.Army}
//	set := order.Set{
//		ArmyGraph: g,
//		Positions: []board.Position{
//			{Territory: bud, Unit: ubud},
//			{Territory: gal, Unit: ugal},
//		}}
//	set.AddMove(order.Move{Country: "France", UnitType: board.Army, From: bud, To: gal})
//	set.AddMove(order.Move{Country: "Austria-Hungary", UnitType: board.Army, From: gal, To: vie})
//	r, err := set.Resolve()
//	is.NoErr(err)
//	is.Equal(2, len(r))
//	is.True(r[0].Success)
//	is.True(r[1].Success)
//}
//
//func TestResolveArmyMovesToOccTerrAndOccTerrMoves2(t *testing.T) {
//	is := is.New(t)
//	g := armyGraph()
//	bud := board.Territory{Abbr: "bud", Name: "Budapest"}
//	gal := board.Territory{Abbr: "gal", Name: "Galicia"}
//	vie := board.Territory{Abbr: "vie", Name: "Vienna"}
//	ubud := board.Unit{Country: "France", Type: board.Army}
//	ugal := board.Unit{Country: "France", Type: board.Army}
//	set := order.Set{
//		ArmyGraph: g,
//		Positions: []board.Position{
//			{Territory: bud, Unit: ubud},
//			{Territory: gal, Unit: ugal},
//		}}
//	set.AddMove(order.Move{Country: "France", UnitType: board.Army, From: bud, To: gal})
//	set.AddMove(order.Move{Country: "France", UnitType: board.Army, From: gal, To: vie})
//	r, err := set.Resolve()
//	is.NoErr(err)
//	is.Equal(2, len(r))
//	is.True(r[0].Success)
//	is.True(r[1].Success)
//}
//
//func TestResolveArmyMovesToOccTerrAndOccTerrMoves3(t *testing.T) {
//	is := is.New(t)
//	g := armyGraph()
//	bud := board.Territory{Abbr: "bud", Name: "Budapest"}
//	gal := board.Territory{Abbr: "gal", Name: "Galicia"}
//	vie := board.Territory{Abbr: "vie", Name: "Vienna"}
//	boh := board.Territory{Abbr: "boh", Name: "Vienna"}
//	ugal := board.Unit{Country: "Russia", Type: board.Army}
//	ubud := board.Unit{Country: "Austria-Hungary", Type: board.Army}
//	uvie := board.Unit{Country: "Austria-Hungary", Type: board.Army}
//	set := order.Set{
//		ArmyGraph: g,
//		Positions: []board.Position{
//			{Territory: bud, Unit: ubud},
//			{Territory: gal, Unit: ugal},
//			{Territory: vie, Unit: uvie},
//		}}
//	set.AddMove(order.Move{Country: "Russia", UnitType: board.Army, From: gal, To: boh})
//	set.AddMove(order.Move{Country: "Austria-Hungary", UnitType: board.Army, From: vie, To: boh})
//	set.AddMove(order.Move{Country: "Austria-Hungary", UnitType: board.Army, From: bud, To: vie})
//	r, err := set.Resolve()
//	is.NoErr(err)
//	is.False(r[0].Success)
//	is.False(r[1].Success)
//	is.False(r[2].Success)
//}
