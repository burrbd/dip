package game_test

import (
	"testing"

	"github.com/burrbd/diplomacy/internal/game"
	"github.com/burrbd/diplomacy/internal/game/order"
	"github.com/burrbd/diplomacy/internal/game/order/board"

	"github.com/cheekybits/is"
)

var (
	par     = board.Territory{Abbr: "par", Name: "Paris"}
	handler = game.MainPhaseHandler{
		ArmyGraph: mockGraph{
			IsNeighbourFunc: func(t1, t2 string) (bool, error) { return true, nil },
		},
	}
)

type call struct {
	unit *board.Unit
	terr board.Territory
}

var calls = make([]call, 0)

type mockPositionMap struct {
	MoveFunc        func(unit *board.Unit, territory board.Territory, strength int)
	BounceFunc      func(unit *board.Unit)
	SetDefeatedFunc func(unit *board.Unit)
	UnitsFunc       func() []*board.Unit
	ConflictFunc    func() []*board.Unit
}

func (m mockPositionMap) Move(unit *board.Unit, territory board.Territory, strength int) {
	m.MoveFunc(unit, territory, strength)
}

func (m mockPositionMap) Bounce(unit *board.Unit) {
	m.BounceFunc(unit)
}

func (m mockPositionMap) SetDefeated(unit *board.Unit) {
	m.SetDefeatedFunc(unit)
}

func (m mockPositionMap) Units() []*board.Unit {
	return m.UnitsFunc()
}

func (m mockPositionMap) Conflict() []*board.Unit {
	return m.ConflictFunc()
}

func TestOrderHandler_HandleMove(t *testing.T) {
	is := is.New(t)

	h := game.MainPhaseHandler{

		ArmyGraph: mockGraph{
			IsNeighbourFunc: func(t1, t2 string) (bool, error) { return true, nil },
		},
	}

	set := order.Set{}

	m1, err := order.Decode("A Bel-Hol")
	is.NoErr(err)
	mm1 := m1.(order.Move)
	set.AddMove(mm1)
	u1 := &board.Unit{Territory: mm1.From}

	act := call{}
	positions := &mockPositionMap{
		MoveFunc: func(unit *board.Unit, territory board.Territory, strength int) {
			act.unit = unit
			act.terr = territory
		},
		UnitsFunc: func() []*board.Unit { return []*board.Unit{u1} },
	}

	h.ApplyOrders(set, positions)
	is.Equal(u1, act.unit)
	is.Equal(mm1.To, act.terr)
}

func TestOrderHandler_Handle_NotNeighbor_DoesNotCallMove(t *testing.T) {
	is := is.New(t)

	var isNeighbourCalled bool

	h := game.MainPhaseHandler{
		ArmyGraph: mockGraph{
			IsNeighbourFunc: func(t1, t2 string) (bool, error) { isNeighbourCalled = true; return false, nil },
		},
	}

	set := order.Set{}

	m1, err := order.Decode("A Bel-Hol")
	is.NoErr(err)
	mm1 := m1.(order.Move)
	set.AddMove(mm1)

	u1 := &board.Unit{Territory: mm1.From}

	positions := &mockPositionMap{
		MoveFunc:  func(unit *board.Unit, territory board.Territory, strength int) { is.Fail("unexpected Move() call") },
		UnitsFunc: func() []*board.Unit { return []*board.Unit{u1} },
	}

	h.ApplyOrders(set, positions)

	is.True(isNeighbourCalled)
}

func TestSet_Strength(t *testing.T) {
	is := is.New(t)

	u1 := &board.Unit{Territory: bud}
	u2 := &board.Unit{Territory: vie}
	u3 := &board.Unit{Territory: boh}

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

func TestSet_Strength_WhenSupportIsCut(t *testing.T) {
	// bud -> gal
	// vie s bud -> gal
	// boh -> vie
	is := is.New(t)
	orders := order.Set{}
	u1 := &board.Unit{Territory: bud}
	u2 := &board.Unit{Territory: vie}
	u3 := &board.Unit{Territory: boh}
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

func TestSet_Strength_WhenSupportIsCutByAttackedUnit(t *testing.T) {
	// bud -> gal
	// vie s bud -> gal
	// gal -> vie

	// gal can't cut support because the support is for attack against gal
	is := is.New(t)
	u1 := &board.Unit{Territory: bud}
	u2 := &board.Unit{Territory: vie}
	u3 := &board.Unit{Territory: gal}

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
