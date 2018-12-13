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
			if unit.Position.Abbr != from || unit.PrevPosition != nil {
				continue
			}
			neighbours, _ := h.ArmyGraph.IsNeighbour(from, to)
			if neighbours && from == unit.Position.Abbr {
				h.setMoveStrength(unit, move, orders)
				positions.Move(unit, move.To)
			}
		}
	}
}

func (h MainPhaseHandler) ResolveOrderConflicts(positions board.PositionManager) error {
	for {
		units := positions.Conflict()
		if len(units) == 0 {
			return nil
		}
		defeated := false
		sort.Sort(board.UnitsByStrength(units))
		if units[0].Strength > units[1].Strength {
			units, defeated = units[1:], true
		}
		for _, u := range units {
			if u == nil {
				continue
			}
			if u.AtOrigin() {
				positions.SetDefeated(u, defeated)
			} else {
				positions.Bounce(u, *u.PrevPosition)
			}
		}
	}
}

func (h MainPhaseHandler) setMoveStrength(u *board.Unit, move *order.Move, orders order.Set) {
	for _, support := range orders.MoveSupports {
		if support.Move.From.Abbr == move.From.Abbr &&
			support.Move.To.Abbr == move.To.Abbr &&
			!h.moveSupportCut(*support, orders) {
			u.Strength++
		}
	}
}

func (h MainPhaseHandler) moveSupportCut(support order.MoveSupport, orders order.Set) bool {
	for _, cutMove := range orders.Moves {
		if cutMove.To == support.By && cutMove.From != support.Move.To {
			return true
		}
	}
	return false
}
