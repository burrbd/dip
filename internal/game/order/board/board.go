package board

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
	Units map[string][]*Unit
}

func NewPositions() Positions {
	return Positions{Units: make(map[string][]*Unit)}
}

func (p Positions) Add(territory string, u *Unit) {
	if _, ok := p.Units[territory]; !ok {
		p.Units[territory] = make([]*Unit, 0)
	}
	p.Units[territory] = append(p.Units[territory], u)
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
