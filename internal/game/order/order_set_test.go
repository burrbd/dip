package order_test

import (
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
)

func newPositions() board.Positions {
	return board.NewPositions([]board.Territory{bud, gal, vie, boh})

}

func TestSet_MoveSupportCount(t *testing.T) {
	is := is.New(t)
	u1 := &board.Unit{Position: gal}
	u2 := &board.Unit{Position: vie}
	u3 := &board.Unit{Position: boh}

	positions := newPositions()
	u1.PrevPosition = &bud
	positions.Add(gal, u1)
	positions.Add(vie, u2)
	positions.Add(boh, u3)

	orders := order.Set{}
	m := order.Move{From: bud, To: gal}
	orders.AddMove(m)
	orders.AddMoveSupport(order.MoveSupport{Move: m, By: vie})
	is.Equal(1, orders.Strength(u1))
	orders.AddMoveSupport(order.MoveSupport{Move: m, By: boh})
	is.Equal(2, orders.Strength(u1))
}
