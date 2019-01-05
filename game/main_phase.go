package game

import (
	"sort"

	"github.com/burrbd/dip/game/order"
	"github.com/burrbd/dip/game/order/board"
)

type Handler interface {
	ApplyOrders(order.Set, board.Manager)
	ResolveOrders(board.Manager)
}

type MainPhaseHandler struct {
	ArmyGraph board.Graph
}

func (h MainPhaseHandler) ApplyOrders(orders order.Set, positions board.Manager) {
	all := positions.Units()
	for _, unit := range all {
		for _, move := range orders.Moves {
			from, to := move.From.Abbr, move.To.Abbr
			neighbours, _ := h.ArmyGraph.IsNeighbour(from, to)
			if from != unit.Position().Territory.Abbr || !neighbours {
				continue
			}
			positions.Move(unit, move.To, h.strength(unit, move, orders))
			break
		}
	}
}

func (h MainPhaseHandler) ResolveOrders(positions board.Manager) {
	for {
		units := positions.Conflict()
		if units == nil {
			return
		}
		var defeated bool
		sort.Sort(board.UnitPositionsByStrength(units))
		if units[0].Position().Strength > units[1].Position().Strength {
			units, defeated = units[1:], true
		}
		for _, u := range units {
			if !u.AtOrigin() {
				positions.Bounce(u)
			} else if defeated {
				positions.SetDefeated(u)
			}
		}
	}
}

func (h MainPhaseHandler) strength(u *board.Unit, move *order.Move, orders order.Set) int {
	strength := 0
	for _, support := range orders.MoveSupports {
		if support.Move.From.Abbr == move.From.Abbr &&
			support.Move.To.Abbr == move.To.Abbr &&
			!h.moveSupportCut(*support, orders) {
			strength++
		}
	}
	return strength
}

func (h MainPhaseHandler) moveSupportCut(support order.MoveSupport, orders order.Set) bool {
	for _, cutMove := range orders.Moves {
		if cutMove.To == support.By && cutMove.From != support.Move.To {
			return true
		}
	}
	return false
}
