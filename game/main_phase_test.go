package game_test

import (
	"testing"

	"github.com/cheekybits/is"

	"github.com/burrbd/dip/game"
	"github.com/burrbd/dip/game/order"
	"github.com/burrbd/dip/game/order/board"
)

var (
	bud = board.Territory{Abbr: "bud", Name: "Budapest"}
	gal = board.Territory{Abbr: "gal", Name: "Galicia"}
	vie = board.Territory{Abbr: "vie", Name: "Vienna"}
	boh = board.Territory{Abbr: "boh", Name: "Bohemia"}
	lon = board.Territory{Abbr: "lon", Name: "London"}
)

var handler = game.MainPhaseHandler{ArmyGraph: mockGraph{IsNeighbourFunc: func(t1, t2 string) (bool, error) { return true, nil }}}

type call struct {
	unit *board.Unit
	terr board.Territory
}

type mockPositionMap struct {
	MoveFunc        func(unit *board.Unit, territory board.Territory, strength int)
	HoldFunc        func(unit *board.Unit, strength int)
	BounceFunc      func(unit *board.Unit)
	SetDefeatedFunc func(unit *board.Unit)
	UnitsFunc       func() []*board.Unit
	ConflictFunc    func() []*board.Unit
}

func (m mockPositionMap) Move(unit *board.Unit, territory board.Territory, strength int) {
	m.MoveFunc(unit, territory, strength)
}

func (m mockPositionMap) Hold(unit *board.Unit, strength int) { m.HoldFunc(unit, strength) }

func (m mockPositionMap) Bounce(unit *board.Unit) { m.BounceFunc(unit) }

func (m mockPositionMap) SetDefeated(unit *board.Unit) { m.SetDefeatedFunc(unit) }

func (m mockPositionMap) Units() []*board.Unit { return m.UnitsFunc() }

func (m mockPositionMap) Conflict() []*board.Unit { return m.ConflictFunc() }

func TestMainPhaseHandler_ApplyOrders_HandleMove(t *testing.T) {
	is := is.New(t)
	move := order.Move{UnitType: board.Army, From: board.Territory{Abbr: "bel"}, To: board.Territory{Abbr: "hol"}}
	set := order.Set{Moves: []order.Move{move}}
	u := &board.Unit{PhaseHistory: []board.Position{{Territory: move.From}}}
	act := call{}
	positions := &mockPositionMap{
		MoveFunc: func(unit *board.Unit, territory board.Territory, strength int) {
			act.unit = unit
			act.terr = territory
		},
		UnitsFunc: func() []*board.Unit { return []*board.Unit{u} },
	}
	handler.ApplyOrders(set, positions)
	is.Equal(u, act.unit)
	is.Equal(move.To, act.terr)
}

func TestMainPhaseHandler_ApplyOrders_HandleHoldStrength(t *testing.T) {
	is := is.New(t)

	u1 := &board.Unit{PhaseHistory: []board.Position{{Territory: vie, Cause: board.Added}}}
	u2 := &board.Unit{PhaseHistory: []board.Position{{Territory: bud, Cause: board.Added}}}
	orders := order.Set{}
	h := order.Hold{At: bud}
	orders.AddHold(h)
	orders.AddHoldSupport(order.HoldSupport{Hold: h, By: vie})
	var called bool
	positionMap := &mockPositionMap{
		HoldFunc:  func(unit *board.Unit, strength int) { called = true; is.Equal(1, strength) },
		UnitsFunc: func() []*board.Unit { return []*board.Unit{u1, u2} },
	}

	handler.ApplyOrders(orders, positionMap)
	is.True(called)
}

// newAddedUnit
func newAddedUnit(terr board.Territory, cause int) *board.Unit {
	return nil
}

