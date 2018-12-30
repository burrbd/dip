package board

import (
	"sort"
	"strings"
)

type PositionManager interface {
	Units() []*Unit
	Move(unit *Unit, territory Territory, strength int)
	Bounce(unit *Unit)
	SetDefeated(unit *Unit)
	Conflict() []*Unit
}

type PositionMap struct {
	conflicts        map[string][]*Unit
	counterConflicts map[string][]*Unit
}

func NewPositionMap(units []*Unit) PositionMap {
	conflicts := make(map[string][]*Unit)
	for _, u := range units {
		if _, ok := conflicts[u.Position().Territory.Abbr]; !ok {
			conflicts[u.Position().Territory.Abbr] = make([]*Unit, 0)
		}
		conflicts[u.Position().Territory.Abbr] = append(conflicts[u.Position().Territory.Abbr], u)
	}
	return PositionMap{
		conflicts:        conflicts,
		counterConflicts: make(map[string][]*Unit)}
}

func (m PositionMap) Units() []*Unit {
	ret := make([]*Unit, 0)
	for _, units := range m.conflicts {
		ret = append(ret, units...)
	}
	return ret
}

func (m PositionMap) Conflict() []*Unit {
	for _, units := range m.counterConflicts {
		nonRetreatingUnits := unitFilter(units, func(u *Unit) bool { return !u.Defeated() })
		if len(nonRetreatingUnits) == 2 {
			unitsCopy := make([]*Unit, 2)
			copy(unitsCopy, units)
			return unitsCopy
		}
	}
	for _, units := range m.conflicts {
		nonRetreatingUnits := unitFilter(units, func(u *Unit) bool { return !u.Defeated() })
		conflicts := len(nonRetreatingUnits)
		if conflicts > 1 {
			unitsCopy := make([]*Unit, conflicts)
			copy(unitsCopy, units)
			return unitsCopy
		}
	}
	return nil
}

func (m PositionMap) Move(u *Unit, next Territory, strength int) {
	m.del(u)
	u.PhaseHistory = append(u.PhaseHistory, Position{
		Territory: next, Strength: strength, Cause: Moved})
	m.add(u)
	m.addCounterConflict(u)
}

func (m PositionMap) Bounce(u *Unit) {
	prev := u.PrevPosition()
	if prev == nil {
		return
	}
	next := prev.Territory
	m.delCounterConflict(u)
	m.del(u)
	u.PhaseHistory = append(u.PhaseHistory, Position{
		Territory: next, Strength: 0, Cause: Bounced})
	m.add(u)
}

func (m PositionMap) SetDefeated(u *Unit) {
	u.PhaseHistory = append(u.PhaseHistory, Position{
		Territory: u.Position().Territory, Cause: Defeated})
}

func (m PositionMap) add(u *Unit) {
	terr := u.Position().Territory.Abbr
	if _, ok := m.conflicts[terr]; !ok {
		m.conflicts[terr] = make([]*Unit, 0)
	}
	m.conflicts[terr] = append(m.conflicts[terr], u)
}

func (m PositionMap) del(u *Unit) {
	terr := u.Position().Territory.Abbr
	units, ok := m.conflicts[terr]
	if !ok {
		return
	}
	for i, unit := range units {
		if u == unit {
			units = removeIndex(i, units)
		}
	}
	m.conflicts[terr] = units
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
	if u.Position().Cause != Moved || u.PrevPosition() == nil {
		return
	}
	key := pairKey(u.Position().Territory, u.PrevPosition().Territory)
	if _, ok := m.counterConflicts[key]; !ok {
		m.counterConflicts[key] = make([]*Unit, 0, 2)
	}
	m.counterConflicts[key] = append(m.counterConflicts[key], u)
}

func (m PositionMap) delCounterConflict(u *Unit) {
	if u.PrevPosition() == nil {
		return
	}
	key := pairKey(u.PrevPosition().Territory, u.Position().Territory)
	units := m.counterConflicts[key]
	for i, cu := range units {
		if u == cu {
			m.counterConflicts[key] = removeIndex(i, units)
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
