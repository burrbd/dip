package order_test

import (
	"sort"
	"testing"

	"github.com/burrbd/diplomacy/internal/game/order"
	"github.com/burrbd/diplomacy/internal/game/order/board"

	"github.com/cheekybits/is"
)

var (
	bud = board.Territory{Abbr: "bud", Name: "Budapest"}
	gal = board.Territory{Abbr: "gal", Name: "Galicia"}
	vie = board.Territory{Abbr: "vie", Name: "Vienna"}
	boh = board.Territory{Abbr: "boh", Name: "Bohemia"}
	lon = board.Territory{Abbr: "lon", Name: "London"}
	par = board.Territory{Abbr: "par", Name: "Paris"}
)

func TestSet_Strength(t *testing.T) {
	is := is.New(t)
	u1 := &board.Unit{Position: gal}
	u2 := &board.Unit{Position: vie}
	u3 := &board.Unit{Position: boh}

	positions := board.NewPositions()
	u1.PrevPosition = &bud
	positions.Add(u1)
	positions.Add(u2)
	positions.Add(u3)

	orders := order.Set{}
	m := order.Move{From: bud, To: gal}
	orders.AddMove(m)
	orders.AddMoveSupport(order.MoveSupport{Move: m, By: vie})
	is.Equal(1, orders.Strength(u1))
	orders.AddMoveSupport(order.MoveSupport{Move: m, By: boh})
	is.Equal(2, orders.Strength(u1))
}

func TestSet_Strength_WhenSupportIsCut(t *testing.T) {
	// bud -> gal
	// vie s bud -> gal
	// boh -> vie
	is := is.New(t)
	orders := order.Set{}
	move := order.Move{From: bud, To: gal}
	orders.AddMove(move)
	orders.AddMoveSupport(order.MoveSupport{Move: move, By: vie})
	orders.AddMove(order.Move{From: boh, To: vie})
	is.Equal(0, orders.Strength(&board.Unit{Position: gal, PrevPosition: &bud}))
}

func TestSet_Strength_WhenSupportIsCutByAttackedUnit(t *testing.T) {
	// bud -> gal
	// vie s bud -> gal
	// gal -> vie

	// gal can't cut support because the support is for attack against gal
	is := is.New(t)
	orders := order.Set{}
	move := order.Move{From: bud, To: gal}
	orders.AddMove(move)
	orders.AddMoveSupport(order.MoveSupport{Move: move, By: vie})
	orders.AddMove(order.Move{From: gal, To: vie})
	is.Equal(1, orders.Strength(&board.Unit{Position: gal, PrevPosition: &bud}))
}

func TestSet_Strength_WhenUnitIsAlreadyMovedToOrigin_ReturnsZero(t *testing.T) {
	is := is.New(t)
	orders := order.Set{}
	move := order.Move{From: bud, To: gal}
	orders.AddMove(move)
	orders.AddMoveSupport(order.MoveSupport{Move: move, By: vie})
	is.Equal(0, orders.Strength(&board.Unit{Position: bud, PrevPosition: &bud}))
}

func TestSet_ByStrength(t *testing.T) {
	// A Bud-Gal
	// A Vie S Bud-Gal
	// A Par-Lon
	is := is.New(t)
	orders := order.Set{}
	supportedMove := order.Move{From: bud, To: gal}
	orders.AddMove(order.Move{From: bud, To: gal})
	u1 := &board.Unit{Position: gal, PrevPosition: &bud}
	orders.AddMoveSupport(order.MoveSupport{Move: supportedMove, By: vie})
	u2 := &board.Unit{Position: vie}
	orders.AddMove(order.Move{From: par, To: lon})
	u3 := &board.Unit{Position: lon, PrevPosition: &par}

	units := []*board.Unit{u3, u1, u2}
	sort.Sort(orders.ByStrength(units))

	is.Equal(u1, units[0])
}
