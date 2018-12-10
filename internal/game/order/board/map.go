package board

import (
	"sort"
	"strings"
)

type PositionMap struct {
	Units            map[string][]*Unit
	CounterConflicts map[string][]*Unit
}

func NewPositionMap() PositionMap {
	return PositionMap{
		Units:            make(map[string][]*Unit),
		CounterConflicts: make(map[string][]*Unit)}
}

func (m PositionMap) AllUnits() []*Unit {
	ret := make([]*Unit, 0)
	for _, units := range m.Units {
		ret = append(ret, units...)
	}
	return ret
}

func (m PositionMap) GetConflict() []*Unit {
	for _, units := range m.CounterConflicts {
		nonRetreatingUnits := unitFilter(units, func(u *Unit) bool { return u != nil && !u.Defeated })
		if len(nonRetreatingUnits) == 2 {
			unitsCopy := make([]*Unit, 2)
			copy(unitsCopy, units)
			return unitsCopy
		}
	}
	for _, units := range m.Units {
		nonRetreatingUnits := unitFilter(units, func(u *Unit) bool { return u != nil && !u.Defeated })
		conflicts := len(nonRetreatingUnits)
		if conflicts > 1 {
			unitsCopy := make([]*Unit, conflicts)
			copy(unitsCopy, units)
			return unitsCopy
		}
	}
	return nil
}

func (m PositionMap) Add(u *Unit) {
	terr := u.Position.Abbr
	if _, ok := m.Units[terr]; !ok {
		m.Units[terr] = make([]*Unit, 0)
	}
	m.Units[terr] = append(m.Units[terr], u)
}

func (m PositionMap) Del(u *Unit) {
	terr := u.Position.Abbr
	units, ok := m.Units[terr]
	if !ok {
		return
	}
	for i, unit := range m.Units[terr] {
		if u == unit {
			units = removeIndex(i, units)
		}
	}
	m.Units[terr] = units
}

func (m PositionMap) Move(u *Unit, next Territory) {
	m.update(u, next)
	m.addCounterConflict(u)
}

func (m PositionMap) Bounce(u *Unit, next Territory) {
	m.delCounterConflict(u)
	m.update(u, next)
}

func (m PositionMap) update(u *Unit, next Territory) {
	m.Del(u)
	u.SetNewPosition(next)
	m.Add(u)
}

func unitFilter(units []*Unit, f func(*Unit) bool) []*Unit {
	filtered := make([]*Unit, 0)
	for _, u := range units {
		if f(u) {
			filtered = append(filtered, u)
		}
	}
	return filtered
}

func (m PositionMap) addCounterConflict(u *Unit) {
	key := pairKey(*u.PrevPosition, u.Position)
	if _, ok := m.CounterConflicts[key]; !ok {
		m.CounterConflicts[key] = make([]*Unit, 0, 2)
	}
	m.CounterConflicts[key] = append(m.CounterConflicts[key], u)
}

func (m PositionMap) delCounterConflict(u *Unit) {
	if u.PrevPosition == nil {
		return
	}
	key := pairKey(*u.PrevPosition, u.Position)
	units := m.CounterConflicts[key]
	for i, cu := range units {
		if u == cu {
			m.CounterConflicts[key] = removeIndex(i, units)
		}
	}
}

func pairKey(t1, t2 Territory) string {
	s := []string{t1.Abbr, t2.Abbr}
	sort.Strings(s)
	return strings.Join(s, "")
}

func removeIndex(i int, units []*Unit) []*Unit {
	copy(units[i:], units[i+1:])
	units[len(units)-1] = nil
	return units[:len(units)-1]
}
