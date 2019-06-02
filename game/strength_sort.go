package game

import "github.com/burrbd/dip/game/order/board"

// UnitPositionsByStrength returns a unit strength sorter
func UnitPositionsByStrength(manager board.Manager, units []*board.Unit) strengthSorter {
	return strengthSorter{manager, units}
}

type strengthSorter struct {
	manager board.Manager
	units   []*board.Unit
}

func (s strengthSorter) Len() int {
	return len(s.units)
}

func (s strengthSorter) Swap(i, j int) {
	s.units[i], s.units[j] = s.units[j], s.units[i]
}

func (s strengthSorter) Less(i, j int) bool {
	if s.manager.Position(s.units[i]) == nil {
		return false
	}
	if s.manager.Position(s.units[j]) == nil {
		return true
	}
	return s.manager.Position(s.units[i]).Strength > s.manager.Position(s.units[j]).Strength
}
