package board

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
	Strength     int
	Defeated     bool
}

func (u *Unit) SetNewPosition(terr Territory) {
	if u.PrevPosition == nil {
		prev := u.Position
		u.PrevPosition = &prev
	}
	u.Position = terr
}

func (u *Unit) AtOrigin() bool {
	return u.PrevPosition == nil || *u.PrevPosition == u.Position
}

func UnitsByStrength(units []*Unit) strengthSorter {
	return strengthSorter{units}
}

type strengthSorter struct {
	units []*Unit
}

func (s strengthSorter) Len() int {
	return len(s.units)
}

func (s strengthSorter) Swap(i, j int) {
	s.units[i], s.units[j] = s.units[j], s.units[i]
}

func (s strengthSorter) Less(i, j int) bool {
	return s.units[i].Strength > s.units[j].Strength
}