func TestMainPhaseHandler_ApplyOrders_DoesNotCallMoveWhenNotNeighbour(t *testing.T) {
	is := is.New(t)
	var isNeighbourCalled bool
	notNeighbourHandler := game.MainPhaseHandler{ArmyGraph: mockGraph{
		IsNeighbourFunc: func(t1, t2 string) (bool, error) { isNeighbourCalled = true; return false, nil }},
	}
	set := order.Set{Moves: []order.Move{{}}}
	u := &board.Unit{}
	positions := &mockPositionMap{
		MoveFunc:  func(unit *board.Unit, territory board.Territory, strength int) { is.Fail("unexpected Move() call") },
		UnitsFunc: func() []*board.Unit { return []*board.Unit{u} },
	}
	notNeighbourHandler.ApplyOrders(set, positions)
	is.True(isNeighbourCalled)
}

func TestMainPhaseHandler_ApplyOrders_Strength(t *testing.T) {
	is := is.New(t)

	u1 := &board.Unit{PhaseHistory: []board.Position{{Territory: bud}}}
	u2 := &board.Unit{PhaseHistory: []board.Position{{Territory: vie}}}
	u3 := &board.Unit{PhaseHistory: []board.Position{{Territory: boh}}}
	orders := order.Set{}
	m := order.Move{From: bud, To: gal}
	orders.AddMove(m)
	orders.AddMoveSupport(order.MoveSupport{Move: m, By: vie})
	orders.AddMoveSupport(order.MoveSupport{Move: m, By: boh})

	positionMap := &mockPositionMap{
		MoveFunc:  func(unit *board.Unit, territory board.Territory, strength int) { is.Equal(2, strength) },
		UnitsFunc: func() []*board.Unit { return []*board.Unit{u1, u2, u3} },
	}

	handler.ApplyOrders(orders, positionMap)
}

func TestMainPhaseHandler_ApplyOrders_Strength_WhenSupportIsCut(t *testing.T) {
	// bud -> gal
	// vie s bud -> gal
	// boh -> vie
	is := is.New(t)
	orders := order.Set{}
	u1 := &board.Unit{PhaseHistory: []board.Position{{Territory: bud}}}
	u2 := &board.Unit{PhaseHistory: []board.Position{{Territory: vie}}}
	u3 := &board.Unit{PhaseHistory: []board.Position{{Territory: boh}}}
	move := order.Move{From: bud, To: gal}
	orders.AddMove(move)
	orders.AddMoveSupport(order.MoveSupport{Move: move, By: vie})
	orders.AddMove(order.Move{From: boh, To: vie})

	positions := &mockPositionMap{
		MoveFunc:  func(unit *board.Unit, territory board.Territory, strength int) {},
		UnitsFunc: func() []*board.Unit { return []*board.Unit{u1, u2, u3} },
	}

	handler.ApplyOrders(orders, positions)

	is.Equal(0, u1.Position().Strength)
}

func TestMainPhaseHandler_ApplyOrders_Strength_WhenSupportIsCutByAttackedUnit(t *testing.T) {
	// bud -> gal
	// vie s bud -> gal
	// gal -> vie

	// gal can't cut support because the support is for attack against gal
	is := is.New(t)
	u1 := &board.Unit{PhaseHistory: []board.Position{{Territory: bud}}}
	u2 := &board.Unit{PhaseHistory: []board.Position{{Territory: vie}}}
	u3 := &board.Unit{PhaseHistory: []board.Position{{Territory: gal}}}

	orders := order.Set{}
	move := order.Move{From: bud, To: gal}
	orders.AddMove(move)
	orders.AddMoveSupport(order.MoveSupport{Move: move, By: vie})
	orders.AddMove(order.Move{From: gal, To: vie})

	moveStrength := map[*board.Unit]int{}
	positions := &mockPositionMap{
		MoveFunc:  func(unit *board.Unit, territory board.Territory, strength int) { moveStrength[unit] = strength },
		UnitsFunc: func() []*board.Unit { return []*board.Unit{u1, u2, u3} },
	}

	handler.ApplyOrders(orders, positions)

	is.Equal(2, len(moveStrength))
	is.Equal(1, moveStrength[u1])
	is.Equal(0, moveStrength[u3])
}
