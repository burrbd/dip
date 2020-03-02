package game_test

import (
	"sort"
	"testing"

	"github.com/burrbd/dip/game"
	"github.com/burrbd/dip/game/order/board"
	"github.com/cheekybits/is"
)

func TestUnitsByStrength(t *testing.T) {
	is := is.New(t)
	m := board.NewPositionManager()
	u1 := &board.Unit{}
	u2 := &board.Unit{}
	u3 := &board.Unit{}
	units := []*board.Unit{u1, u2, u3}
	m.AddUnit(u1, board.Territory{})
	m.AddUnit(u2, board.Territory{})
	m.AddUnit(u3, board.Territory{})
	m.Hold(u1, 3)
	m.Hold(u2, 2)
	m.Hold(u3, 1)

	sort.Sort(game.UnitPositionsByStrength(m, units))

	is.Equal(u1, units[0])
	is.Equal(u2, units[1])
	is.Equal(u3, units[2])
}

func TestNilUnitsCoverage(t *testing.T) {
	is := is.New(t)
	m := mockManager{positionFunc: func(u *board.Unit) *board.Position {
		if u == nil {
			return nil
		}
		return &board.Position{}
	}}

	t.Run("nil i sort value", func(t *testing.T) {
		s := game.UnitPositionsByStrength(m, []*board.Unit{nil, &board.Unit{}})
		is.False(s.Less(0, 1))
	})

	t.Run("nil j sort value", func(t *testing.T) {
		s := game.UnitPositionsByStrength(m, []*board.Unit{&board.Unit{}, nil})
		is.True(s.Less(0, 1))
	})
}

type mockManager struct {
	board.Manager
	positionFunc func(u *board.Unit) *board.Position
}

func (m mockManager) Position(u *board.Unit) *board.Position {
	return m.positionFunc(u)
}
