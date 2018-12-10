package game_test

import (
	"testing"

	"github.com/burrbd/diplomacy/internal/game"
	"github.com/burrbd/diplomacy/internal/game/order"
	"github.com/burrbd/diplomacy/internal/game/order/board"

	"github.com/cheekybits/is"
)

type call struct {
	unit *board.Unit
	terr board.Territory
}

var calls = make([]call, 0)

type mockPositionMap struct {
	MoveFunc     func(unit *board.Unit, territory board.Territory)
	BounceFunc   func(unit *board.Unit, territory board.Territory)
	AllUnitsFunc func() []*board.Unit
}

func (m mockPositionMap) Move(unit *board.Unit, territory board.Territory) {
	m.MoveFunc(unit, territory)
}

func (m mockPositionMap) Bounce(unit *board.Unit, territory board.Territory) {
	m.BounceFunc(unit, territory)
}

func (m mockPositionMap) AllUnits() []*board.Unit {
	return m.AllUnitsFunc()
}

func TestOrderHandler_HandleMove(t *testing.T) {
	is := is.New(t)

	handler := game.OrderHandler{

		ArmyGraph: mockGraph{
			IsNeighbourFunc: func(t1, t2 string) (bool, error) { return true, nil },
		},
	}

	set := order.Set{}

	m1, err := order.Decode("A Bel-Hol")
	is.NoErr(err)
	mm1 := m1.(order.Move)
	set.AddMove(mm1)
	u1 := &board.Unit{Position: mm1.From}

	act := call{}
	recorder := &mockPositionMap{
		MoveFunc: func(unit *board.Unit, territory board.Territory) {
			act.unit = unit
			act.terr = territory
		},
		AllUnitsFunc: func() []*board.Unit { return []*board.Unit{u1} },
	}

	handler.Handle(set, recorder)
	is.Equal(u1, act.unit)
	is.Equal(mm1.To, act.terr)
}

func TestOrderHandler_Handle_NotNeighbor_DoesNotCallMove(t *testing.T) {
	is := is.New(t)

	var isNeighbourCalled bool

	handler := game.OrderHandler{
		ArmyGraph: mockGraph{
			IsNeighbourFunc: func(t1, t2 string) (bool, error) { isNeighbourCalled = true; return false, nil },
		},
	}

	set := order.Set{}

	m1, err := order.Decode("A Bel-Hol")
	is.NoErr(err)
	mm1 := m1.(order.Move)
	set.AddMove(mm1)

	u1 := &board.Unit{Position: mm1.From}

	recorder := &mockPositionMap{
		MoveFunc:     func(unit *board.Unit, territory board.Territory) { is.Fail("unexpected Move() call") },
		AllUnitsFunc: func() []*board.Unit { return []*board.Unit{u1} },
	}

	handler.Handle(set, recorder)

	is.True(isNeighbourCalled)
}

var (
	par     = board.Territory{Abbr: "par", Name: "Paris"}
	handler = game.OrderHandler{
		ArmyGraph: mockGraph{
			IsNeighbourFunc: func(t1, t2 string) (bool, error) { return true, nil },
		},
	}
)

func TestSet_Strength(t *testing.T) {
	is := is.New(t)

	u1 := &board.Unit{Position: bud}
	u2 := &board.Unit{Position: vie}
	u3 := &board.Unit{Position: boh}

	positions := board.NewPositionMap()
	positions.Add(u1)
	positions.Add(u2)
	positions.Add(u3)

	orders := order.Set{}
	m := order.Move{From: bud, To: gal}
	orders.AddMove(m)
	orders.AddMoveSupport(order.MoveSupport{Move: m, By: vie})
	orders.AddMoveSupport(order.MoveSupport{Move: m, By: boh})

	recorder := &mockPositionMap{
		MoveFunc:     func(unit *board.Unit, territory board.Territory) {},
		AllUnitsFunc: func() []*board.Unit { return []*board.Unit{u1, u2, u3} },
	}

	handler.Handle(orders, recorder)
	is.Equal(2, u1.Strength)
}

func TestSet_Strength_WhenSupportIsCut(t *testing.T) {
	// bud -> gal
	// vie s bud -> gal
	// boh -> vie
	is := is.New(t)
	orders := order.Set{}
	u1 := &board.Unit{Position: bud}
	u2 := &board.Unit{Position: vie}
	u3 := &board.Unit{Position: boh}
	move := order.Move{From: bud, To: gal}
	orders.AddMove(move)
	orders.AddMoveSupport(order.MoveSupport{Move: move, By: vie})
	orders.AddMove(order.Move{From: boh, To: vie})

	recorder := &mockPositionMap{
		MoveFunc:     func(unit *board.Unit, territory board.Territory) {},
		AllUnitsFunc: func() []*board.Unit { return []*board.Unit{u1, u2, u3} },
	}

	handler.Handle(orders, recorder)

	is.Equal(0, u1.Strength)
}

func TestSet_Strength_WhenSupportIsCutByAttackedUnit(t *testing.T) {
	// bud -> gal
	// vie s bud -> gal
	// gal -> vie

	// gal can't cut support because the support is for attack against gal
	is := is.New(t)
	u1 := &board.Unit{Position: bud}
	u2 := &board.Unit{Position: vie}
	u3 := &board.Unit{Position: gal}

	orders := order.Set{}
	move := order.Move{From: bud, To: gal}
	orders.AddMove(move)
	orders.AddMoveSupport(order.MoveSupport{Move: move, By: vie})
	orders.AddMove(order.Move{From: gal, To: vie})

	recorder := &mockPositionMap{
		MoveFunc:     func(unit *board.Unit, territory board.Territory) {},
		AllUnitsFunc: func() []*board.Unit { return []*board.Unit{u1, u2, u3} },
	}

	handler.Handle(orders, recorder)

	is.Equal(1, u1.Strength)
}
