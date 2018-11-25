package game

import (
	"github.com/burrbd/diplomacy/internal/game/order"
	"github.com/burrbd/diplomacy/internal/game/order/board"
)

type RetreatPhaseResolver struct{}

func (r RetreatPhaseResolver) Resolve(set order.Set, positions board.Positions) (board.Positions, error) {
	return board.Positions{}, nil
}

type BuildPhaseResolver struct{}

func (r BuildPhaseResolver) Resolve(set order.Set, positions board.Positions) (board.Positions, error) {
	return board.Positions{}, nil
}
