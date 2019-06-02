package game

import (
	"github.com/burrbd/dip/game/order"
	"github.com/burrbd/dip/game/order/board"
)

// OrderHandler applies orders to the board.
type OrderHandler struct {
	Matcher order.Matcher
}

// ApplyOrders applies orders for a turn. Satisifies Handler interface
func (h OrderHandler) ApplyOrders(orders order.Set, manager board.Manager) {
	for unit, pos := range manager.Positions() {
		for _, move := range orders.Moves {
			if !h.Matcher.MatchMove(move, pos, "") {
				continue
			}
			manager.Move(unit, move.To, h.moveStrength(move, orders))
			break
		}
		for _, hold := range orders.Holds {
			if !h.Matcher.MatchHold(hold, pos, "") {
				continue
			}
			manager.Hold(unit, h.holdStrength(hold, orders))
			break
		}
	}
}

func (h OrderHandler) moveStrength(move order.Move, orders order.Set) (strength int) {
	for _, support := range orders.MoveSupports {
		if h.Matcher.MatchMoveSupport(support, move) && !h.Matcher.MoveSupportCut(support, orders.Moves) {
			strength++
		}
	}
	return
}

func (h OrderHandler) holdStrength(hold order.Hold, orders order.Set) (strength int) {
	for _, support := range orders.HoldSupports {
		if h.Matcher.MatchHoldSupport(support, hold) && !h.Matcher.HoldSupportCut(support, orders.Moves) {
			strength++
		}
	}
	return
}
