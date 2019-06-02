package order

import (
	"github.com/burrbd/dip/game/order/board"
)

// Matcher matches orders against unit positions.
type Matcher interface {
	MatchMove(Move, board.Position, string) bool
	MatchHold(Hold, board.Position, string) bool
	MatchMoveSupport(MoveSupport, Move) bool
	MatchHoldSupport(HoldSupport, Hold) bool
	MoveSupportCut(MoveSupport, []Move) bool
	HoldSupportCut(HoldSupport, []Move) bool
}

// PositionMatcher implements Matcher.
type PositionMatcher struct {
	ArmyGraph board.Graph
}

// MatchMove matches move order against a position.
func (m PositionMatcher) MatchMove(move Move, pos board.Position, country string) bool {
	if move.Country != country {
		return false
	}
	from, to := move.From.Abbr, move.To.Abbr
	neighbours, _ := m.ArmyGraph.IsNeighbour(from, to)
	return neighbours && pos.Cause == board.Added && from == pos.Territory.Abbr
}

// MatchHold matches hold order against a position.
func (m PositionMatcher) MatchHold(hold Hold, pos board.Position, country string) bool {
	return country == hold.Country &&
		pos.Cause == board.Added &&
		pos.Territory.Abbr == hold.At.Abbr
}

// MatchMoveSupport matches move support against move.
func (m PositionMatcher) MatchMoveSupport(sup MoveSupport, move Move) bool {
	return sup.Move.From.Abbr == move.From.Abbr &&
		sup.Move.To.Abbr == move.To.Abbr
}

// MatchHoldSupport matches hold support against hold.
func (m PositionMatcher) MatchHoldSupport(sup HoldSupport, hold Hold) bool {
	return sup.Hold.At.Abbr == hold.At.Abbr
}

// MoveSupportCut checks if a move support has been cut.
func (m PositionMatcher) MoveSupportCut(sup MoveSupport, moves []Move) bool {
	for _, cut := range moves {
		if cut.To == sup.By && cut.From != sup.Move.To {
			return true
		}
	}
	return false
}

// HoldSupportCut checks if a hold support has been cut.
func (m PositionMatcher) HoldSupportCut(sup HoldSupport, moves []Move) bool {
	// TODO: implement
	return false
}
