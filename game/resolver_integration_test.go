package game_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/burrbd/dip/game"
	"github.com/burrbd/dip/game/order"
	"github.com/burrbd/dip/game/order/board"
	"github.com/cheekybits/is"
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
		description: "given a unit moves to a non-contiguous territory, then the move will be invalid",
		orders: []*result{
			{order: "A Vie-Lon", position: "vie"},
		},
	},
	{
		description: "given a supported unit holds and is attacked by unit with equal strength, " +
			"then attacking unit bounces",
		orders: []*result{
			{order: "A Vie H", position: "vie"},
			{order: "A Bud S A Vie", position: "bud"},
			{order: "A Boh-Vie", position: "boh"},
			{order: "A Tyr S A Boh-Vie", position: "tyr"},
		},
	},
	{
		description: "given a supported attack where the support cutter is itself bounced, " +
			"then support is still cut and attack fails (DATC 6.D.9)",
		orders: []*result{
			{order: "A Boh-Vie", position: "boh"},
			{order: "A Gal S A Boh-Vie", position: "gal"},
			{order: "A Vie H", position: "vie"},
			{order: "A Bud-Gal", position: "bud"},
			{order: "A Sil-Gal", position: "sil"},
		},
	},
	{
		description: "given a supported attack where both support cutters tie at the supporter's territory, " +
			"then support is still cut and attack fails (DATC 6.D.9)",
		orders: []*result{
			{order: "A Boh-Gal", position: "boh"},
			{order: "A Vie S A Boh-Gal", position: "vie"},
			{order: "A Gal H", position: "gal"},
			{order: "A Bud-Vie", position: "bud"},
			{order: "A Tri-Vie", position: "tri"},
		},
	},
	{
		description: "given a lone attack on a supporter that itself bounces, " +
			"then support is still cut and attack on supported territory fails (DATC 6.D.9)",
		orders: []*result{
			{order: "A Gal-Vie", position: "gal"},
			{order: "A Boh S A Gal-Vie", position: "boh"},
			{order: "A Vie H", position: "vie"},
			{order: "A Mun-Boh", position: "mun"},
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
	orders      []*result
	focus       bool
}

func TestMainPhaseResolver_ResolveCases(t *testing.T) {
	country := "a_country"

	for _, spec := range filter(specs) {
		t.Run(spec.description, func(t *testing.T) {

			is := is.New(t)

			logTableHeading(t)
			positionManager := board.NewPositionManager()

			orders := order.Set{}

			for _, result := range spec.orders {
				o, err := order.Decode(result.order, country)
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

				u := &board.Unit{Country: country}
				positionManager.AddUnit(u, board.LookupTerritory(terr.Abbr))
				result.unit = u
			}

			validator := order.NewValidator(board.CreateArmyGraph())
			orderHandler := game.OrderHandler{
				Validator: validator,
			}
			orderHandler.ApplyOrders(orders, positionManager)
			game.ResolveOrders(positionManager)

			for _, result := range spec.orders {
				logTableRow(t, *result)
				is.NotNil(result.unit)
				is.Equal(result.defeated, positionManager.Defeated(result.unit))
				is.Equal(result.position, positionManager.Position(result.unit).Territory.Abbr)
			}
			is.Equal(len(spec.orders), len(positionManager.Positions()))
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

func logTableHeading(t *testing.T) {
	t.Helper()
	t.Log("  | order             | result | defeated |")
	t.Log("  +---------------------------------------+")
}

func logTableRow(t *testing.T, o result) {
	t.Helper()
	t.Logf("  | %s%s| %s    | %t%s|",
		o.order,
		strings.Repeat(" ", 18-len(o.order)),
		o.position,
		o.defeated,
		strings.Repeat(" ", 9-len(fmt.Sprintf("%t", o.defeated))))
}
