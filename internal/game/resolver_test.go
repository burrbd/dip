package game_test

import (
	"testing"

	"github.com/burrbd/diplomacy/internal/game"

	"github.com/cheekybits/is"

	"github.com/burrbd/diplomacy/internal/game/order"
	"github.com/burrbd/diplomacy/internal/game/order/board"
)

var cases = []orderCase{
	{
		description: "given a simple move, then the order is resolved",
		givenMap:    []string{"bud", "vie"},
		orders: []orderResult{
			{order: "A Bud-Vie", result: "vie"},
		},
	},
	{
		description: "given two units attack same territory without support, then both units bounce back",
		givenMap:    []string{"gal", "bud", "vie"},
		orders: []orderResult{
			{order: "A Bud-Vie", result: "bud"},
			{order: "A Gal-Vie", result: "gal"},
		},
	},
	{
		description: "given units attack in circular chain without support, then all attacking units bounce back",
		givenMap:    []string{"gal", "bud", "vie"},
		orders: []orderResult{
			{order: "A Bud-Gal", result: "bud"},
			{order: "A Gal-Vie", result: "gal"},
			{order: "A Vie H", result: "vie"},
		},
	},
	{
		description: "given unit attacks unsupported territory with support, then attacking unit wins",
		givenMap:    []string{"gal", "vie", "boh", "bud"},
		orders: []orderResult{
			{order: "A Gal-Vie", result: "vie"},
			{order: "A Boh S A Gal-Vie", result: "boh"},
			{order: "A Bud-Vie", result: "bud"},
		},
	},
	{
		description: "given unit holds territory and is not dislodged, then unit remains on territory",
		givenMap:    []string{"vie"},
		orders: []orderResult{
			{order: "A Vie H", result: "vie"},
		},
	},
}

func TestMainPhaseResolver_Resolve_OnlyMovesToNeighbouringTerritory(t *testing.T) {
	is := is.New(t)
	graph := mockGraph{
		IsNeighbourFunc: func(t1, t2 string) (bool, error) { return false, nil },
	}
	resolver := game.MainPhaseResolver{ArmyGraph: graph}

	unit := &board.Unit{Position: gal}

	positions := newPositions()
	positions.Add(unit)

	orders := order.Set{}
	orders.AddMove(order.Move{From: gal, To: lon})

	resolved, err := resolver.Resolve(orders, positions)

	is.NoErr(err)
	is.Nil(resolved.Units["lon"])
	is.Equal(unit, resolved.Units["gal"][0])
}

type orderResult struct {
	order   string
	result  string
	retreat bool
}

type orderCase struct {
	description string
	givenMap    []string
	orders      []orderResult
	focus       bool
}

func TestMainPhaseResolver_ResolveCases(t *testing.T) {
	is := is.New(t)
	graph := mockGraph{
		IsNeighbourFunc: func(t1, t2 string) (bool, error) { return true, nil },
	}

	focused := focusedCases(cases)
	if len(focused) > 0 {
		cases = focused
	}

	for _, c := range cases {
		t.Log(c.description)
		territories := make([]board.Territory, 0)
		for _, t := range c.givenMap {
			territories = append(territories, board.Territory{Abbr: t})
		}
		positions := board.NewPositions(territories)

		orders := order.Set{}
		resolver := game.MainPhaseResolver{ArmyGraph: graph}

		units := make([]*board.Unit, 0)
		for _, c := range c.orders {
			t.Logf("\t%s", c.order)
			o, _ := order.Decode(c.order)
			var t board.Territory
			switch v := o.(type) {
			case order.Move:
				t = v.From
				orders.AddMove(v)
			case order.Hold:
				t = v.Pos
				orders.AddHold(v)
			case order.MoveSupport:
				t = v.By
				orders.AddMoveSupport(v)
			case order.MoveConvoy:
				t = v.By
				orders.AddMoveConvoy(v)
			}
			u := &board.Unit{Position: t}
			units = append(units, u)
			positions.Add(u)
		}

		resolvedPositions, err := resolver.Resolve(orders, positions)
		is.NoErr(err)

		for i, u := range units {
			terr := c.orders[i].result
			is.Equal(u, resolvedPositions.Units[terr][0])
			is.Equal(1, len(resolvedPositions.Units[terr]))
		}

		positionTotal := 0
		for _, positionUnits := range resolvedPositions.Units {
			positionTotal += len(positionUnits)
		}
		is.Equal(positionTotal, len(units))
	}
}

func focusedCases(cases []orderCase) []orderCase {
	focused := make([]orderCase, 0)
	for _, c := range cases {
		if c.focus {
			focused = append(focused, c)
		}
	}
	return focused
}

var (
	bud = board.Territory{Abbr: "bud", Name: "Budapest"}
	gal = board.Territory{Abbr: "gal", Name: "Galicia"}
	vie = board.Territory{Abbr: "vie", Name: "Vienna"}
	boh = board.Territory{Abbr: "boh", Name: "Bohemia"}
	lon = board.Territory{Abbr: "lon", Name: "London"}
)

type mockGraph struct {
	IsNeighbourFunc func(t1, t2 string) (bool, error)
}

func (g mockGraph) IsNeighbour(t1, t2 string) (bool, error) {
	return g.IsNeighbourFunc(t1, t2)
}

func newPositions() board.Positions {
	return board.NewPositions([]board.Territory{bud, gal, vie, boh, lon})

}
