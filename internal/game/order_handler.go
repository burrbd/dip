package game

import (
	"github.com/burrbd/diplomacy/internal/game/order"
	"github.com/burrbd/diplomacy/internal/game/order/board"
)

type ConflictRecorder interface {
	Move(unit *board.Unit, territory board.Territory)
	Bounce(unit *board.Unit, territory board.Territory)
	AllUnits() []*board.Unit
}

type OrderHandler struct {
	ArmyGraph board.Graph
}

func (h OrderHandler) Handle(orders order.Set, recorder ConflictRecorder) {
	all := recorder.AllUnits()
	for _, unit := range all {
		for _, move := range orders.Moves {
			from, to := move.From.Abbr, move.To.Abbr
			if unit.Position.Abbr != from || unit.PrevPosition != nil {
				continue
			}
			neighbours, _ := h.ArmyGraph.IsNeighbour(from, to)
			if neighbours && from == unit.Position.Abbr {
				h.setMoveStrength(unit, move, orders)
				recorder.Move(unit, move.To)
			}
		}
	}
}

func (h OrderHandler) setMoveStrength(u *board.Unit, move *order.Move, orders order.Set) {
	for _, support := range orders.MoveSupports {
		if support.Move.From.Abbr == move.From.Abbr &&
			support.Move.To.Abbr == move.To.Abbr &&
			!h.moveSupportCut(*support, orders) {
			u.Strength++
		}
	}
}

func (h OrderHandler) moveSupportCut(support order.MoveSupport, orders order.Set) bool {
	for _, cutMove := range orders.Moves {
		if cutMove.To == support.By && cutMove.From != support.Move.To {
			return true
		}
	}
	return false
}
