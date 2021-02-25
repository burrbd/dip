package order_test

import (
	"testing"

	"github.com/burrbd/dip/game/order"
	"github.com/burrbd/dip/game/order/board"
	"github.com/cheekybits/is"
)

func TestValidator_ValidateMove(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		is := is.New(t)
		v := order.NewValidator(&mockSimpleGraph{
			hasEdgeBetweenFunc: func(_, _ int64) bool {
				return true
			},
		})
		u := board.Unit{Country: "fr"}
		m := order.Move{Country: "fr"}
		is.NoErr(v.ValidateMove(u, m))
	})

	t.Run("country mismatch", func(t *testing.T) {
		is := is.New(t)
		v := order.NewValidator(&mockSimpleGraph{})
		u := board.Unit{Country: "fr"}
		m := order.Move{Country: "bogus"}
		is.Err(v.ValidateMove(u, m))
	})

	t.Run("no edge between territories", func(t *testing.T) {
		is := is.New(t)
		g := &mockSimpleGraph{
			hasEdgeBetweenFunc: func(xid, yid int64) bool {
				return false
			},
		}
		v := order.NewValidator(g)
		u := board.Unit{Country: "fr"}
		m := order.Move{Country: "fr", From: board.Territory{}, To: board.Territory{}}
		is.Err(v.ValidateMove(u, m))
		is.True(g.called)
	})
}

func TestValidator_ValidateMoveSupport(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		is := is.New(t)
		v := order.NewValidator(&mockSimpleGraph{
			hasEdgeBetweenFunc: func(_, _ int64) bool {
				return true
			},
		})
		u := board.Unit{Country: "fr"}
		m := order.MoveSupport{Country: "fr"}
		is.NoErr(v.ValidateMoveSupport(u, m))
	})

	t.Run("country mismatch", func(t *testing.T) {
		is := is.New(t)
		v := order.NewValidator(&mockSimpleGraph{})

		u := board.Unit{Country: "fr"}
		s := order.MoveSupport{Country: "bougus"}

		is.Err(v.ValidateMoveSupport(u, s))
	})

	t.Run("no territory edge to support", func(t *testing.T) {
		is := is.New(t)
		g := &mockSimpleGraph{
			hasEdgeBetweenFunc: func(xid, yid int64) bool {
				return false
			},
		}
		v := order.NewValidator(g)
		u := board.Unit{}
		s := order.MoveSupport{
			By:   board.Territory{},
			Move: order.Move{From: board.Territory{}, To: board.Territory{}},
		}

		is.Err(v.ValidateMoveSupport(u, s))
		is.True(g.called)
	})
}

type mockSimpleGraph struct {
	called             bool
	hasEdgeBetweenFunc func(xid, yid int64) bool
}

func (g *mockSimpleGraph) HasEdgeBetween(xid, yid int64) bool {
	g.called = true
	return g.hasEdgeBetweenFunc(xid, yid)
}
