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
	_ = g.AddVertex(gal)
	_ = g.AddVertex(bud)
	_ = g.AddVertex(vie)
	_ = g.AddEdge(gal, bud)
	_ = g.AddEdge(gal, vie)
	_ = g.AddEdge(bud, vie)
	return g
}

func TestResolveArmyMovesSuccess(t *testing.T) {
	is := is.New(t)
	g := armyGraph()
	bud := board.Territory{Abbr: "bud", Name: "Budapest"}
	gal := board.Territory{Abbr: "gal", Name: "Galicia"}
	unit := board.Unit{Country: "France", Type: board.Army}
	set := order.Set{ArmyGraph: g, Positions: []board.Position{{Territory: bud, Unit: unit}}}
	set.AddMove(order.Move{Country: "France", UnitType: board.Army, From: bud, To: gal})
	r, err := set.Resolve()
	is.NoErr(err)
	is.True(r[0].Success)
}

func TestResolveArmyMovesWrongCountry(t *testing.T) {
	is := is.New(t)
	g := armyGraph()
	bud := board.Territory{Abbr: "bud", Name: "Budapest"}
	gal := board.Territory{Abbr: "gal", Name: "Galicia"}
	unit := board.Unit{Country: "France", Type: board.Army}
	set := order.Set{ArmyGraph: g, Positions: []board.Position{{Territory: bud, Unit: unit}}}
	set.AddMove(order.Move{Country: "France", UnitType: board.Army, From: bud, To: gal})
	r, err := set.Resolve()
	is.NoErr(err)
	is.True(r[0].Success)
}

func TestResolveArmyMovesToOccupiedSpace(t *testing.T) {
	is := is.New(t)
	g := armyGraph()
	bud := board.Territory{Abbr: "bud", Name: "Budapest"}
	gal := board.Territory{Abbr: "gal", Name: "Galicia"}
	ubud := board.Unit{Country: "France", Type: board.Army}
	ugal := board.Unit{Country: "Austria-Hungary", Type: board.Army}
	set := order.Set{
		ArmyGraph: g,
		Positions: []board.Position{
			{Territory: bud, Unit: ubud},
			{Territory: gal, Unit: ugal},
		}}
	set.AddMove(order.Move{Country: "France", UnitType: board.Army, From: bud, To: gal})
	r, err := set.Resolve()
	is.NoErr(err)
	is.False(r[0].Success)
}

func TestResolveArmyMovesToOccTerrAndOccTerrMoves(t *testing.T) {
	is := is.New(t)
	g := armyGraph()
	bud := board.Territory{Abbr: "bud", Name: "Budapest"}
	gal := board.Territory{Abbr: "gal", Name: "Galicia"}
	vie := board.Territory{Abbr: "vie", Name: "Vienna"}
	ubud := board.Unit{Country: "France", Type: board.Army}
	ugal := board.Unit{Country: "Austria-Hungary", Type: board.Army}
	set := order.Set{
		ArmyGraph: g,
		Positions: []board.Position{
			{Territory: bud, Unit: ubud},
			{Territory: gal, Unit: ugal},
		}}
	set.AddMove(order.Move{Country: "France", UnitType: board.Army, From: bud, To: gal})
	set.AddMove(order.Move{Country: "Austria-Hungary", UnitType: board.Army, From: gal, To: vie})
	r, err := set.Resolve()
	is.NoErr(err)
	is.True(r[0].Success)
	is.True(r[1].Success)
}

func TestResolveArmyMovesToOccTerrAndOccTerrMoves2(t *testing.T) {
	is := is.New(t)
	g := armyGraph()
	bud := board.Territory{Abbr: "bud", Name: "Budapest"}
	gal := board.Territory{Abbr: "gal", Name: "Galicia"}
	vie := board.Territory{Abbr: "vie", Name: "Vienna"}
	ubud := board.Unit{Country: "France", Type: board.Army}
	ugal := board.Unit{Country: "France", Type: board.Army}
	set := order.Set{
		ArmyGraph: g,
		Positions: []board.Position{
			{Territory: bud, Unit: ubud},
			{Territory: gal, Unit: ugal},
		}}
	set.AddMove(order.Move{Country: "France", UnitType: board.Army, From: bud, To: gal})
	set.AddMove(order.Move{Country: "France", UnitType: board.Army, From: gal, To: vie})
	r, err := set.Resolve()
	is.NoErr(err)
	is.True(r[0].Success)
	is.True(r[1].Success)
}
