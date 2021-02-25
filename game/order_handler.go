package game

import (
	"github.com/burrbd/dip/game/order"
	"github.com/burrbd/dip/game/order/board"
)

type validator interface {
	ValidateMove(board.Unit, order.Move) error
	ValidateMoveSupport(board.Unit, order.MoveSupport) error
}

// OrderHandler applies orders to the board.
type OrderHandler struct {
	Validator validator
}

// ApplyOrders applies orders for a turn. Satisifies Handler interface
func (h OrderHandler) ApplyOrders(orders order.Set, manager board.Manager) {
boardPositionLoop:
	for unit, pos := range manager.Positions() {
		for _, move := range orders.Moves {
			if matchMoveToPosition(pos, move, "") {
				manager.Move(unit, move.To, h.moveStrength(move, orders))
				continue boardPositionLoop
			}
		}
		manager.Hold(unit, h.positionStrength(pos, orders))
	}
}

func (h OrderHandler) moveStrength(move order.Move, orders order.Set) (strength int) {
	for _, support := range orders.MoveSupports {
		if support.Move.From.Abbr == move.From.Abbr &&
			support.Move.To.Abbr == move.To.Abbr &&
			!moveSupportCut(support, orders.Moves) {
			strength++
		}
	}
	return
}

func (h OrderHandler) positionStrength(pos board.Position, orders order.Set) (strength int) {
	for _, support := range orders.HoldSupports {
		if support.Hold.At.Abbr == pos.Territory.Abbr &&
			!holdSupportCut(support, orders.Moves) {
			strength++
		}
	}
	return
}

func matchMoveToPosition(pos board.Position, move order.Move, country string) bool {
	return move.Country == country && move.From.Abbr == pos.Territory.Abbr
}

func moveSupportCut(sup order.MoveSupport, moves []order.Move) bool {
	for _, cut := range moves {
		if sup.By.Is(cut.To) && sup.Move.To.IsNot(cut.From) {
			return true
		}
	}
	return false
}

func holdSupportCut(sup order.HoldSupport, moves []order.Move) bool {
	for _, cut := range moves {
		if sup.By.Is(cut.To) {
			return true
		}
	}
	return false
}
