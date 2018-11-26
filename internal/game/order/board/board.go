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

func (p Positions) Add(t Territory, u *Unit) {
	if _, ok := p.Units[t.Abbr]; !ok {
		p.Units[t.Abbr] = make([]*Unit, 0)
	}
	p.Units[t.Abbr] = append(p.Units[t.Abbr], u)
}

func (p Positions) Del(t Territory, u *Unit) error {
	units, ok := p.Units[t.Abbr]
	if !ok {
		return fmt.Errorf("no units in t %s", t)
	}
	for i, unit := range p.Units[t.Abbr] {
		if u == unit {
			copy(units[i:], units[i+1:])
			units[len(units)-1] = nil
			units = units[:len(units)-1]
		}
	}
	p.Units[t.Abbr] = units
	return nil
}

func (p Positions) Update(prev, next Territory, u *Unit) error {
	if err := p.Del(prev, u); err != nil {
		return err
	}
	p.Add(next, u)
	if u.PrevPositions == nil {
		u.PrevPositions = make([]Territory, 0, 1)
	}
	u.PrevPositions = append(u.PrevPositions, prev)
	return nil
}

type Territory struct {
	Abbr string
	Name string
}

func (t Territory) ID() string {
	return t.Abbr
}

type Unit struct {
	Country       string
	Type          UnitType
	PrevPositions []Territory
}
