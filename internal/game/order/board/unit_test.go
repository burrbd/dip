package board_test

import (
	"sort"
	"testing"

	"github.com/cheekybits/is"

	"github.com/burrbd/diplomacy/internal/game/order/board"
)

func TestUnitsByStrength(t *testing.T) {
	// A Bud-Gal
	// A Vie S Bud-Gal
	// A Par-Lon
	is := is.New(t)
	u1 := &board.Unit{Strength: 3}
	u2 := &board.Unit{Strength: 2}
	u3 := &board.Unit{Strength: 1}

	units := []*board.Unit{u2, u1, u3}
	sort.Sort(board.UnitsByStrength(units))

	is.Equal(u1, units[0])
	is.Equal(u2, units[1])
	is.Equal(u3, units[2])
}
