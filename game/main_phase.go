package game

import (
	"sort"

	"github.com/burrbd/dip/game/order"
	"github.com/burrbd/dip/game/order/board"
)

// Handler main game phase order handler and resolver
type Handler interface {
	ApplyOrders(order.Set, board.Manager)
	ResolveOrders(board.Manager)
}

// MainPhaseHandler implments Handler
type MainPhaseHandler struct {
	ArmyGraph board.Graph
}

// ApplyOrders applies orders for a turn. Satisifies Handler interface
func (h MainPhaseHandler) ApplyOrders(orders order.Set, positions board.Manager) {
	all := positions.Units()
	for _, unit := range all {
		for _, move := range orders.Moves {
			from, to := move.From.Abbr, move.To.Abbr
			neighbours, _ := h.ArmyGraph.IsNeighbour(from, to)
			if !neighbours || from != unit.Position().Territory.Abbr {
				continue
			}
			strength := h.strength(move, orders)
			positions.Move(unit, move.To, strength)
			break
		}
		for _, hold := range orders.Holds {
			pos := unit.Position()
			if pos == nil || pos.Cause != board.Added || unit.Position().Territory.Abbr != hold.At.Abbr {
				continue
			}
			strength := h.holdStrength(hold, orders)
			positions.Hold(unit, strength)
			break
		}
	}
}

// ResolveOrders resolves orders for a turn. Satisifies Handler interface
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

func (h MainPhaseHandler) strength(move order.Move, orders order.Set) int {
	strength := 0
	for _, support := range orders.MoveSupports {
		if support.Move.From.Abbr == move.From.Abbr &&
			support.Move.To.Abbr == move.To.Abbr &&
			!h.moveSupportCut(support, orders) {
			strength++
		}
	}
	return strength
}

func (h MainPhaseHandler) holdStrength(hold order.Hold, orders order.Set) int {
	strength := 0
	for _, support := range orders.HoldSupports {
		if support.Hold.At.Abbr == hold.At.Abbr {
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
