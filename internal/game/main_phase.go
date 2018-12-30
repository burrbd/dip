package game

import (
	"sort"

	"github.com/burrbd/diplomacy/internal/game/order"
	"github.com/burrbd/diplomacy/internal/game/order/board"
)

type Handler interface {
	ApplyOrders(order.Set, board.PositionManager)
	ResolveOrderConflicts(board.PositionMap) (board.PositionMap, error)
}

type MainPhaseHandler struct {
	ArmyGraph board.Graph
}

func (h MainPhaseHandler) ApplyOrders(orders order.Set, positions board.PositionManager) {
	all := positions.Units()
	for _, unit := range all {
		for _, move := range orders.Moves {
			from, to := move.From.Abbr, move.To.Abbr
			neighbours, _ := h.ArmyGraph.IsNeighbour(from, to)
			if unit.PrevPosition() != nil || from != unit.Position().Territory.Abbr || !neighbours {
				continue
			}
			positions.Move(unit, move.To, h.strength(unit, move, orders))
		}
	}
}

func (h MainPhaseHandler) ResolveOrderConflicts(positions board.PositionManager) error {
	for {
		units := positions.Conflict()
		if len(units) == 0 {
			return nil
		}
		var defeated bool
		sort.Sort(board.UnitPositionsByStrength(units))
		if units[0].Position().Strength > units[1].Position().Strength {
			units, defeated = units[1:], true
		}
		for _, u := range units {
			if u == nil {
				continue
			}
			if u.AtOrigin() && defeated {
				positions.SetDefeated(u)
			} else if !u.AtOrigin() {
				positions.Bounce(u)
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
