package board_test

import (
	"testing"

	"github.com/burrbd/dip/game/order/board"
	"github.com/cheekybits/is"
)

func TestCreateArmyGraph(t *testing.T) {
	lu := board.LookupTerritory
	is := is.New(t)

	g := board.CreateArmyGraph()

	is.True(g.HasEdgeBetween(lu("alb").ID(), lu("gre").ID()))
	is.True(g.HasEdgeBetween(lu("gre").ID(), lu("alb").ID()))

	is.False(g.HasEdgeBetween(lu("mar").ID(), lu("mos").ID()))
}
