package board

import (
	"fmt"
	"sort"
	"strings"
)

type Graph interface {
	IsNeighbour(t1, t2 string) (bool, error)
}

type UnitType string

const (
	Army  UnitType = "army"
	Fleet UnitType = "fleet"
)

type Territory struct {
	Abbr string
	Name string
}

func (t Territory) ID() string {
	return t.Abbr
}

type Unit struct {
	Country      string
	Type         UnitType
	Position     Territory
	PrevPosition *Territory
	Defeated     bool
}

func (u *Unit) SetNewPosition(terr Territory) {
	if u.PrevPosition == nil {
		prev := u.Position
		u.PrevPosition = &prev
	}
	u.Position = terr
}

func (u *Unit) OriginalPosition() bool {
	return u.PrevPosition == nil || *u.PrevPosition == u.Position
}

type Positions struct {
	Units            map[string][]*Unit
	CounterConflicts map[string][]*Unit
}

func NewPositions() Positions {
	return Positions{
		Units:            make(map[string][]*Unit),
		CounterConflicts: make(map[string][]*Unit)}
}

func (p Positions) Add(u *Unit) {
	terr := u.Position.Abbr
	if _, ok := p.Units[terr]; !ok {
		p.Units[terr] = make([]*Unit, 0)
	}
	p.Units[terr] = append(p.Units[terr], u)
}

func (p Positions) Del(u *Unit) error {
	terr := u.Position.Abbr
	units, ok := p.Units[terr]
	if !ok {
		return fmt.Errorf("no units in t %s", terr)
	}
	for i, unit := range p.Units[terr] {
		if u == unit {
			units = removeIndex(i, units)
		}
	}
	p.Units[terr] = units
	return nil
}

func (p Positions) Move(u *Unit, next Territory) error {
	if err := p.update(u, next); err != nil {
		return err
	}
	p.addCounterConflict(u)
	return nil
}

func (p Positions) Bounce(u *Unit, next Territory) error {
	p.delCounterConflict(u)
	if err := p.update(u, next); err != nil {
		return err
	}
	return nil
}

func (p Positions) update(u *Unit, next Territory) error {
	if err := p.Del(u); err != nil {
		return err
	}
	u.SetNewPosition(next)
	p.Add(u)
	return nil
}

func (p Positions) ConflictHandler(f func([]*Unit)) {
	for _, units := range p.CounterConflicts {
		nonRetreatingUnits := unitFilter(units, func(u *Unit) bool { return u != nil && !u.Defeated })
		if len(nonRetreatingUnits) == 2 {
			f(nonRetreatingUnits)
		}
	}
	for _, units := range p.Units {
		nonRetreatingUnits := unitFilter(units, func(u *Unit) bool { return u != nil && !u.Defeated })
		if len(nonRetreatingUnits) > 1 {
			f(nonRetreatingUnits)
		}
	}
}

func (p Positions) ConflictCount() int {
	conflicts := 0
	p.ConflictHandler(func(_ []*Unit) { conflicts++ })
	return conflicts
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

func (p Positions) addCounterConflict(u *Unit) {
	key := pairKey(*u.PrevPosition, u.Position)
	if _, ok := p.CounterConflicts[key]; !ok {
		p.CounterConflicts[key] = make([]*Unit, 0, 2)
	}
	p.CounterConflicts[key] = append(p.CounterConflicts[key], u)
}

func (p Positions) delCounterConflict(u *Unit) {
	key := pairKey(*u.PrevPosition, u.Position)
	units := p.CounterConflicts[key]
	for i, cu := range units {
		if u == cu {
			p.CounterConflicts[key] = removeIndex(i, units)
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
