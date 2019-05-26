package game_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/cheekybits/is"

	"github.com/burrbd/dip/game"
	"github.com/burrbd/dip/game/order"
	"github.com/burrbd/dip/game/order/board"
)

var specs = []spec{
	{
		description: "given a unit moves unchallenged, then unit changes territory",
		orders: []*result{
			{order: "A Bud-Vie", position: "vie"},
		},
	},
	{
		description: "given two units attack same territory without support, then neither unit wins territory",
		orders: []*result{
			{order: "A Bud-Vie", position: "bud"},
			{order: "A Gal-Vie", position: "gal"},
		},
	},
	{
		description: "given units attack in circular chain without support, then all attacking units bounce back",
		orders: []*result{
			{order: "A Bud-Gal", position: "bud"},
			{order: "A Gal-Vie", position: "gal"},
			{order: "A Vie H", position: "vie"},
		},
	},
	{
		description: "given two units attack an empty territory, then supported attack wins",
		orders: []*result{
			{order: "A Gal-Vie", position: "vie"},
			{order: "A Boh S A Gal-Vie", position: "boh"},
			{order: "A Bud-Vie", position: "bud"},
		},
	},
	{
		description: "given two units attack an empty territory, then unit with greatest support wins",
		orders: []*result{
			{order: "A Gal-Vie", position: "vie"},
			{order: "A Boh S A Gal-Vie", position: "boh"},
			{order: "A Tri S A Gal-Vie", position: "tri"},
			{order: "A Bud-Vie", position: "bud"},
			{order: "A Tyr S A Bud-Vie", position: "tyr"},
		},
	},
	{
		description: "given unit holds territory, then unit remains on territory",
		orders: []*result{
			{order: "A Vie H", position: "vie"},
		},
	},
	{
		description: "given unit attacks territory and defending territory attacks support, " +
			"then attacking unit still wins",
		orders: []*result{
			{order: "A Gal-Vie", position: "vie"},
			{order: "A Boh S A Gal-Vie", position: "boh"},
			{order: "A Vie-Boh", position: "vie", defeated: true},
		},
	},
	{
		description: "given two units attack each other (counterattack), then both units bounce",
		orders: []*result{
			{order: "A Vie-Bud", position: "vie"},
			{order: "A Bud-Vie", position: "bud"},
		},
	},
	{
		description: "given a counterattack, and another attacks one counterattack party," +
			"then all units bounce",
		orders: []*result{
			{order: "A Vie-Bud", position: "vie"},
			{order: "A Bud-Vie", position: "bud"},
			{order: "A Boh-Vie", position: "boh"},
		},
	},
	{
		description: "given a counterattack, and another unit attacks one counterattack party with support, " +
			"then supported unit wins",
		orders: []*result{
			{order: "A Vie-Bud", position: "vie", defeated: true},
			{order: "A Bud-Vie", position: "bud"},
			{order: "A Boh-Vie", position: "vie"},
			{order: "A Tyr S A Boh-Vie", position: "tyr"},
		},
	},
	{
		description: "given a counterattack and a supported second attack, where one counterattack party has support, " +
			"then all units bounce",
		orders: []*result{
			{order: "A Vie-Bud", position: "vie"},
			{order: "A Bud-Vie", position: "bud"},
			{order: "A Sil S A Bud-Vie", position: "sil"},
			{order: "A Boh-Vie", position: "boh"},
			{order: "A Tyr S A Boh-Vie", position: "tyr"},
		},
	},
	{
		description: "given a unit holds and another unit supports holding unit," +
			"then both units remain in position",
		orders: []*result{
			{order: "A Vie H", position: "vie"},
			{order: "A Bud S A Vie", position: "bud"},
		},
	},
	{
		description: "given a unit holds and is supported, and is attacked by unit with equal strength, " +
			"then attacking unit bounces",
		orders: []*result{
			{order: "A Vie H", position: "vie"},
			{order: "A Bud S A Vie", position: "bud"},
			{order: "A Boh-Vie", position: "boh"},
			{order: "A Tyr S A Boh-Vie", position: "tyr"},
		},
	},
}

type result struct {
	order    string
	position string
	defeated bool
	unit     *board.Unit
}

type spec struct {
	description string
	givenMap    []string
	orders      []*result
	focus       bool
}

func TestMainPhaseResolver_ResolveCases(t *testing.T) {
	is := is.New(t)
	graph := mockGraph{
		IsNeighbourFunc: func(t1, t2 string) (bool, error) { return true, nil },
	}

	for _, spec := range filter(specs) {
		t.Run(spec.description, func(t *testing.T) {
			logTableHeading(t)
			positionManager := board.NewPositionManager()

			orders := order.Set{}

			for _, result := range spec.orders {
				o, err := order.Decode(result.order)
				is.NoErr(err)
				var terr board.Territory
				switch v := o.(type) {
				case order.Move:
					terr = v.From
					orders.AddMove(v)
				case order.Hold:
					terr = v.At
					orders.AddHold(v)
				case order.MoveSupport:
					terr = v.By
					orders.AddMoveSupport(v)
				case order.HoldSupport:
					terr = v.By
					orders.AddHoldSupport(v)
				case order.MoveConvoy:
					terr = v.By
					orders.AddMoveConvoy(v)
				}

				u := &board.Unit{}
				positionManager.AddUnit(u, terr)
				result.unit = u
			}

			orderHandler := game.MainPhaseHandler{ArmyGraph: graph}
			orderHandler.ApplyOrders(orders, positionManager)
			orderHandler.ResolveOrders(positionManager)

			for _, result := range spec.orders {
				logTableRow(t, *result)
				is.NotNil(result.unit)
				is.Equal(result.defeated, result.unit.Defeated())
				is.Equal(result.position, result.unit.Position().Territory.Abbr)
			}
			is.Equal(len(spec.orders), len(positionManager.Units()))
		})
	}
}

func filter(specs []spec) []spec {
	focused := make([]spec, 0)
	for _, c := range specs {
		if c.focus {
			focused = append(focused, c)
		}
	}
	if len(focused) == 0 {
		return specs
	}
	return focused
}

type mockGraph struct {
	IsNeighbourFunc func(t1, t2 string) (bool, error)
}

func (g mockGraph) IsNeighbour(t1, t2 string) (bool, error) {
	return g.IsNeighbourFunc(t1, t2)
}

func logTableHeading(t *testing.T) {
	t.Log("  | order             | result | defeated |")
	t.Log("  +---------------------------------------+")
}

func logTableRow(t *testing.T, o result) {
	t.Logf("  | %s%s| %s    | %t%s|",
		o.order,
		strings.Repeat(" ", 18-len(o.order)),
		o.position,
		o.defeated,
		strings.Repeat(" ", 9-len(fmt.Sprintf("%t", o.defeated))))
}
