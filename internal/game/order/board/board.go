package board

import "fmt"

type Graph interface {
	IsNeighbour(t1, t2 string) (bool, error)
}

type UnitType string

const (
	Army  UnitType = "army"
	Fleet UnitType = "fleet"
)

type PositionMap interface {
	Get(string) *Unit
	Add(string, *Unit) error
	Del(string) error
}

type Positions struct {
	territories map[string]Territory
	Units       map[string][]*Unit
}

func NewPositions(territories []Territory) Positions {
	territoryMap := make(map[string]Territory)
	for _, v := range territories {
		territoryMap[v.Abbr] = v
	}
	return Positions{
		Units:       make(map[string][]*Unit),
		territories: territoryMap}
}

func (p Positions) Territory(abbr string) (Territory, error) {
	t, ok := p.territories[abbr]
	if !ok {
		return Territory{}, fmt.Errorf("unknown territory %s", abbr)
	}
	return t, nil
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
	prev := u.Position
	u.PrevPosition = &prev
	u.Position = next
	p.Add(u)
	return nil
}

func (p Positions) ConflictHandler(f func(Territory, []*Unit)) {
	for terr, units := range p.Units {
		t, err := p.Territory(terr)
		if err != nil {
			continue
		}
		if len(units) > 1 {
			f(t, units)
		}
	}
}

func (p Positions) ConflictCount() int {
	conflicts := 0
	p.ConflictHandler(func(terr Territory, units []*Unit) { conflicts++ })
	return conflicts
}

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
}
