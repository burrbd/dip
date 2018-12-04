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
	MustRetreat  bool
}

type PositionMap interface {
	Get(string) *Unit
	Add(string, *Unit) error
	Del(string) error
}

type Positions struct {
	Units                  map[string][]*Unit
	CounterAttackConflicts map[string][]*Unit
}

func NewPositions() Positions {
	return Positions{
		Units: make(map[string][]*Unit),
		CounterAttackConflicts: make(map[string][]*Unit)}
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
			copy(units[i:], units[i+1:])
			units[len(units)-1] = nil
			units = units[:len(units)-1]
		}
	}
	p.Units[terr] = units
	return nil
}

func (p Positions) Update(u *Unit, next Territory) error {
	if err := p.Del(u); err != nil {
		return err
	}
	if u.PrevPosition == nil {
		prev := u.Position
		u.PrevPosition = &prev
		u.Position = next
		p.addCounterattack(u)
	} else {
		p.removeCounterattack(u)
		u.Position = next
	}
	p.Add(u)
	return nil
}

func (p Positions) ConflictHandler(f func([]*Unit)) {
	for _, units := range p.CounterAttackConflicts {
		nonRetreatingUnits := unitFilter(units, func(u *Unit) bool { return u != nil && !u.MustRetreat })
		if len(nonRetreatingUnits) == 2 {
			f(nonRetreatingUnits)
		}
	}
	for _, units := range p.Units {
		nonRetreatingUnits := unitFilter(units, func(u *Unit) bool { return u != nil && !u.MustRetreat })
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

func (p Positions) addCounterattack(u *Unit) {
	key := counterAttackKey(*u.PrevPosition, u.Position)
	if _, ok := p.CounterAttackConflicts[key]; !ok {
		p.CounterAttackConflicts[key] = make([]*Unit, 0, 2)
	}
	p.CounterAttackConflicts[key] = append(p.CounterAttackConflicts[key], u)
}

func (p Positions) removeCounterattack(u *Unit) {
	key := counterAttackKey(*u.PrevPosition, u.Position)
	units := p.CounterAttackConflicts[key]
	for i, cu := range units {
		if u == cu {
			copy(units[i:], units[i+1:])
			units[len(units)-1] = nil
			units = units[:len(units)-1]
		}
	}
	p.CounterAttackConflicts[key] = units
}

func counterAttackKey(t1, t2 Territory) string {
	s := []string{t1.Abbr, t2.Abbr}
	sort.Strings(s)
	return strings.Join(s, "")
}
