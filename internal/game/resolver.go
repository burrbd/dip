package game

import (
	"sort"

	"github.com/burrbd/diplomacy/internal/game/order/board"
)

type Resolver interface {
	Resolve(board.PositionMap) (board.PositionMap, error)
}

type MainPhaseResolver struct {
}

func (r MainPhaseResolver) Resolve(recorder board.PositionMap) (board.PositionMap, error) {
	for {
		units := recorder.GetConflict()
		if len(units) == 0 {
			return recorder, nil
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
				u.Defeated = defeated
			} else {
				recorder.Bounce(u, *u.PrevPosition)
			}
		}
	}
}
