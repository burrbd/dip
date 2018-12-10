package game

import (
	"github.com/burrbd/diplomacy/internal/game/order"
	"github.com/burrbd/diplomacy/internal/game/order/board"
)

type RetreatPhaseResolver struct{}

func (r RetreatPhaseResolver) Resolve(set order.Set, positions board.PositionMap) (board.PositionMap, error) {
	return board.PositionMap{}, nil
}

type BuildPhaseResolver struct{}

func (r BuildPhaseResolver) Resolve(set order.Set, positions board.PositionMap) (board.PositionMap, error) {
	return board.PositionMap{}, nil
}
