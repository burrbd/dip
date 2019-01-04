package board

type territoryConflicts map[string][]*Unit

func (conflicts territoryConflicts) add(u *Unit) {
	terr := u.Position().Territory.Abbr
	if _, ok := conflicts[terr]; !ok {
		conflicts[terr] = make([]*Unit, 0)
	}
	conflicts[terr] = append(conflicts[terr], u)
}

func (conflicts territoryConflicts) del(u *Unit) {
	terr := u.Position().Territory.Abbr
	units, ok := conflicts[terr]
	if !ok {
		return
	}
	for i, unit := range units {
		if u == unit {
			units = removeIndex(i, units)
		}
	}
	conflicts[terr] = units
}

func removeIndex(i int, units []*Unit) []*Unit {
	copy(units[i:], units[i+1:])
	units[len(units)-1] = nil
	return units[:len(units)-1]
}
