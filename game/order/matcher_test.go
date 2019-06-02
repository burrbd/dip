package order_test

import (
	"testing"

	"github.com/burrbd/dip/game/order"
	"github.com/burrbd/dip/game/order/board"
)

func TestPositionMatcher_MatchMove(t *testing.T) {
	specs := []struct {
		desc         string
		from         board.Territory
		to           board.Territory
		isNeighbor   bool
		posTerritory board.Territory
		posCause     board.PositionEvent
		moveCountry  string
		country      string
		match        bool
		graphInvoked bool
	}{
		{
			desc:         "success",
			from:         board.Territory{Abbr: "a"},
			to:           board.Territory{Abbr: "b"},
			isNeighbor:   true,
			posTerritory: board.Territory{Abbr: "a"},
			posCause:     board.Added,
			moveCountry:  "fr",
			country:      "fr",
			match:        true,
			graphInvoked: true,
		},
		{
			desc:         "country mismatch",
			from:         board.Territory{Abbr: "a"},
			to:           board.Territory{Abbr: "b"},
			posTerritory: board.Territory{Abbr: "a"},
			posCause:     board.Added,
			moveCountry:  "de",
			country:      "fr",
			match:        false,
			graphInvoked: false,
		},
		{
			desc:         "position mismatch",
			from:         board.Territory{Abbr: "a"},
			to:           board.Territory{Abbr: "b"},
			posTerritory: board.Territory{Abbr: "zzzz"},
			posCause:     board.Added,
			moveCountry:  "de",
			country:      "de",
			match:        false,
			graphInvoked: true,
		},
		{
			desc:         "not neighbor",
			from:         board.Territory{Abbr: "a"},
			to:           board.Territory{Abbr: "b"},
			isNeighbor:   false,
			posTerritory: board.Territory{Abbr: "a"},
			posCause:     board.Added,
			moveCountry:  "fr",
			country:      "fr",
			match:        false,
			graphInvoked: true,
		},
		{
			desc:         "not added",
			from:         board.Territory{Abbr: "a"},
			to:           board.Territory{Abbr: "b"},
			isNeighbor:   true,
			posTerritory: board.Territory{Abbr: "a"},
			posCause:     board.Moved,
			moveCountry:  "fr",
			country:      "fr",
			match:        false,
			graphInvoked: true,
		},
	}

	for i, spec := range specs {
		t.Run(spec.desc, func(t *testing.T) {
			graph := &mockGraph{isNeighbourFunc: func(_, _ string) (bool, error) {
				return spec.isNeighbor, nil
			}}
			m := order.PositionMatcher{ArmyGraph: graph}
			move := order.Move{From: spec.from, To: spec.to, Country: spec.moveCountry}
			pos := board.Position{Territory: spec.posTerritory, Cause: spec.posCause}
			if m.MatchMove(move, pos, spec.country) != spec.match {
				t.Errorf("[%d] unexpected match move result: match != %t", i, spec.match)
			}
			if !graph.invoked == spec.graphInvoked {
				t.Errorf("[%d] expected graph to be invoked", i)
			}
		})
	}
}

func TestPositionMatcher_MatchHold(t *testing.T) {
	specs := []struct {
		desc         string
		at           board.Territory
		posTerritory board.Territory
		posCause     board.PositionEvent
		holdCountry  string
		country      string
		match        bool
	}{
		{
			desc:         "success",
			at:           board.Territory{Abbr: "a"},
			posTerritory: board.Territory{Abbr: "a"},
			posCause:     board.Added,
			holdCountry:  "fr",
			country:      "fr",
			match:        true,
		},
		{
			desc:         "country mismatch",
			at:           board.Territory{Abbr: "a"},
			posTerritory: board.Territory{Abbr: "a"},
			posCause:     board.Added,
			holdCountry:  "fr",
			country:      "de",
			match:        false,
		},
		{
			desc:         "position mismatch",
			at:           board.Territory{Abbr: "a"},
			posTerritory: board.Territory{Abbr: "b"},
			posCause:     board.Added,
			holdCountry:  "fr",
			country:      "fr",
			match:        false,
		},
		{
			desc:         "not added",
			at:           board.Territory{Abbr: "a"},
			posTerritory: board.Territory{Abbr: "a"},
			posCause:     board.Moved,
			holdCountry:  "fr",
			country:      "fr",
			match:        false,
		},
	}

	for i, spec := range specs {
		t.Run(spec.desc, func(t *testing.T) {
			m := order.PositionMatcher{}
			hold := order.Hold{At: spec.at, Country: spec.holdCountry}
			pos := board.Position{Territory: spec.posTerritory, Cause: spec.posCause}
			if m.MatchHold(hold, pos, spec.country) != spec.match {
				t.Errorf("[%d] unexpected match hold result: match != %t", i, spec.match)
			}
		})
	}
}

type mockGraph struct {
	invoked         bool
	isNeighbourFunc func(t1, t2 string) (bool, error)
}

func (g *mockGraph) IsNeighbour(t1, t2 string) (bool, error) {
	g.invoked = true
	return g.isNeighbourFunc(t1, t2)
}
