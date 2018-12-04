package game_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/burrbd/diplomacy/internal/game"

	"github.com/cheekybits/is"

	"github.com/burrbd/diplomacy/internal/game/order"
	"github.com/burrbd/diplomacy/internal/game/order/board"
)

var cases = []orderCase{
	{
		description: "given a unit moves unchallenged, then unit changes territory",
		orders: []*orderResult{
			{order: "A Bud-Vie", result: "vie"},
		},
	},
	{
		description: "given two units attack same territory without support, then neither unit wins territory",
		orders: []*orderResult{
			{order: "A Bud-Vie", result: "bud"},
			{order: "A Gal-Vie", result: "gal"},
		},
	},
	{
		description: "given units attack in circular chain without support, then all attacking units bounce back",
		orders: []*orderResult{
			{order: "A Bud-Gal", result: "bud"},
			{order: "A Gal-Vie", result: "gal"},
			{order: "A Vie H", result: "vie"},
		},
	},
	{
		description: "given two units attack an empty territory, then supported attack wins",
		orders: []*orderResult{
			{order: "A Gal-Vie", result: "vie"},
			{order: "A Boh S A Gal-Vie", result: "boh"},
			{order: "A Bud-Vie", result: "bud"},
		},
	},
	{
		description: "given two units attack an empty territory, then unit with greatest support wins",
		orders: []*orderResult{
			{order: "A Gal-Vie", result: "vie"},
			{order: "A Boh S A Gal-Vie", result: "boh"},
			{order: "A Tri S A Gal-Vie", result: "tri"},
			{order: "A Bud-Vie", result: "bud"},
			{order: "A Tyr S A Bud-Vie", result: "tyr"},
		},
	},
	{
		description: "given unit holds territory, then unit remains on territory",
		orders: []*orderResult{
			{order: "A Vie H", result: "vie"},
		},
	},
	{
		description: "given unit attacks territory and defending territory attacks support, then attacking " +
			"unit still wins",
		orders: []*orderResult{
			{order: "A Gal-Vie", result: "vie"},
			{order: "A Boh S A Gal-Vie", result: "boh"},
			{order: "A Vie-Boh", result: "vie", defeated: true},
		},
	},
	{
		description: "given two units attack each other (counterattack), then both units bounce",
		orders: []*orderResult{
			{order: "A Vie-Bud", result: "vie"},
			{order: "A Bud-Vie", result: "bud"},
		},
	},
	{
		description: "given a counterattack, and another attacks one counterattack party," +
			"then all units bounce",
		orders: []*orderResult{
			{order: "A Vie-Bud", result: "vie"},
			{order: "A Bud-Vie", result: "bud"},
			{order: "A Boh-Vie", result: "boh"},
		},
	},
	{
		description: "given a counterattack, and another unit attacks one counterattack party " +
			"with support then supported unit wins",
		orders: []*orderResult{
			{order: "A Vie-Bud", result: "vie", defeated: true},
			{order: "A Bud-Vie", result: "bud"},
			{order: "A Boh-Vie", result: "vie"},
			{order: "A Tyr S A Boh-Vie", result: "tyr"},
		},
	},
	{
		description: "given a counterattack and a supported second attack, where one counterattack party has support, " +
			"then all units bounce",
		orders: []*orderResult{
			{order: "A Vie-Bud", result: "vie"},
			{order: "A Bud-Vie", result: "bud"},
			{order: "A Sil S A Bud-Vie", result: "sil"},
			{order: "A Boh-Vie", result: "boh"},
			{order: "A Tyr S A Boh-Vie", result: "tyr"},
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

	positions := board.NewPositions()
	positions.Add(unit)

	orders := order.Set{}
	orders.AddMove(order.Move{From: gal, To: lon})

	resolved, err := resolver.Resolve(orders, positions)

	is.NoErr(err)
	is.Nil(resolved.Units["lon"])
	is.Equal(unit, resolved.Units["gal"][0])
}

type orderResult struct {
	order    string
	result   string
	defeated bool
	unit     *board.Unit
}

type orderCase struct {
	description string
	givenMap    []string
	orders      []*orderResult
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

	for i, orderCase := range cases {
		logTableHeading(t, orderCase.description, i)

		positions := board.NewPositions()
		orders := order.Set{}
		resolver := game.MainPhaseResolver{ArmyGraph: graph}

		for _, orderResult := range orderCase.orders {
			o, err := order.Decode(orderResult.order)
			is.NoErr(err)
			var terr board.Territory
			switch v := o.(type) {
			case order.Move:
				terr = v.From
				orders.AddMove(v)
			case order.Hold:
				terr = v.Pos
				orders.AddHold(v)
			case order.MoveSupport:
				terr = v.By
				orders.AddMoveSupport(v)
			case order.MoveConvoy:
				terr = v.By
				orders.AddMoveConvoy(v)
			}

			u := &board.Unit{Position: terr}

			positions.Add(u)
			orderResult.unit = u

			logTableRow(t, *orderResult)
		}

		resolvedPositions, err := resolver.Resolve(orders, positions)
		is.NoErr(err)

		for _, orderResult := range orderCase.orders {
			is.NotNil(orderResult.unit)
			is.Equal(orderResult.defeated, orderResult.unit.Defeated)
			is.Equal(orderResult.result, orderResult.unit.Position.Abbr)
		}

		positionTotal := 0
		for _, positionUnits := range resolvedPositions.Units {
			positionTotal += len(positionUnits)
		}
		is.Equal(positionTotal, len(orderCase.orders))
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

func logTableHeading(t *testing.T, desc string, i int) {
	if i != 0 {
		t.Log("")
	}
	t.Log(desc)
	t.Log("")
	t.Log("  | order             | result | defeated |")
	t.Log("  +---------------------------------------+")
}

func logTableRow(t *testing.T, o orderResult) {
	t.Logf("  | %s%s| %s    | %t%s|",
		o.order,
		strings.Repeat(" ", 18-len(o.order)),
		o.result,
		o.defeated,
		strings.Repeat(" ", 9-len(fmt.Sprintf("%t", o.defeated))))
}
