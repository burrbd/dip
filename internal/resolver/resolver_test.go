package resolver_test

import (
	"testing"

	"github.com/cheekybits/is"

	"github.com/burrbd/diplomacy/internal/resolver/order"
	"github.com/burrbd/diplomacy/internal/resolver/order/board"
	"github.com/burrbd/kit/graph"
)

func armyGraph() *graph.Simple {
	g := graph.NewSimple()
	bud := board.Territory{Abbr: "bud", Name: "Budapest"}
	gal := board.Territory{Abbr: "gal", Name: "Galicia"}
	vie := board.Territory{Abbr: "vie", Name: "Vienna"}
	boh := board.Territory{Abbr: "boh", Name: "Vienna"}
	_ = g.AddVertex(gal)
	_ = g.AddVertex(bud)
	_ = g.AddVertex(vie)
	_ = g.AddVertex(boh)
	_ = g.AddEdge(gal, bud)
	_ = g.AddEdge(gal, vie)
	_ = g.AddEdge(bud, vie)
	_ = g.AddEdge(boh, vie)
	_ = g.AddEdge(boh, gal)
	return g
}

func TestResolveArmyMovesSuccess(t *testing.T) {
	is := is.New(t)
	g := armyGraph()
	bud := board.Territory{Abbr: "bud", Name: "Budapest"}
	gal := board.Territory{Abbr: "gal", Name: "Galicia"}
	unit := board.Unit{Country: "France", Type: board.Army}
	set := order.Set{ArmyGraph: g, Positions: []*board.Position{
		{Territory: bud, Unit: unit, Previous: make([]board.Territory, 0)},
	}}
	set.AddMove(order.Move{Country: "France", UnitType: board.Army, From: bud, To: gal})
	r, err := set.Resolve()
	is.NoErr(err)
	is.True(r[0].Success)
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
