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
	PhaseHistory []Position
}

type Position struct {
	Territory Territory
	Strength  int
	Cause     PositionEvent
}

type PositionEvent int

const (
	Added PositionEvent = iota
	Held
	Moved
	Bounced
	Defeated
)

func (u *Unit) AtOrigin() bool {
	pos := u.Position()
	if pos == nil {
		return false
	}
	return pos.Territory.Abbr == u.PhaseHistory[0].Territory.Abbr
}

func (u *Unit) Defeated() bool {
	return u.Position().Cause == Defeated
}

func (u *Unit) Position() *Position {
	n := len(u.PhaseHistory)
	if n == 0 {
		return nil
	}
	return &u.PhaseHistory[n-1]
}

func (u *Unit) PrevPosition() *Position {
	n := len(u.PhaseHistory)
	if n > 1 {
		return &u.PhaseHistory[n-2]
	}
	return nil
}

func UnitPositionsByStrength(units []*Unit) strengthSorter {
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
	if s.units[i].Position() == nil {
		return false
	}
	if s.units[j].Position() == nil {
		return true
	}
	return s.units[i].Position().Strength > s.units[j].Position().Strength
}
